package types

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"zfofa/backend/core/output"
)

/*
----------------------------------------------------------------------------------------------------------
Asset: 资产信息
----------------------------------------------------------------------------------------------------------
*/

type Asset struct {
	Icon           string `json:"icon" desc:"ICON"`
	Host           string `json:"host" desc:"主机名"`
	Ip             string `json:"ip" desc:"IP地址"`
	Port           int    `json:"port" desc:"端口"`
	Title          string `json:"title" desc:"站点标题"`
	Server         string `json:"server" desc:"服务"`
	Country        string `json:"country" desc:"国家/城市"`
	Organization   string `json:"organization" desc:"组织"`
	Header         string `json:"header" desc:"响应头"`
	Certificate    string `json:"certificate" desc:"证书"`
	Alive          bool   `json:"alive" desc:"是否存活"`
	LastUpdateTime string `json:"lastUpdateTime" desc:"最后更新时间"`
}

// GetDesc 通过结构体字段名，获取字段描述信息，不区分大小写
func (a Asset) GetDesc(field string) (string, bool) {
	assetType := reflect.TypeOf(a)

	// 实现不区分大小写比较
	for i := 0; i < assetType.NumField(); i++ {
		if strings.EqualFold(assetType.Field(i).Name, field) {
			field = assetType.Field(i).Name
		}
	}

	if sf, ok := assetType.FieldByName(field); ok {
		return sf.Tag.Get("desc"), true
	}

	return "", false
}

func (a Asset) GetDescString(field string) string {
	if s, ok := a.GetDesc(field); ok {
		return s
	}
	return ""
}

// GetDescArr 获取全部字段描述
func (a Asset) GetDescArr() []string {
	return []string{
		a.GetDescString("Icon"), a.GetDescString("Host"),
		a.GetDescString("Ip"), a.GetDescString("Port"),
		a.GetDescString("Title"), a.GetDescString("Server"),
		a.GetDescString("Country"), a.GetDescString("Organization"),
		a.GetDescString("Header"), a.GetDescString("Certificate"),
		a.GetDescString("Alive"), a.GetDescString("LastUpdateTime"),
	}
}

// GetFields 获取全部字段值
func (a Asset) GetFields() []string {
	port := strconv.Itoa(a.Port)
	alive := strconv.FormatBool(a.Alive)
	return []string{
		a.Icon, a.Host, a.Ip, port, a.Title, a.Server, a.Country,
		a.Organization, a.Header, a.Certificate, alive, a.LastUpdateTime,
	}
}

func (a Asset) Compare(asset Asset) bool {
	if strings.EqualFold(a.Host, asset.Host) &&
		strings.EqualFold(a.Ip, asset.Ip) &&
		a.Port == asset.Port {
		return true
	}

	return false
}

func (a Asset) String() string {
	var header, cert string
	if len(a.Header) > 8 {
		header = a.Header[:8]
	}
	if len(a.Certificate) > 8 {
		cert = a.Certificate[:8]
	}
	return fmt.Sprintf("{%s %s %d %s %s %s %s %s... %s... %v %s}",
		a.Host, a.Ip, a.Port, a.Title, a.Server, a.Country, a.Organization,
		header, cert, a.Alive, a.LastUpdateTime)
}

/*
----------------------------------------------------------------------------------------------------------
 保存资产数据
----------------------------------------------------------------------------------------------------------
*/

func SaveAssets(outfile string, assets []Asset) {
	index := strings.LastIndex(outfile, ".") + 1
	suffix := outfile[index:]

	if len(assets) == 0 {
		return
	}

	switch suffix {
	case "txt":
		tw, _ := output.CreateTxtWriter(outfile, false)
		defer tw.Close()

		for _, asset := range assets {
			tw.WriteStringLn(asset.Host)
		}

	case "csv":
		cw, _ := output.CreateCsvWriter(outfile)
		defer cw.Close()

		cw.WriteRow(assets[0].GetDescArr())
		for _, asset := range assets {
			cw.WriteRow(asset.GetFields())
		}

	case "json":
		jw, _ := output.CreateJsonWriter(outfile)
		defer jw.Close()

		jw.Write(assets)
	}
}
