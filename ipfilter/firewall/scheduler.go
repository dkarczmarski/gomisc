package firewall

import (
	"context"
	"log"
	"sync"
	"time"
)

func RunDeleteOutOfDateTask(ctx context.Context, wg *sync.WaitGroup, service *Service) {
	wg.Add(1)
	go runDeleteOutOfDateTask(ctx, wg, service)
}

func runDeleteOutOfDateTask(ctx context.Context, wg *sync.WaitGroup, service *Service) {
loop:
	for {
		deleted, err := func() ([]IPEntry, error) {
			srvCtx, srvCtxCancel := context.WithTimeout(ctx, 10*time.Second)
			defer srvCtxCancel()
			return service.DeleteOutOfDateCtx(srvCtx, 15*time.Second)
		}()
		if err != nil {
			log.Print(err)
		}
		if len(deleted) > 0 {
			log.Printf("deleted out-of-date entries: %+v", deleted)
		}

		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			log.Printf("out-of-date scheduler: %v", ctx.Err())
			break loop
		}
	}

	log.Printf("firewall entries: %+v", service.List())
	wg.Done()
}
