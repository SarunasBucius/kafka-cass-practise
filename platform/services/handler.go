package services

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// ListenHTTP listens and serves http requests
func ListenHTTP(ctx context.Context, srv *http.Server, cancel context.CancelFunc, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()
	if err := srv.ListenAndServe(); err != nil {
		cancel()
	}
}

// SetRoutes sets routes for http.ListenAndServe
func SetRoutes(k *kcp.Kcp) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/visited", visitedHandler(k))
	return r
}

func visitedHandler(k *kcp.Kcp) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := k.ProduceVisit(); err != nil {
			http.Error(w, "unexpected error occured", http.StatusInternalServerError)
			return
		}
	}
}
