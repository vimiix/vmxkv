package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// 优先级： 命令行 > 环境变量 > config.json

type Config struct {
	DBFile   string `json:"db_file" flag:"db-file" cfg:"db_file"`
	LogDir   string `json:"log_dir" flag:"log-dir" cfg:"log_dir"`
	Socket   string `json:"socket" flag:"socket" cfg:"socket"`
	PidFile  string `json:"pid_file" flag:"pid-file" cfg:"pid_file"`
	Port     int64  `json:"port" flag:"port" cfg:"port"`
	ReadOnly bool   `json:"read_only" flag:"ro" cfg:"read_only"`
}

func NewConfig() (cfg *Config) {
	return &Config{
		DBFile:  "vmxkv.db",
		LogDir:  "/var/log/vmxkv",
		Socket:  "/var/lib/vmxkv/vmxkv.sock",
		PidFile: "/var/run/vmxkvd/vmxkvd.pid",
		Port:    19323,
	}
}

func ParseConfig(filepath string, cfg interface{}) {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Fail to read config file :%v\n", err)
	}
	if err = json.Unmarshal(bs, cfg); err != nil {
		log.Fatalf("Fail to unmarshal config file :%v\n", err)
	}
}
