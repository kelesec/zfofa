package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	app2 "zfofa/backend/app"
	"zfofa/backend/core/filelog"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := app2.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "zfofa v0.1 by kele",
		Width:  900,
		Height: 620,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			// 设置最小尺寸
			runtime.WindowSetMinSize(ctx, 900, 620)

			// 传递上下文
			app.Startup(ctx)
			filelog.FileContext = ctx
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		fmt.Printf("Error: %s", err)
		filelog.Fatalf("Error: %s", err)
	}
}
