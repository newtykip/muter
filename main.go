package main

import (
	"log"

	"github.com/kardianos/service"
)

const (
	ServiceName string = "dev.newty.muter"
)

type muter struct {
	exitChan chan struct{}
	logger   service.Logger
}

func main() {
	// create service
	m := &muter{
		exitChan: make(chan struct{}),
	}
	svc, err := service.New(m, &service.Config{
		Name: ServiceName,
	})
	if err != nil {
		log.Fatal(err)
	}

	// collect errors
	errorChan := make(chan error, 5)
	m.logger, err = svc.Logger(errorChan)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errorChan
			if err != nil {
				log.Print(err)
			}
		}
	}()

	// start service
	if err := svc.Run(); err != nil {
		m.logger.Error(err)
	}
}
