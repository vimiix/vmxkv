package conf

import (
	"flag"
)

func NewFlagSet(cfg *Config) (fs *flag.FlagSet) {
	fs = flag.NewFlagSet("", flag.ExitOnError)

	fs.Bool("v", false, "Show version")
	fs.String("c", "", "Config file path")
	fs.Int64("port", cfg.Port, "Port number to use for connection")
	fs.String("db-file", cfg.DBFile, "File to store data")
	fs.String("log-dir", cfg.LogDir, "Dir to store log")
	fs.String("socket", cfg.Socket, "Unix socket file")
	fs.String("pid-file", cfg.PidFile, "Process pid file path")
	fs.Bool("ro", false, "Read-only mode")
	return
}
