package services

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// SetRoutes sets routes for http.ListenAndServe.
func SetRoutes(h Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/visits", postVisitHandler(h)).Methods("POST")
	r.HandleFunc("/api/visits", getVisitsHandler(h)).Methods("GET")
	r.HandleFunc("/api/visits/{ip}", getVisitsByIPHandler(h)).Methods("GET")
	return r
}

func postVisitHandler(h Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		if err := h.ProduceVisit(ip); err != nil {
			http.Error(w, "unexpected error occured", http.StatusInternalServerError)
			return
		}
	}
}

func getVisitsHandler(h Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := make(map[string]string)
		for f, val := range r.URL.Query() {
			filter[f] = val[0]
		}
		visits, err := h.GetVisits(filter)
		if err != nil {
			http.Error(w, "unexpected error occured", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(visits); err != nil {
			http.Error(w, "unexpected error occured", http.StatusInternalServerError)
			return
		}
	}
}

func getVisitsByIPHandler(h Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filter := make(map[string]string)
		for f, val := range r.URL.Query() {
			filter[f] = val[0]
		}
		visits, err := h.GetVisitsByIP(vars["ip"], filter)
		if err != nil {
			http.Error(w, "unexpected error occured", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(visits); err != nil {
			http.Error(w, "unexpected error occured", http.StatusInternalServerError)
			return
		}
	}
}
