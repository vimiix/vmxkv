package server

import (
	"fmt"

	"github.com/vimiix/vmxkv/conf"
	"github.com/vimiix/vmxkv/internal/logging"
)

const (
	dbFileMode = 0666
)

type Server struct {
	cfg    *conf.Config
	db     *DB
	stopCh chan int
}

func New(cfg *conf.Config) (s *Server) {
	if cfg == nil {
		cfg = conf.NewConfig()
	}
	db, err := OpenDB(cfg.DBFile, dbFileMode)
	if err != nil {
		logging.Fatal(fmt.Sprintf("failed opening db [%s]: %s\n", cfg.DBFile, err))
	}
	s = &Server{
		cfg:    cfg,
		db:     db,
		stopCh: make(chan int),
	}
	return
}

func (s *Server) Serve() (err error) {
	// TODO http api and cli
	//s.db.Put(1, 1)
	//s.db.Put(2, 3)
	fmt.Println(s.db.Get(2))
	<-s.stopCh
	return
}

func (s *Server) Stop() (err error) {
	s.db.Close()
	close(s.stopCh)
	return
}
