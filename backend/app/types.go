package app

import "fmt"

type QuerySetting struct {
	Keywords     string `json:"keywords"`
	Cookie       string `json:"cookie"`
	Threads      int    `json:"threads"`
	AssetsNumber int    `json:"assetsNumber"`
	CheckAlive   bool   `json:"checkAlive"`
}

func (q QuerySetting) String() string {
	return fmt.Sprintf(`{Keywords: "%s", Cookie: "%s", Threads: %d, AssetsNumber: %d, CheckAlive: %v}`+"\n",
		q.Keywords, q.Cookie, q.Threads, q.AssetsNumber, q.CheckAlive)
}
