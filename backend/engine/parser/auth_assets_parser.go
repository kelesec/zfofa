package parser

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"zfofa/backend/engine/types"
)

var (
	//matchItemsReg = regexp.MustCompile(`hsxa-meta-data-item">[\s\S]+?el-scrollbar__bar`)
	matchItemsReg = regexp.MustCompile(`hsxa-meta-data-item">[\s\S]+?class="el-scrollbar__bar is-vertical"[\s\S]+?<!----></div><!----></div>`)

	iconReg        = regexp.MustCompile(`class="el-image hsxa-favicon".*?src="([^"]+)"`)
	hostReg        = regexp.MustCompile(`class="hsxa-host".*?href="([^"]+)"`)
	ipReg          = regexp.MustCompile(`class="hsxa-meta-data-list-main-left hsxa-fl"[\s\S]+?<span style="display:none;">([^<]+)</`)
	portReg        = regexp.MustCompile(`class="hsxa-port".*?>([^>]+)<`)
	titleReg       = regexp.MustCompile(`class="hsxa-meta-data-list-main-left hsxa-fl".*?<p.*?>([^<]+)</p>`)
	serverReg      = regexp.MustCompile(`class="hsxa-list-icon-show"></span>([^<]+)<`)
	countryReg     = regexp.MustCompile(`class="el-tooltip item".*?<a.*?">([^<]+).*?<a.*?">([^<]+).*?<a.*?">([^<]+)<`)
	orgReg         = regexp.MustCompile(`<span>[^A].*?class="el-tooltip hsxa-jump-a item".*?>([^<]+)<`)
	headerReg      = regexp.MustCompile(`class="el-scrollbar__view".*?<span>([^<]+)<`)
	certificateReg = regexp.MustCompile(`class="hsxa-certs-detail">([^<]+)<`)
	lastUpTimeReg  = regexp.MustCompile(`class="el-tooltip hsxa-jump-a item"[\s\S]+?<p>.*?<span>([^<]+)<`)
)

func AuthAssetsParser(body []byte) types.ParseResult {
	items := matchItemsReg.FindAllSubmatch(body, -1)
	var assets []types.Asset
	for _, item := range items {
		assets = append(assets, types.Asset{
			Icon:           matchIcon(item[0]),
			Host:           matchHost(item[0]),
			Ip:             matchIp(item[0]),
			Port:           matchPort(item[0]),
			Title:          matchTitle(item[0]),
			Server:         matchServer(item[0]),
			Country:        matchCountry(item[0]),
			Organization:   matchOrg(item[0]),
			Header:         matchHeader(item[0]),
			Certificate:    matchCertificate(item[0]),
			LastUpdateTime: matchLastUpTime(item[0]),
		})
	}
	return types.ParseResult{Assets: assets}
}

func matchIcon(b []byte) string {
	//match := iconReg.FindSubmatch(b)
	//if len(match) == 0 {
	//	return ""
	//}
	//return strings.TrimSpace(string(match[1]))
	return ""
}

func matchHost(b []byte) string {
	match := hostReg.FindSubmatch(b)
	if len(match) == 0 {
		return ""
	}
	return strings.TrimSpace(string(match[1]))
}

func matchIp(b []byte) string {
	match := ipReg.FindSubmatch(b)
	if len(match) == 0 {
		return ""
	}
	return strings.TrimSpace(string(match[1]))
}

func matchPort(b []byte) int {
	match := portReg.FindSubmatch(b)
	if len(match) == 0 {
		return 0
	}

	port, err := strconv.Atoi(strings.TrimSpace(string(match[1])))
	if err != nil {
		return 0
	}
	return port
}

func matchTitle(b []byte) string {
	match := titleReg.FindSubmatch(b)
	if len(match) == 0 {
		return ""
	}
	return strings.TrimSpace(string(match[1]))
}

func matchServer(b []byte) string {
	match := serverReg.FindSubmatch(b)
	if len(match) == 0 {
		return ""
	}
	return strings.TrimSpace(string(match[1]))
}

func matchCountry(b []byte) string {
	match := countryReg.FindSubmatch(b)
	if len(match) < 3 {
		return ""
	}
	return fmt.Sprintf("%s / %s / %s",
		strings.TrimSpace(string(match[1])),
		strings.TrimSpace(string(match[2])),
		strings.TrimSpace(string(match[3])))
}

func matchOrg(b []byte) string {
	match := orgReg.FindSubmatch(b)
	if len(match) == 0 {
		return ""
	}
	return strings.TrimSpace(string(match[1]))
}

func matchHeader(b []byte) string {
	match := headerReg.FindSubmatch(b)
	if len(match) == 0 {
		return ""
	}
	return html.UnescapeString(strings.TrimSpace(string(match[1])))
}

func matchCertificate(b []byte) string {
	match := certificateReg.FindSubmatch(b)
	if len(match) == 0 {
		return ""
	}
	return strings.TrimSpace(string(match[1]))
}

func matchLastUpTime(b []byte) string {
	match := lastUpTimeReg.FindSubmatch(b)
	if len(match) == 0 {
		return ""
	}
	return strings.TrimSpace(string(match[1]))
}
