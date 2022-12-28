package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"tcp-proxy/modules/log"
)

type ClientConfig struct {
	ServerAddr string  `json:"server_addr"`
	Proxy      []Proxy `json:"proxy"`
	Log        bool    `json:"log"`
	LogPath    string  `json:"log_path"`
}

type Proxy struct {
	Name       string `json:"name"`
	ProxyAddr  string `json:"proxy_addr"`  // 本地需要代理的地址
	RemotePort uint16 `json:"remote_port"` // 远程代理端口
}

var Config ClientConfig

func exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

func init() {
	var config string

	flag.StringVar(&config, "c", "", "选择配置文件.")
	flag.Parse()
	exePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panicf("获取执行路径失败：%s", err.Error())
	}
	if config == "" && exists(path.Join(exePath, "config.json")) {
		config = path.Join(exePath, "config.json")
	}

	if config != "" {
		bs, err := ioutil.ReadFile(config)
		if err != nil {
			log.Panicf("配置读取文件失败：%s", err.Error())
		}
		err = json.Unmarshal(bs, &Config)
		if err != nil {
			log.Panicf("配置文件解析失败：%s", err.Error())
		}
	} else {
		log.Println("请选择配置文件")
	}

	if Config.Log {
		if Config.LogPath != "" {
			log.SaveLog(Config.LogPath)
		} else {
			log.SaveLog(exePath)
		}
	}

}
