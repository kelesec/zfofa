package check_alive

import (
	"context"
	"zfofa/backend/engine/types"
)

type CheckAlive struct {
	MaxCheckAliveWorkers int
}

func (c *CheckAlive) Run(ctx context.Context, assets []types.Asset, assetOutChan chan types.Asset) {
	inChan := make(chan types.Asset, 50)
	caw := CheckAliveWorker{
		AssetInChan: inChan,
	}

	// 提交任务
	go func() {
		for _, asset := range assets {
			inChan <- asset
		}
	}()

	// 开始任务
	if c.MaxCheckAliveWorkers <= 0 {
		c.MaxCheckAliveWorkers = 100
	}

	for i := 0; i < c.MaxCheckAliveWorkers; i++ {
		go caw.Worker(ctx, assetOutChan)
	}
}
