package gutils

import (
	"os"
	"os/signal"
)

func WaitForCtrlC() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<- sig

	signal.Stop(sig)
}
