package app

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"zfofa/backend/check_alive"
	"zfofa/backend/core/conf"
	"zfofa/backend/core/fetch"
	"zfofa/backend/core/filelog"
	"zfofa/backend/engine"
	"zfofa/backend/engine/types"
)

// App struct
type App struct {
	ctx        context.Context
	config     *conf.Config
	Assets     []types.Asset
	StopCancel context.CancelFunc
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved,
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// 检测配置文件
	config, err := conf.ImportConf()
	if err != nil {
		filelog.Fatalf(err.Error())
	}

	// 没有修改的值替换成空
	if config.FofaToolConf.HttpProxy == "http://ip:port" {
		config.FofaToolConf.HttpProxy = ""
	}

	if config.Userinfo.FofaToken == "fofa_token" {
		config.Userinfo.FofaToken = ""
	}

	a.config = config
}

// StartQuery 开始查询
func (a *App) StartQuery(setting QuerySetting) []types.Asset {
	// 每次开始前清空之前的资产和还原配置
	a.Assets = []types.Asset{}
	a.Startup(a.ctx)

	// 基本参数
	fst := types.FofaSearchType{
		Url:         "https://api.fofa.info",
		Keywords:    setting.Keywords,
		SearchCount: setting.AssetsNumber,
		WorkCount:   setting.Threads,
		Mode:        types.UnAuthMode,
		FofaConf:    a.config,
	}

	if setting.Cookie != "" {
		a.config.Userinfo.FofaToken = setting.Cookie
	}

	// 模式选择
	if a.config.Userinfo.FofaToken != "" {
		fst.Mode = types.AuthMode
		fst.Url = "https://fofa.info"
	}

	// 代理配置
	if a.config.FofaToolConf.HttpProxy != "" {
		p := a.config.FofaToolConf.HttpProxy
		u, err := url.Parse(p)
		if err != nil {
			fst.Proxy = nil
		}

		port, _ := strconv.Atoi(u.Port())
		fst.Proxy = &fetch.Proxy{
			Host: u.Hostname(),
			Port: port,
		}

		switch u.Scheme {
		case "http":
			fst.Proxy.Protocol = fetch.HttpProtocol
		case "https":
			fst.Proxy.Protocol = fetch.HttpsProtocol
		}
	}

	// 查询
	ctx, cancel := context.WithCancel(context.Background())
	a.StopCancel = cancel
	e := engine.Engine{
		SearchType: &fst,
	}
	assets := e.Run(ctx, cancel)
	assets = types.RemoveDuplicateAssets(assets)

	if setting.CheckAlive {
		filelog.Info("Start Check Alive. Assets: %d", len(assets))

		// 探测存活
		out := make(chan types.Asset, 50)
		ck := check_alive.CheckAlive{
			MaxCheckAliveWorkers: a.config.FofaToolConf.MaxCheckAliveWorkers,
		}

		ctx, cancel := context.WithCancel(context.Background())
		ck.Run(ctx, assets, out)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case asset := <-out:
					a.Assets = append(a.Assets, asset)
				}

				if len(a.Assets) == len(assets) {
					cancel()
				}
			}
		}()
		wg.Wait()
	} else {
		a.Assets = append(a.Assets, assets...)
	}

	log.Printf(
		"[OVER] all requests closed, %d assets successfully acquired",
		len(a.Assets))
	filelog.Info(
		"[OVER] all requests closed, %d assets successfully acquired",
		len(a.Assets))
	return a.Assets
}

// StopQuery 停止查询
func (a *App) StopQuery() {
	if a.StopCancel != nil {
		a.StopCancel()
	}
}

// ExportAssets 导出资产，返回导出后的文件路径
func (a *App) ExportAssets(fileType []string) {
	date := time.Now().Format("2006/01/02")
	date = strings.ReplaceAll(date, "/", "-")
	curTime := time.Now().UnixMicro()

	for _, suffix := range fileType {
		filename := fmt.Sprintf("%s-%d.%s", date, curTime, suffix)
		curPath, err := os.Getwd()
		if err != nil {
			filelog.Error("%s", err)
			continue
		}

		realFilePath := filepath.Join(curPath, filename)
		types.SaveAssets(realFilePath, a.Assets)
		filelog.Info("The results are saved in `%s`", realFilePath)
	}
}
