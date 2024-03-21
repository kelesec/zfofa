package unauth

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
	"zfofa/backend/core/fetch"
	"zfofa/backend/core/filelog"
	"zfofa/backend/engine/parser"
	"zfofa/backend/engine/scheduler"
	"zfofa/backend/engine/types"
	"zfofa/backend/engine/worker"
)

type ConcurrentRunner struct {
	WorkerCount    int
	Scheduler      types.Scheduler
	DoneCancelFunc context.CancelFunc
}

func (c *ConcurrentRunner) Run(ctx context.Context, searchType types.FofaSearchType) []types.Asset {
	filelog.Info("Start scraping data. --%s", "UnAuth::ConcurrentRunner::Run")
	log.Printf("Start scraping data. --%s", "UnAuth::ConcurrentRunner::Run")

	// 伪造百度爬虫
	headers := http.Header{}
	headers.Set(
		"User-Agent",
		"Mozilla/5.0 (compatible; Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)",
	)

	// 请求数量在20以内，可直接获取
	perAssetCount := searchType.SearchCount
	if perAssetCount > 0 && perAssetCount <= 20 {
		if perAssetCount < 2 {
			perAssetCount = 10
		}

		req := types.Request{
			Url:       parser.GenerateUrlWithUnAuth(searchType.Keywords, perAssetCount),
			Headers:   headers,
			ParseFunc: parser.UnAuthAssetsParser,
		}
		c.Scheduler.Push(req)
	}

	// 不在20内则构造新的请求：给关键字添加 after before 关键字，每次请求20个
	updateTime := types.CreateLastUpdateTimeWithCurrentTime(
		searchType.FofaConf.FofaToolConf.BetweenAfterTimeAndBeforeTime)
	keywordsWithUpTime := updateTime.GenKeywordsWithUpTime(searchType.Keywords)

	if perAssetCount < 0 || perAssetCount > 20 {
		perAssetCount = 20
		requestsCount := searchType.SearchCount/perAssetCount + 1

		for i := 0; i < requestsCount; i++ {
			req := types.Request{
				Url:       parser.GenerateUrlWithUnAuth(keywordsWithUpTime, perAssetCount),
				Headers:   headers,
				ParseFunc: parser.UnAuthAssetsParser,
			}
			c.Scheduler.Push(req)

			updateTime.SubDays()
			keywordsWithUpTime = updateTime.ReplaceKeywordsWithUpTime(keywordsWithUpTime)
		}
	}

	// 开始任务
	request := fetch.CreateRequest(true, 15*time.Second, searchType.Proxy)
	simpleWorker := worker.SimpleWorker{
		Request: request, Locker: &sync.Mutex{},
		ParseResultOutChan: make(chan types.ParseResult, cap(c.Scheduler.GetRequestInChan()))}
	for i := 0; i < c.WorkerCount; i++ {
		go simpleWorker.Worker(ctx, &c.Scheduler)
	}

	// 解析结果
	tryRequestBlankCount := 0
	var assets []types.Asset
	for {
		select {
		case result := <-simpleWorker.ParseResultOutChan:
			if len(result.Assets) != 0 {
				tryRequestBlankCount = 0
				assets = append(assets, result.Assets...)
				for _, asset := range result.Assets {
					log.Printf("[RES] %s\n", asset)
					filelog.Info("[RES] %s\n", asset)
				}
			} else if tryRequestBlankCount == searchType.FofaConf.FofaToolConf.MaxTryFetches {
				// 请求多次为空，则可能已经无数据，添加最后一次请求
				keywordsWithUpTime = updateTime.RemoveAfterTime(keywordsWithUpTime)
				req := types.Request{
					Url:       parser.GenerateUrlWithUnAuth(keywordsWithUpTime, perAssetCount),
					Headers:   headers,
					ParseFunc: parser.UnAuthAssetsParser,
				}

				if simpleSchedule, ok := c.Scheduler.(*scheduler.SimpleScheduler); ok {
					simpleSchedule.PushLastRequest(req)
				} else {
					c.Scheduler.Close()
				}

				tryRequestBlankCount++
				break
			} else {
				log.Printf("[BLANK] empty results: %s", result.Url)
				filelog.Info("[BLANK] empty results: %s", result.Url)
				tryRequestBlankCount++
			}

			/**退出条件：
			1、调度器关闭后，全部任务执行结束
			2、资产已经达到需求，并且任务全部执行完成
			*/
			if !simpleWorker.IsActive() {
				c.DoneCancelFunc()
				break
			}

			if len(assets) >= searchType.SearchCount && !simpleWorker.IsActive() {
				c.DoneCancelFunc()
				break
			}

			// 请求为空，但是资产数量不够，且未达到最大请求次数，继续添加请求
			if len(c.Scheduler.GetRequestInChan()) == 0 {
				reqCount := (searchType.SearchCount-len(assets))/perAssetCount + 1

				for i := 0; i < reqCount; i++ {
					req := types.Request{
						Url:       parser.GenerateUrlWithUnAuth(keywordsWithUpTime, perAssetCount),
						Headers:   headers,
						ParseFunc: parser.UnAuthAssetsParser,
					}

					c.Scheduler.Push(req)
					keywordsWithUpTime = updateTime.ReplaceKeywordsWithUpTime(keywordsWithUpTime)
					updateTime.SubDays()
				}
			}
		case <-ctx.Done():
			// 退出
			for {
				if !simpleWorker.IsActive() {
					return assets
				}
			}
		}
	}
}
