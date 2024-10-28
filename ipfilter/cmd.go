package main

import (
	"context"
	"errors"
	"github.com/dkarczmarski/gomisc/ipfilter/firewall"
	"github.com/dkarczmarski/gomisc/ipfilter/htserver"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	mux := http.NewServeMux()

	// todo: concurrent
	service := firewall.NewService(
		firewall.WithTimeFunc(time.Now),
		firewall.WithSudoWrapper(),
	)

	mux.HandleFunc("POST /api/me/add", func(w http.ResponseWriter, r *http.Request) {
		htserver.HandleAddMe(w, r, service)
	})

	mux.HandleFunc("POST /api/me/delete", func(w http.ResponseWriter, r *http.Request) {
		htserver.HandleDeleteMe(w, r, service)
	})

	mux.HandleFunc("POST /api/ip/add", func(w http.ResponseWriter, r *http.Request) {
		htserver.HandleAddIP(w, r, service)
	})

	mux.HandleFunc("POST /api/ip/delete", func(w http.ResponseWriter, r *http.Request) {
		htserver.HandleDeleteIP(w, r, service)
	})

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("templ/index.html"))

		entries := service.List()

		if err := templ.Execute(w, entries); err != nil {
			log.Fatal(err)
		}
	})

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
