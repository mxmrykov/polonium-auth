package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/mxmrykov/polonium-auth/internal/app"
	"github.com/mxmrykov/polonium-auth/internal/config"
)

func main() {
	cfg, ctx := config.Init(), context.Background()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	application, err := app.New(cfg)

	if err != nil {
		log.Fatalln("cannot init app: ", err)
	}

	go func() {
		if err := application.Run(); err != nil {
			fmt.Println("error running app: ", err)
		}
	}()

	<-ctx.Done()
	done := make(chan struct{})
	go func() {
		_ = application.Stop(ctx)
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("application stopped gracefully")
	case <-time.After(10 * time.Second):
		fmt.Println("application shutdown timeout exeeded, forcing quit")
	}
}
