package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// HandleInterrupt sets up signal handling for graceful shutdown on SIGINT and SIGTERM.
// The callback function will be called with the received signal when the process receives an interrupt signal.
// After the callback executes, the process will exit with the appropriate Unix exit code.
func HandleInterrupt(callback func()) {
	HandleSignals(func(sig os.Signal) {
		if callback != nil {
			callback()
		}

		// Calculate exit code following Unix convention: 128 + signal_number
		exitCode := 128
		switch sig {
		case syscall.SIGINT:
			exitCode += int(syscall.SIGINT)
		case syscall.SIGTERM:
			exitCode += int(syscall.SIGTERM)
		}

		os.Exit(exitCode)
	}, syscall.SIGINT, syscall.SIGTERM)
}

// HandleSignals sets up signal handling for specified signals.
// The callback function will be called with the received signal when the process receives any of the specified signals.
func HandleSignals(callback func(os.Signal), signals ...os.Signal) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)

	go func() {
		sig := <-sigChan
		callback(sig)
	}()
}
