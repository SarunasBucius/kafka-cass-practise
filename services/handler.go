package services

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
)

// SetRoutes sets routes for http.ListenAndServe
func SetRoutes(prod *kcp.Produce) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/visited", visitedHandler(prod))
	return r
}

func visitedHandler(prod *kcp.Produce) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := prod.ProduceVisit(); err != nil {
			http.Error(w, "unexpected error occured", http.StatusInternalServerError)
			return
		}
	}
}
