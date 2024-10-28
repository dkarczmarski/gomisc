package htserver

import (
	"github.com/dkarczmarski/gomisc/ipfilter/firewall"
	"html/template"
	"log"
	"net/http"
)

type ServeMux struct {
	mux *http.ServeMux
}

func NewServeMux(service *firewall.Service) *ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/me/add", func(w http.ResponseWriter, r *http.Request) {
		user := HandleBasicAuth(w, r, createAuthFunc(users))
		if user == nil {
			log.Println("User not authorized. Redirect to /login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		HandleAddMe(w, r, service)
	})

	mux.HandleFunc("POST /api/me/delete", func(w http.ResponseWriter, r *http.Request) {
		user := HandleBasicAuth(w, r, createAuthFunc(users))
		if user == nil {
			log.Println("User not authorized. Redirect to /login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		HandleDeleteMe(w, r, service)
	})

	mux.HandleFunc("POST /api/ip/add", func(w http.ResponseWriter, r *http.Request) {
		user := HandleBasicAuth(w, r, createAuthFunc(users))
		if user == nil {
			log.Println("User not authorized. Redirect to /login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		HandleAddIP(w, r, service)
	})

	mux.HandleFunc("POST /api/ip/delete", func(w http.ResponseWriter, r *http.Request) {
		user := HandleBasicAuth(w, r, createAuthFunc(users))
		if user == nil {
			log.Println("User not authorized. Redirect to /login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		HandleDeleteIP(w, r, service)
	})

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		user := HandleBasicAuth(w, r, createAuthFunc(users))
		log.Printf("user: %+v", user)
		if user == nil {
			return
		}

		templ := template.Must(template.ParseFiles("templates/index.html"))

		entries := service.List()

		if err := templ.Execute(w, map[string]interface{}{
			"User":    user,
			"Entries": entries,
		}); err != nil {
			log.Fatal(err)
		}
	})

	return &ServeMux{
		mux: mux,
	}
}

func (srv *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.mux.ServeHTTP(w, r)
}
