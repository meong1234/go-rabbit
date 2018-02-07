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

	osSignals := make(chan os.Signal)
	signal.Notify(osSignals, os.Interrupt)

	select {
	case <-osSignals:
		util.Log.Infof("osSignal Interrupt trigerred")
		return daemon.Stop()
	}
}
