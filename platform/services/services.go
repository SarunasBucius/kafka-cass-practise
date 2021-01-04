package services

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// Handler contains methods to handle request.
type Handler interface {
	ProduceVisit(ip string) error
	GetVisits(filter map[string]string) (kcp.VisitsByIP, error)
	GetVisitsByIP(ip string, filter map[string]string) (kcp.VisitsByIP, error)
}

// ListenHTTP listens and serves http requests.
func ListenHTTP(ctx context.Context, h http.Handler, cancel context.CancelFunc, wg *sync.WaitGroup) {
	srv := &http.Server{
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h,
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
