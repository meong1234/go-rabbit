package application

import (
	"github.com/go-rabbit/util"
	"os"
	"os/signal"
)

func AppRunner(daemon util.Daemon) error {
	err := daemon.Start()
	if err != nil {
		return err
	}

	serverErrors := make(chan error, 1)
	osSignals := make(chan os.Signal)
	signal.Notify(osSignals, os.Interrupt)

	select {
	case err := <-serverErrors:
		return err
	case <-osSignals:
		util.Log.Infof("osSignal Interrupt trigerred")
		return daemon.Stop()
	}
}
