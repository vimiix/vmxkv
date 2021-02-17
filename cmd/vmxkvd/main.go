package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/mreiferson/go-options"

	"github.com/vimiix/vmxkv/conf"
	"github.com/vimiix/vmxkv/internal/logging"
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

func (cfg config) Validate() {
	if v, exists := cfg["port"]; exists {
		if intVal, err := strconv.Atoi(fmt.Sprintf("%v", v)); err == nil {
			cfg["port"] = int64(intVal)
		} else {
			log.Fatalf("failed parsing port: %+v\n", v)
		}
	}
}

func (s *service) Init() (err error) {
	cfg := conf.NewConfig()
	fs := conf.NewFlagSet(cfg)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatalf("failed parsing args: %+v\n", err)
	}
	if fs.Lookup("v").Value.(flag.Getter).Get().(bool) {
		fmt.Println(Version)
		os.Exit(0)
	}
	var cfgMap config
	configFile := fs.Lookup("c").Value.String()
	if configFile != "" {
		conf.ParseConfig(configFile, &cfgMap)
	}
	cfgMap.Validate()
	options.Resolve(cfg, fs, cfgMap)
	s.cfg = cfg
	logging.Init(cfg.LogDir)
	s.server = server.New(cfg)
	return
}

func (s *service) Start() (err error) {
	err = s.server.Serve()
	return
}

func (s *service) Stop() (err error) {
	err = s.server.Stop()
	return
}
