// Package htserver provides application HTTP server implementation.
package htserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
)

type Firewall interface {
	AddIPCtx(ctx context.Context, ip string) error
	DeleteIPCtx(ctx context.Context, ip string) error
}

func HandleAddMe(w http.ResponseWriter, r *http.Request, service Firewall) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(fmt.Errorf("net.SplitHostPort(): %w", err))
		w.WriteHeader(http.StatusInternalServerError)
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
}

func HandleDeleteMe(w http.ResponseWriter, r *http.Request, service Firewall) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(fmt.Errorf("net.SplitHostPort(): %w", err))
		w.WriteHeader(http.StatusInternalServerError)
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
}

func HandleAddIP(w http.ResponseWriter, r *http.Request, service Firewall) {
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
}

func HandleDeleteIP(w http.ResponseWriter, r *http.Request, service Firewall) {
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
}
