package types

import (
	"net/http"
	"net/url"
)

type Request struct {
	Url       string
	Params    url.Values
	Headers   http.Header
	ParseFunc func([]byte) ParseResult
}

type ParseResult struct {
	Assets []Asset
	Url    string
}

// NilParse 用于测试
func NilParse([]byte) ParseResult {
	return ParseResult{
		Assets: []Asset{
			{
				Host:           "https://example.com",
				Ip:             "127.0.0.1",
				Port:           80,
				Title:          "测试数据",
				Server:         "Nginx/1.13.3",
				Country:        "本地数据",
				Organization:   "Test Org",
				Header:         "HTTP/1.1 200 OK...",
				LastUpdateTime: "2024-01-25",
				Alive:          false,
			},
		},
	}
}
