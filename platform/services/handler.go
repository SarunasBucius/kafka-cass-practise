package services

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Handler contains methods to handle request.
type Handler interface {
	ProduceVisit(ip string) error
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
	r.HandleFunc("/api/visited", visitedHandler(h))
	return r
}

func visitedHandler(h Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		if err := h.ProduceVisit(ip); err != nil {
			http.Error(w, "unexpected error occured", http.StatusInternalServerError)
			return
		}
	}
}
