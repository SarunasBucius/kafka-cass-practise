package services

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// Handler contains methods to handle request.
type Handler interface {
	ProduceVisit(ip string) error
	GetVisits(filter map[string]string) (kcp.VisitsByIP, error)
}

// ListenHTTP listens and serves http requests.
func ListenHTTP(ctx context.Context, h Handler, cancel context.CancelFunc, wg *sync.WaitGroup) {
	srv := &http.Server{
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      SetRoutes(h),
	}
	go func() {
		defer wg.Done()
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()
	if err := srv.ListenAndServe(); err != nil {
		cancel()
	}
}

// SetRoutes sets routes for http.ListenAndServe.
func SetRoutes(h Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/visits", postVisitHandler(h)).Methods("POST")
	r.HandleFunc("/api/visits", getVisitsHandler(h)).Methods("GET")
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
