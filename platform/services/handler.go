package services

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

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
