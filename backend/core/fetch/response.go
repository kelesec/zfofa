package fetch

import (
	"net/http"
	"unsafe"
)

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

func (r *Response) BodyToString() (string, error) {
	return *(*string)(unsafe.Pointer(&r.Body)), nil
}
