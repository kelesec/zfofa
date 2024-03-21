package types

type Scheduler interface {
	Configure(chan Request) bool
	Push(Request)
	Pop() *Request
	Close()
	IsClosed() bool
	GetRequestInChan() chan Request
}
