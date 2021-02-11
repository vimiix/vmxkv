package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// 优先级： 命令行 > 环境变量 > config.json

type Config struct {
	DataDir string `json:"data_dir" flag:"data-dir" cfg:"data_dir"`
	LogDir  string `json:"log_dir" flag:"log-dir" cfg:"log_dir"`
	Socket  string `json:"socket" flag:"socket" cfg:"socket"`
	PidFile string `json:"pid_file" flag:"pid-file" cfg:"pid_file"`
	Port    int    `json:"port" flag:"port" cfg:"port"`
}

func NewConfig() (cfg *Config) {
	return &Config{
		DataDir: "/var/lib/vmxkv",
		LogDir:  "/var/log/vmxkv",
		Socket:  "/var/lib/vmxkv/vmxkv.sock",
		PidFile: "/var/run/vmxkvd/vmxkvd.pid",
		Port:    19323,
	}
}

func ParseConfig(filepath string, cfg interface{}) {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Fail to read config file :%v", err)
	}
	if err = json.Unmarshal(bs, cfg); err != nil {
		log.Fatalf("Fail to unmarshal config file :%v", err)
	}
}
