package types

import (
	"fmt"
	"regexp"
	"time"
	"zfofa/backend/core/conf"
	"zfofa/backend/core/fetch"
)

type SearchMode string

var (
	UnAuthMode SearchMode = "UnAuthMode" // 未登录模式
	AuthMode   SearchMode = "AuthMode"   // 登陆模式
)

/*
----------------------------------------------------------------------------------------------------------
FofaSearchType: fofa的查询规则
	Url: https://fofa.info、https://api.fofa.info
	Keywords: 查询关键字（无需编码）
	SearchCount: 需要获取的资产数量
----------------------------------------------------------------------------------------------------------
*/

type FofaSearchType struct {
	Url         string
	Keywords    string
	SearchCount int
	WorkCount   int
	Mode        SearchMode
	Proxy       *fetch.Proxy
	FofaConf    *conf.Config
}

/*
----------------------------------------------------------------------------------------------------------
LastUpdateTime: Fofa语法时间类
----------------------------------------------------------------------------------------------------------
*/

type lastUpdateTime struct {
	days       int
	afterTime  time.Time
	beforeTime time.Time
}

func CreateLastUpdateTime(days int, afterTime time.Time, beforeTime time.Time) *lastUpdateTime {
	return &lastUpdateTime{days: days, afterTime: afterTime, beforeTime: beforeTime}
}

// CreateLastUpdateTimeWithCurrentTime
// afterTime 当前时间前 days 天
// beforeTime 当前时间后 1 天
func CreateLastUpdateTimeWithCurrentTime(days int) *lastUpdateTime {
	if days > 0 {
		days = -days
	}

	return &lastUpdateTime{
		days:       days,
		afterTime:  time.Now().AddDate(0, 0, days+1),
		beforeTime: time.Now().AddDate(0, 0, 1),
	}
}

func (l *lastUpdateTime) SubDays() {
	l.afterTime = l.afterTime.AddDate(0, 0, l.days)
	l.beforeTime = l.beforeTime.AddDate(0, 0, l.days)
}

func (l *lastUpdateTime) GetAfterTime() string {
	return l.afterTime.Format("2006-01-02")
}

func (l *lastUpdateTime) GetBeforeTime() string {
	return l.beforeTime.Format("2006-01-02")
}

var (
	afterReg  = regexp.MustCompile(`after\s*?=\s*?['|"]([^'|^"]+)['|"]`)
	beforeReg = regexp.MustCompile(`before\s*?=\s*?['|"]([^'|^"]+)['|"]`)
)

// GenKeywordsWithUpTime 给关键字添加 after 和 before 关键字
func (l *lastUpdateTime) GenKeywordsWithUpTime(keywords string) string {
	if !afterReg.MatchString(keywords) {
		keywords = fmt.Sprintf(`%s && after="%s"`, keywords, l.GetAfterTime())
	}

	if !beforeReg.MatchString(keywords) {
		keywords = fmt.Sprintf(`%s && before="%s"`, keywords, l.GetBeforeTime())
	}

	return keywords
}

// ReplaceKeywordsWithUpTime 替换 after、before 时间
func (l *lastUpdateTime) ReplaceKeywordsWithUpTime(keywords string) string {
	afKeywords := afterReg.ReplaceAllString(keywords, fmt.Sprintf(`after="%s"`, l.GetAfterTime()))
	bfKeywords := beforeReg.ReplaceAllString(afKeywords, fmt.Sprintf(`before="%s"`, l.GetBeforeTime()))
	return bfKeywords
}

// RemoveAfterTime 移除 after 语法
func (l *lastUpdateTime) RemoveAfterTime(keywords string) string {
	reg := regexp.MustCompile(`&{2}\s*after\s*?=\s*?['|"][^'|^"]+['|"]`)
	return reg.ReplaceAllString(keywords, "")
}

// RemoveBeforeTime 移除 before 语法
func (l *lastUpdateTime) RemoveBeforeTime(keywords string) string {
	reg := regexp.MustCompile(`&{2}\s*before\s*?=\s*?['|"][^'|^"]+['|"]`)
	return reg.ReplaceAllString(keywords, "")
}
