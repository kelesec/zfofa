package scheduler

import (
	"zfofa/backend/engine/types"
)

type SimpleScheduler struct {
	RequestInChan chan types.Request
	chanClosed    bool
	cap           int
}

func CreateSimpleScheduler(cap int) *SimpleScheduler {
	inChan := make(chan types.Request, cap)
	return &SimpleScheduler{RequestInChan: inChan, cap: cap, chanClosed: false}
}

func (s *SimpleScheduler) Configure(inChan chan types.Request) bool {
	if inChan != nil {
		s.RequestInChan = inChan
		return true
	}
	return false
}

func (s *SimpleScheduler) Close() {
	if !s.chanClosed {
		s.chanClosed = true
		close(s.RequestInChan)
	}
}

func (s *SimpleScheduler) IsClosed() bool {
	return s.chanClosed
}

func (s *SimpleScheduler) Push(r types.Request) {
	go func() {
		if !s.chanClosed {
			s.RequestInChan <- r
		}
	}()
}

func (s *SimpleScheduler) PushLastRequest(r types.Request) {
	go func() {
		if !s.chanClosed {
			s.RequestInChan <- r
			s.Close()
		}
	}()
}

func (s *SimpleScheduler) Pop() *types.Request {
	if r, ok := <-s.RequestInChan; ok {
		return &r
	}
	return nil
}

func (s *SimpleScheduler) GetRequestInChan() chan types.Request {
	return s.RequestInChan
}
