package worker

import (
	"context"
	"log"
	"net"
	"sync"
	"zfofa/backend/core/fetch"
	"zfofa/backend/core/filelog"
	"zfofa/backend/engine/types"
)

type SimpleWorker struct {
	Request            *fetch.Request
	activeWorkerCount  int
	Locker             *sync.Mutex
	ParseResultOutChan chan types.ParseResult
}

// setActiveStatus 激活，设置为工作状态
func (s *SimpleWorker) setActiveStatus() {
	s.Locker.Lock()
	s.activeWorkerCount++
	s.Locker.Unlock()
}

// delActiveStatus 销毁，当前工作已经完成
func (s *SimpleWorker) delActiveStatus() {
	s.Locker.Lock()
	s.activeWorkerCount--
	s.Locker.Unlock()
}

// IsActive 判断是否处于工作状态
func (s *SimpleWorker) IsActive() bool {
	return s.activeWorkerCount != 0
}

func (s *SimpleWorker) Worker(ctx context.Context, scheduler *types.Scheduler) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			r := (*scheduler).Pop()
			if r == nil {
				continue
			}

			// 置为工作状态
			s.setActiveStatus()

			resp, err := s.Request.Get(r.Url, r.Params, r.Headers)
			if err != nil {
				log.Printf("fetching fail: (url=%s), %s", r.Url, err)
				filelog.Error("fetching fail: (url=%s), %s", r.Url, err)

				// 遇到请求超时错误，将结果添加到队列中再次请求
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					(*scheduler).Push(*r)
				}

				// 请求执行完成，工作结束
				s.delActiveStatus()
				continue
			}

			if resp.StatusCode != 200 {
				log.Printf("[%d] err status code: %s", resp.StatusCode, r.Url)
				filelog.Error("[%d] err status code: %s", resp.StatusCode, r.Url)

				(*scheduler).Push(*r)

				s.delActiveStatus()
				continue
			}

			// 请求成功，解析内容
			log.Printf("[200] fetched url: %s", r.Url)
			filelog.Info("[200] fetched url: %s", r.Url)
			res := r.ParseFunc(resp.Body)
			res.Url = r.Url
			s.ParseResultOutChan <- res

			s.delActiveStatus()
		}
	}
}
