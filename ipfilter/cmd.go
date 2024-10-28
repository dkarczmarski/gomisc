package main

import (
	"context"
	"errors"
	"github.com/dkarczmarski/gomisc/ipfilter/firewall"
	"github.com/dkarczmarski/gomisc/ipfilter/htserver"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	// todo: concurrent
	service := firewall.NewService(
		firewall.WithTimeFunc(time.Now),
		firewall.WithSudoWrapper(),
	)

	mux := htserver.NewServeMux(service)

	var wg sync.WaitGroup

	firewall.RunDeleteOutOfDateTask(ctx, &wg, service)

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	htserver.RunShutdownListenerTask(ctx, &wg, server)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	wg.Wait()
}
