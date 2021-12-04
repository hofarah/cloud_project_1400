package createService

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func processConfig() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func gracefulShutdown(app *fiber.App) {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM) // catch OS signals
		<-sigint
		if err := app.Shutdown(); err != nil {
			log.Printf("API server Shutdown: %v", err)
		}
	}()
}
