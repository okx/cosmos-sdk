package lcd

import (
	"os"
	"os/signal"
	"syscall"
	//"github.com/cosmos/cosmos-sdk/server"

	// unnamed import of statik for swagger UI support
	_ "github.com/cosmos/cosmos-sdk/client/lcd/statik"
)

func trapSignal(cleanupFunc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		if cleanupFunc != nil {
			cleanupFunc()
		}
		exitCode := 128
		switch sig {
		case syscall.SIGINT:
			exitCode += int(syscall.SIGINT)
		case syscall.SIGTERM:
			exitCode += int(syscall.SIGTERM)
		}
		os.Exit(exitCode)
	}()
}
