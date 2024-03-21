package engine

import (
	"context"
	"zfofa/backend/engine/runner/auth"
	"zfofa/backend/engine/runner/unauth"
	"zfofa/backend/engine/scheduler"
	"zfofa/backend/engine/types"
)

type Engine struct {
	SearchType *types.FofaSearchType
}

func (e *Engine) Run(ctx context.Context, cancelFunc context.CancelFunc) []types.Asset {
	switch e.SearchType.Mode {
	case types.UnAuthMode:
		runner := unauth.ConcurrentRunner{
			WorkerCount:    e.SearchType.WorkCount,
			Scheduler:      scheduler.CreateSimpleScheduler(e.SearchType.WorkCount),
			DoneCancelFunc: cancelFunc,
		}
		return runner.Run(ctx, *e.SearchType)
	case types.AuthMode:
		runner := auth.ConcurrentRunner{
			WorkerCount:    e.SearchType.WorkCount,
			Scheduler:      scheduler.CreateSimpleScheduler(e.SearchType.WorkCount),
			DoneCancelFunc: cancelFunc,
		}
		return runner.Run(ctx, *e.SearchType)
	}

	return nil
}
