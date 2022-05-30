package shutdown

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Gracefully(app *fiber.App, timeout time.Duration) chan struct{} {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	serverShutdown := make(chan struct{}, 1)
	go func() {
		<-sig
		fmt.Println("Gracefully shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		go func() {
			err := app.Shutdown()
			if err != nil {
				log.Printf("Error shutdown server: %v\n", err)
			}
			cancel()
		}()
		select {
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				fmt.Println("Gracefull shutown timed out! Force shutdown")
			case context.Canceled:
				fmt.Println("Gracefull shutdown sucessfull completed")
			}
		}
		serverShutdown <- struct{}{}
	}()
	return serverShutdown
}
