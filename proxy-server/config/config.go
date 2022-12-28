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

type Info struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Log  bool   `json:"log"`
}

var Config Info

func exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

func init() {
	var config string
	var lPort int

	flag.StringVar(&config, "c", "", "选择配置文件.")
	flag.IntVar(&lPort, "l", 5000, "指定端口.")
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
	} else if lPort != 0 {
		Config = Info{
			Host: "0.0.0.0",
			Port: lPort,
		}
	} else {
		Config = Info{
			Host: "0.0.0.0",
			Port: 5000,
		}

	}

	if Config.Log {
		log.SaveLog(exePath)
	}
}
