package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "sync"

    "github.com/mreiferson/go-options"

    "github.com/vimiix/vmxkv/conf"
    "github.com/vimiix/vmxkv/internal/svc"
    "github.com/vimiix/vmxkv/server"
)

func main() {
	s := &service{}
	if err := svc.Run(s); err != nil {
		log.Fatalln("Fatal: " + err.Error())
	}
}

type service struct {
	once   sync.Once
	cfg    *conf.Config
	server *server.Server
}

type config map[string]interface{}

func (s *service) Init() (err error) {
    defaultCfg := conf.NewConfig()
    fs := conf.NewFlagSet(defaultCfg)
    fs.Parse(os.Args[1:])
    if fs.Lookup("v").Value.(flag.Getter).Get().(bool) {
        fmt.Println(Version)
        os.Exit(0)
    }
    var cfgMap config
    configFile := fs.Lookup("c").Value.String()
    if configFile != "" {
        conf.ParseConfig(configFile, &cfgMap)
    }
    fmt.Printf("cfgMap: %+v\n", cfgMap)
    options.Resolve(defaultCfg, fs, cfgMap)
    s.cfg = defaultCfg
    fmt.Printf("final cfg: %+v\n", s.cfg)
	return
}

func (s *service) Start() (err error) {
	return
}

func (s *service) Stop() (err error) {
	return
}
