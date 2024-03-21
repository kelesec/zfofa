package fetch

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	url2 "net/url"
	"time"
)

type ProtocolType string

var (
	HttpProtocol   ProtocolType = "http"
	HttpsProtocol  ProtocolType = "https"
	Socks4Protocol ProtocolType = "socks4"
	Socks5Protocol ProtocolType = "socks5"
)

// Proxy 代理设置，暂时只支持HTTP代理
type Proxy struct {
	Protocol ProtocolType
	Host     string
	Port     int
	Username string
	Password string
}

/*
----------------------------------------------------------------------------------------------------------
request
----------------------------------------------------------------------------------------------------------
*/

type Request struct {
	client         *http.Client
	allowRedirects bool
}

func CreateRequest(AllowRedirects bool, timeout time.Duration, proxy *Proxy) *Request {
	var client *http.Client
	if AllowRedirects {
		client = &http.Client{
			Timeout: timeout,
		}
	} else {
		client = &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	// 禁用SSL
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 50,
	}

	// 配置代理
	if proxy != nil {
		u := fmt.Sprintf("%s://%s:%d", proxy.Protocol, proxy.Host, proxy.Port)
		proxyUrl, err := url2.Parse(u)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	client.Transport = transport
	return &Request{client: client, allowRedirects: AllowRedirects}
}

// AllowRedirects 是否允许重定向
func (r *Request) AllowRedirects() {
	if !r.allowRedirects {
		r.allowRedirects = true
		r.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

// DisableRedirects 禁止重定向
func (r *Request) DisableRedirects() {
	if r.allowRedirects {
		r.allowRedirects = false
		r.client.CheckRedirect = http.Client{}.CheckRedirect
	}
}

// SetProxy 设置请求代理
func (r *Request) SetProxy(proxy *Proxy) {
	u := fmt.Sprintf("%s://%s:%d", proxy.Protocol, proxy.Host, proxy.Port)
	proxyUrl, err := url2.Parse(u)
	if err != nil {
		return
	}

	r.client.Transport = &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 50,
		Proxy:               http.ProxyURL(proxyUrl),
	}
}

// AddCookie 添加Cookie
func (r *Request) AddCookie(url string, cookies []*http.Cookie) {
	u, _ := url2.Parse(url)
	r.client.Jar.SetCookies(u, cookies)
}

// Get 发送GET请求
func (r *Request) Get(url string, params url2.Values, headers http.Header) (*Response, error) {
	reqUrl, err := url2.Parse(url)
	if err != nil {
		return nil, err
	}

	url = fmt.Sprintf("%s://%s%s", reqUrl.Scheme, reqUrl.Host, reqUrl.Path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// 添加请求参数、请求头
	if params == nil {
		params = url2.Values{}
	}

	if reqUrl.RawQuery != "" {
		values, _ := url2.ParseQuery(reqUrl.RawQuery)
		for k, vs := range values {
			params[k] = append(params[k], vs...)
		}
	}

	req.URL.RawQuery = params.Encode()
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	for key, values := range headers {
		for _, v := range values {
			req.Header.Set(key, v)
		}
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}, nil
}

/*
----------------------------------------------------------------------------------------------------------
session
----------------------------------------------------------------------------------------------------------
*/

func CreateSession(AllowRedirects bool, timeout time.Duration, proxy *Proxy) *Request {
	var client *http.Client
	jar, _ := cookiejar.New(nil)

	if AllowRedirects {
		client = &http.Client{
			Timeout: timeout,
			Jar:     jar,
		}
	} else {
		client = &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Jar: jar,
		}
	}

	// 禁用SSL
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 50,
	}

	// 配置代理
	if proxy != nil {
		u := fmt.Sprintf("%s://%s:%d", proxy.Protocol, proxy.Host, proxy.Port)
		proxyUrl, err := url2.Parse(u)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	client.Transport = transport
	return &Request{client: client, allowRedirects: AllowRedirects}
}
