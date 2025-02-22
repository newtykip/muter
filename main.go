package main

import (
	"log"

	"github.com/go-ole/go-ole"
	"github.com/kardianos/service"
	"github.com/moutend/go-wca/pkg/wca"
)

const (
	ServiceName string = "dev.newty.muter"
)

type muter struct {
	exitChan chan struct{}
	logger   service.Logger
}

func main() {
	// // create service
	// m := &muter{
	// 	exitChan: make(chan struct{}),
	// }
	// svc, err := service.New(m, &service.Config{
	// 	Name: ServiceName,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // collect errors
	// errorChan := make(chan error, 5)
	// m.logger, err = svc.Logger(errorChan)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// go func() {
	// 	for {
	// 		err := <-errorChan
	// 		if err != nil {
	// 			log.Print(err)
	// 		}
	// 	}
	// }()

	// // start service
	// if err := svc.Run(); err != nil {
	// 	m.logger.Error(err)
	// }

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

	// get initial mute state
	var muted bool
	if err := aev.GetMute(&muted); err != nil {
		log.Fatal(err)
	}
	log.Printf("initial mute state: %v", muted)
}
