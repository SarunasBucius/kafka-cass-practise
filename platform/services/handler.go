package services

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// ListenHTTP listens and serves http requests
func ListenHTTP(ctx context.Context, k *kcp.Kcp, cancel context.CancelFunc, wg *sync.WaitGroup) {
	srv := &http.Server{
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      SetRoutes(k),
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

// SetRoutes sets routes for http.ListenAndServe
func SetRoutes(k *kcp.Kcp) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/visited", visitedHandler(k))
	return r
}

func visitedHandler(k *kcp.Kcp) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		if err := k.ProduceVisit(ip); err != nil {
			http.Error(w, "unexpected error occured", http.StatusInternalServerError)
			return
		}
	}
}
