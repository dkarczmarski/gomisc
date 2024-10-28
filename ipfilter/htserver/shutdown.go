package htserver

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

func RunShutdownListenerTask(ctx context.Context, wg *sync.WaitGroup, server *http.Server) {
	wg.Add(1)
	go RunShutdownListener(ctx, wg, server)
}

func RunShutdownListener(ctx context.Context, wg *sync.WaitGroup, server *http.Server) {
	<-ctx.Done()

	log.Println("the server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)
	log.Println("server is shut down")

	wg.Done()
}
