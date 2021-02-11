package svc

import (
	"os"
	"os/signal"
	"syscall"
)

// 创建一个变量，方便测试时 mock
var signalNotify = signal.Notify

type Service interface {
	Init() error
	Start() error
	Stop() error
}

func Run(s Service, sig ...os.Signal) error {
	if err := s.Init(); err != nil {
		return err
	}

	var errCh = make(chan error, 1)
	go func() {
		errCh <- s.Start()
	}()

	if len(sig) == 0 {
		sig = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	signalCh := make(chan os.Signal, 1)
	signalNotify(signalCh, sig...)

	select {
	case err := <-errCh:
		return err
	case <-signalCh:
	}
	return s.Stop()

}
