package conf

import (
    "flag"
)

func NewFlagSet(cfg *Config) (fs *flag.FlagSet) {
    fs = flag.NewFlagSet("vmxkv", flag.ExitOnError)

    fs.Bool("v", false, "Show version")
    fs.String("c", "", "Config file path")
    fs.Int("port", cfg.Port, "Port number to use for connection")
    fs.String("data-dir", cfg.DataDir, "Dir to store data")
    fs.String("log-dir", cfg.LogDir, "Dir to store log")
    fs.String("socket", cfg.Socket, "Unix socket file")
    fs.String("pid-file", cfg.PidFile, "Process pid file path")
    return
}