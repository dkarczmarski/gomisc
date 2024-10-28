package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dkarczmarski/gomisc/ipfilter/firewall"
	"html/template"
	"log"
	"net"
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
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(fmt.Errorf("net.SplitHostPort(): %w", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ip := host
		log.Printf("ip: %v", ip)

		if err := service.AddIPCtx(r.Context(), ip); err != nil {
			log.Println(fmt.Errorf("service.AddIPCtx(): %w", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("POST /api/me/delete", func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(fmt.Errorf("net.SplitHostPort(): %w", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ip := host
		log.Printf("ip: %v", ip)

		if err := service.DeleteIPCtx(r.Context(), ip); err != nil {
			log.Println(fmt.Errorf("service.DeleteIPCtx(): %w", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("POST /api/ip/add", func(w http.ResponseWriter, r *http.Request) {
		ip := r.FormValue("ip")
		if len(ip) == 0 {
			log.Println("no param: ip")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("ip: %v", ip)

		if err := service.AddIPCtx(r.Context(), ip); err != nil {
			log.Println(fmt.Errorf("service.AddIPCtx(): %w", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("POST /api/ip/delete", func(w http.ResponseWriter, r *http.Request) {
		ip := r.FormValue("ip")
		if len(ip) == 0 {
			log.Println("no param: ip")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("ip: %v", ip)

		if err := service.DeleteIPCtx(r.Context(), ip); err != nil {
			log.Println(fmt.Errorf("service.DeleteIPCtx(): %w", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("templ/index.html"))

		entries := service.List()

		if err := templ.Execute(w, entries); err != nil {
			log.Fatal(err)
		}
	})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
	loop:
		for {
			deleted, err := func() ([]firewall.IPEntry, error) {
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
	}()

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	wg.Add(1)
	go func() {
		<-ctx.Done()
		log.Println("the server is shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
		log.Println("server is shut down")
		wg.Done()
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	wg.Wait()
}
