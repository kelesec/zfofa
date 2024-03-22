package check_alive

import (
	"context"
	"strings"
	"zfofa/backend/core/filelog"
	"zfofa/backend/engine/types"
)

type CheckAliveWorker struct {
	AssetInChan chan types.Asset
}

func (c *CheckAliveWorker) Worker(ctx context.Context, assetOutChan chan types.Asset) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			asset, ok := <-c.AssetInChan
			if !ok {
				continue
			}

			// 存活探测
			if checkAliveWithIcmpEcho(asset.Ip) {
				filelog.Info("[ICMP] %s is alive", asset.Ip)
				asset.Alive = true
			} else if checkAliveWithTcp(asset.Ip, asset.Port) {
				filelog.Info("[TCP] %s:%d is alive", asset.Ip, asset.Port)
				asset.Alive = true
			} else if checkAliveWithUdp(asset.Ip, asset.Port) {
				filelog.Info("[UDP] %s:%d is alive", asset.Ip, asset.Port)
				asset.Alive = true
			} else if strings.EqualFold(asset.Host[0:4], "http") &&
				checkAliveWithHttp(asset.Host) {
				filelog.Info("[HTTP] %s is alive", asset.Host)
				asset.Alive = true
			} else {
				filelog.Error("[ICMP/HTTP/UDP/TCP] %s is not alive", asset.Ip)
			}

			assetOutChan <- asset
		}
	}
}
