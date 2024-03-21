package parser

import (
	"encoding/json"
	"fmt"
	"zfofa/backend/engine/types"
)

type respAsset struct {
	Icon    string `json:"favicon_base64"`
	Host    string `json:"host"`
	Ip      string `json:"ip"`
	Port    int    `json:"port"`
	Title   string `json:"title"`
	Server  string `json:"server"`
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
	AsnOrg  string `json:"asn_org"`
	Header  string `json:"header"`
	Cert    string `json:"cert"`
	Mtime   string `json:"mtime"`
}

type data struct {
	Assets []respAsset `json:"assets"`
}

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    data   `json:"data"`
}

func UnAuthAssetsParser(body []byte) types.ParseResult {
	var resp response

	err := json.Unmarshal(body, &resp)
	if err != nil || resp.Code != 0 || len(resp.Data.Assets) == 0 {
		return types.ParseResult{}
	}

	var assets []types.Asset
	for _, data := range resp.Data.Assets {
		assets = append(assets, types.Asset{
			Icon:           data.Icon,
			Host:           data.Host,
			Ip:             data.Ip,
			Port:           data.Port,
			Title:          data.Title,
			Server:         data.Server,
			Country:        fmt.Sprintf(`%s / %s / %s`, data.Country, data.Region, data.City),
			Organization:   data.AsnOrg,
			Header:         data.Header,
			Certificate:    data.Cert,
			LastUpdateTime: data.Mtime,
		})
	}

	return types.ParseResult{Assets: assets}
}
