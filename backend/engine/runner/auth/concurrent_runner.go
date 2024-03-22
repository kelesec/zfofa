package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
	"zfofa/backend/core/fetch"
	"zfofa/backend/core/filelog"
	"zfofa/backend/engine/parser"
	"zfofa/backend/engine/types"
	"zfofa/backend/engine/worker"
)

type ConcurrentRunner struct {
	WorkerCount    int
	Scheduler      types.Scheduler
	DoneCancelFunc context.CancelFunc
}

func (c *ConcurrentRunner) Run(ctx context.Context, searchType types.FofaSearchType) []types.Asset {
	//获取请求URL
	reqUrls := parser.GenerateUrlsWithAuth(searchType)
	if reqUrls == nil {
		log.Println(errors.New("error: get request url fail"))
		filelog.Error(errors.New("error: get request url fail").Error())
		return nil
	}

	log.Printf("Start scraping data. --%s", "Auth::ConcurrentRunner::Run")
	filelog.Info("Start scraping data. --%s", "Auth::ConcurrentRunner::Run")

	// 提交任务
	for _, reqUrl := range reqUrls {
		r := types.Request{
			Url:       reqUrl,
			ParseFunc: parser.AuthAssetsParser,
		}
		c.Scheduler.Push(r)
	}

	// 构造cookie和请求session
	var fofaCookies = []*http.Cookie{{Name: "fofa_token", Value: searchType.FofaConf.Userinfo.FofaToken}}
	session := fetch.CreateSession(true, 15*time.Second, searchType.Proxy)
	session.AddCookie(searchType.Url, fofaCookies)

	// 开始任务
	simpleWorker := worker.SimpleWorker{
		Request: session, Locker: &sync.Mutex{},
		ParseResultOutChan: make(chan types.ParseResult, cap(c.Scheduler.GetRequestInChan()))}

	for i := 0; i < c.WorkerCount; i++ {
		go simpleWorker.Worker(ctx, &c.Scheduler)
	}

	// 解析结果
	var assets []types.Asset
	var resNumber int
	for {
		select {
		case result := <-simpleWorker.ParseResultOutChan:
			assets = append(assets, result.Assets...)
			for _, asset := range result.Assets {
				log.Printf("[RES] %s\n", asset)
				filelog.Info("[RES] %s\n", asset)
			}

			resNumber++
		case <-ctx.Done():
			for {
				if !simpleWorker.IsActive() {
					return assets
				}
			}
		}

		if len(reqUrls) == resNumber {
			c.DoneCancelFunc()
		}
	}
}
