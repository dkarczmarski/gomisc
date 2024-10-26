package main

import (
	"github.com/dkarczmarski/gomisc/ipfilter/firewall"
	"html/template"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	service := firewall.Service{
		WrapperCmd: "sudo",
	}

	mux.HandleFunc("POST /api/ip/add", func(w http.ResponseWriter, r *http.Request) {
		//todo: get client IP address

		ip := r.FormValue("ip")
		log.Printf("ip: %v", ip)

		if err := service.AddIP(ip); err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("POST /api/ip/delete", func(w http.ResponseWriter, r *http.Request) {
		//todo: get client IP address

		ip := r.FormValue("ip")
		log.Printf("ip: %v", ip)

		if err := service.DeleteIP(ip); err != nil {
			log.Println(err)
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

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
