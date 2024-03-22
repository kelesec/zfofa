package conf

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Userinfo struct {
	FofaToken string `yaml:"fofa_token"`
}

type FofaToolConf struct {
	// unauth模式请求到空资产时，最大重复请求次数
	MaxTryFetches int `yaml:"max_try_fetches"`

	// unauth模式中 after、before 时间间隔天数
	BetweenAfterTimeAndBeforeTime int `yaml:"between_after_time_and_before_time"`

	// 设置HTTP代理
	HttpProxy string `yaml:"http_proxy"`

	// 存活检测时，创建的协程数
	MaxCheckAliveWorkers int `yaml:"max_check_alive_workers"`
}

type Config struct {
	Userinfo     Userinfo     `yaml:"userinfo"`
	FofaToolConf FofaToolConf `yaml:"fofa_tool_conf"`
}

// ExportConf 导出配置文件
func ExportConf() (bool, error) {
	conf := Config{
		Userinfo: Userinfo{
			FofaToken: "fofa_token",
		},
		FofaToolConf: FofaToolConf{
			MaxTryFetches:                 7,
			BetweenAfterTimeAndBeforeTime: 7,
			HttpProxy:                     "http://ip:port",
			MaxCheckAliveWorkers:          100,
		},
	}

	confBytes, err := yaml.Marshal(conf)
	if err != nil {
		return false, err
	}

	if err = os.WriteFile("config.yaml", confBytes, 0666); err != nil {
		return false, err
	}

	return true, nil
}

// ImportConf 导入配置文件信息
func ImportConf() (*Config, error) {
	file, err := os.Open("config.yaml")
	if os.IsNotExist(err) {
		ok, _ := ExportConf()
		if !ok {
			return nil, errors.New("error: 配置文件不存在，且创建失败")
		}
		return nil, errors.New("error: 配置文件不存在，已重新创建，修改配置后继续")
	}
	defer file.Close()

	confBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("error: 配置文件读取失败")
	}

	conf := &Config{}
	err = yaml.Unmarshal(confBytes, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
