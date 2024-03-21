package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	url2 "net/url"
	"time"
	"zfofa/backend/core/crypto"
	"zfofa/backend/core/fetch"
	"zfofa/backend/engine/types"
)

const (
	appId      = "9e9fb94330d97833acfbc041ee1a76793f1bc691"
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAv0xjefuBTF6Ox940ZqLLUFFBDtTcB9dAfDjWgyZ2A55K+VdG
c1L5LqJWuyRkhYGFTlI4K5hRiExvjXuwIEed1norp5cKdeTLJwmvPyFgaEh7Ow19
Tu9sTR5hHxThjT8ieArB2kNAdp8Xoo7O8KihmBmtbJ1umRv2XxG+mm2ByPZFlTdW
RFU38oCPkGKlrl/RzOJKRYMv10s1MWBPY6oYkRiOX/EsAUVae6zKRqNR2Q4HzJV8
gOYMPvqkau8hwN8i6r0z0jkDGCRJSW9djWk3Byi3R2oSdZ0IoS+91MFtKvWYdnNH
2Ubhlnu1P+wbeuIFdp2u7ZQOtgPX0mtQ263e5QIDAQABAoIBAD67GwfeTMkxXNr3
5/EcQ1XEP3RQoxLDKHdT4CxDyYFoQCfB0e1xcRs0ywI1be1FyuQjHB5Xpazve8lG
nTwIoB68E2KyqhB9BY14pIosNMQduKNlygi/hKFJbAnYPBqocHIy/NzJHvOHOiXp
dL0AX3VUPkWW3rTAsar9U6aqcFvorMJQ2NPjijcXA0p1MlZAZKODO2wqidfQ487h
xy0ZkriYVi419j83a1cCK0QocXiUUeQM6zRNgQv7LCmrFo2X4JEzlujEveqvsDC4
MBRgkK2lNH+AFuRwOEr4PIlk9rrpHA4O1V13P3hJpH5gxs5oLLM1CWWG9YWLL44G
zD9Tm8ECgYEA8NStMXyAmHLYmd2h0u5jpNGbegf96z9s/RnCVbNHmIqh/pbXizcv
mMeLR7a0BLs9eiCpjNf9hob/JCJTms6SmqJ5NyRMJtZghF6YJuCSO1MTxkI/6RUw
mrygQTiF8RyVUlEoNJyhZCVWqCYjctAisEDaBRnUTpNn0mLvEXgf1pUCgYEAy1kE
d0YqGh/z4c/D09crQMrR/lvTOD+LRMf9lH+SkScT0GzdNIT5yuscRwKsnE6SpC5G
ySJFVhCnCBsQqq+ohsrXt8a99G7ePTMSAGK3QtC7QS3liDmvPBk6mJiLrKiRAZos
vgPg7nTP8VuF0ZIKzkdWbGoMyNxVFZXovQ8BYxECgYBvCR9xGX4Qy6KiDlV18wNu
ElYkxVqFBBE0AJRg/u+bnQ9jWhi2zxLa1eWZgtss80c876I8lbkGNWedOVZioatm
MFLC4bFalqyZWyO7iP7i60LKvfDJfkOSlDUu3OikahFOiqyG1VBz4+M4U500alIU
AVKD14zTTZMopQSkgUXsoQKBgHd8RgiD3Qde0SJVv97BZzP6OWw5rqI1jHMNBK72
SzwpdxYYcd6DaHfYsNP0+VIbRUVdv9A95/oLbOpxZNi2wNL7a8gb6tAvOT1Cvggl
+UM0fWNuQZpLMvGgbXLu59u7bQFBA5tfkhLr5qgOvFIJe3n8JwcrRXndJc26OXil
0Y3RAoGAJOqYN2CD4vOs6CHdnQvyn7ICc41ila/H49fjsiJ70RUD1aD8nYuosOnj
wbG6+eWekyLZ1RVEw3eRF+aMOEFNaK6xKjXGMhuWj3A9xVw9Fauv8a2KBU42Vmcd
t4HRyaBPCQQsIoErdChZj8g7DdxWheuiKoN4gbfK4W1APCcuhUA=
-----END RSA PRIVATE KEY-----`
)

// GenerateUrlWithUnAuth 获取请求URL
func GenerateUrlWithUnAuth(keywords string, assetCount int) (url string) {
	ts := time.Now().Unix()
	keywords = base64.StdEncoding.EncodeToString([]byte(keywords))
	sign := crypto.GenSignWithRsaSha256(
		fmt.Sprintf("fullfalsepage1qbase64%ssize%dts%d", keywords, assetCount, ts), privateKey)

	// 对keywords和sign进行URL编码，修复 `[-9] API校验密匙错` 问题
	return fmt.Sprintf(
		"https://api.fofa.info/v1/search?qbase64=%s&full=false&page=1&size=%d&ts=%d&sign=%s&app_id=%s",
		url2.QueryEscape(keywords), assetCount, ts, url2.QueryEscape(sign), appId)
}

/*
----------------------------------------------------------------------------------------------------------
Data: 资产总数、独立IP数量
StatsType: 响应情况
----------------------------------------------------------------------------------------------------------
*/

type Data struct {
	Size        int `json:"size"`         // 资产总数
	DistinctIps int `json:"distinct_ips"` // 独立IP数量
}

type StatsType struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 响应信息
	Data    Data   `json:"data"`    // 资产情况
}

// GenerateUrlsWithAuth 获取请求的URL
func GenerateUrlsWithAuth(searchType types.FofaSearchType) (urls []string) {
	keywords := base64.StdEncoding.EncodeToString([]byte(searchType.Keywords))
	ts := time.Now().Unix()
	sign := crypto.GenSignWithRsaSha256(
		fmt.Sprintf("fullfalseqbase64%sts%d", keywords, ts), privateKey)

	url := fmt.Sprintf(
		"https://api.fofa.info/v1/search/stats?qbase64=%s&full=false&fields=&ts=%d&sign=%s&app_id=%s",
		keywords, ts, sign, appId)

	// 增加代理
	req := fetch.CreateRequest(true, 10*time.Second, searchType.Proxy)

	resp, err := req.Get(url, nil, nil)
	if err != nil {
		return nil
	}

	stats := StatsType{}
	err = json.Unmarshal(resp.Body, &stats)
	if err != nil || stats.Code != 0 {
		return nil
	}

	size := stats.Data.Size
	if searchType.SearchCount <= 0 || searchType.SearchCount > size {
		if size <= 10 {
			urls = append(urls, fmt.Sprintf("%s/result?qbase64=%s", searchType.Url, keywords))
			return
		}

		pageSize := size / 10
		lastPageCount := size % 10
		urls = append(urls, fmt.Sprintf("%s/result?qbase64=%s", searchType.Url, keywords))

		for i := 2; i <= pageSize; i++ {
			urls = append(urls,
				fmt.Sprintf("%s/result?qbase64=%s&page=%d&page_size=10", searchType.Url, keywords, i))
		}

		if lastPageCount != 0 {
			urls = append(urls,
				fmt.Sprintf("%s/result?qbase64=%s&page=%d&page_size=%d",
					searchType.Url, keywords, pageSize+1, lastPageCount))
		}
		return
	}

	// 不超过总数，且不是获取全部数据的情况
	pageSize := searchType.SearchCount / 10
	lastPageCount := searchType.SearchCount % 10
	urls = append(urls, fmt.Sprintf("%s/result?qbase64=%s", searchType.Url, keywords))
	if pageSize <= 0 {
		return
	}

	for i := 2; i <= pageSize; i++ {
		urls = append(urls,
			fmt.Sprintf("%s/result?qbase64=%s&page=%d&page_size=10", searchType.Url, keywords, i))
	}

	if lastPageCount == 0 {
		return
	}

	if lastPageCount < 2 {
		lastPageCount = 2
	}
	urls = append(urls, fmt.Sprintf("%s/result?qbase64=%s&page=%d&page_size=%d",
		searchType.Url, keywords, pageSize+1, lastPageCount))
	return
}
