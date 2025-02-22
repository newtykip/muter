package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/kardianos/service"
	"github.com/moutend/go-wca/pkg/wca"
)

const (
	Port     = 3034
	PollRate = 500 * time.Second
)

var mute bool

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
	return nil
}

func trackMute(aev *wca.IAudioEndpointVolume, logger service.Logger) {
	ticker := time.NewTicker(PollRate)
	for range ticker.C {
		// get current mute state
		var currentMute bool
		if err := aev.GetMute(&currentMute); err != nil {
			logger.Error(err)
		}

		// if mute state has changed, update and broadcast
		if currentMute != mute {
			mute = currentMute
			logger.Infof("Mute state changed to %v", mute)
		}
	}
}

// todo: detect updated audio devices
func (m *muter) run() {
	// initialize COM
	if err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED); err != nil {
		log.Fatal(err)
	}
	defer ole.CoUninitialize()

	// get audio device enumerator
	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		log.Fatal(err)
	}
	defer mmde.Release()

	// get default capture device
	var mmd *wca.IMMDevice
	if err := mmde.GetDefaultAudioEndpoint(wca.ECapture, wca.EMultimedia, &mmd); err != nil {
		log.Fatal(err)
	}
	defer mmd.Release()

	// activate audio endpoint volume
	var aev *wca.IAudioEndpointVolume
	if err := mmd.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev); err != nil {
		log.Fatal(err)
	}
	defer aev.Release()

	// track current mute state
	go trackMute(aev, m.logger)

	// start websocket server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleWebsocket(w, r, m.logger)
	})
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", Port), nil)
}
