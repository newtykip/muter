package main

import "github.com/kardianos/service"

func (m *muter) Start(_ service.Service) error {
	if service.Interactive() {
		m.logger.Info("Running in terminal")
	} else {
		m.logger.Info("Running under service manager")
	}

	// start should not block, work is done async
	go m.run()
	return nil
}

func (m *muter) Stop(_ service.Service) error {
	m.logger.Info("muter shutting down!")
	close(m.exitChan)
	return nil
}

func (m *muter) run() {
	// todo: do work here
}
