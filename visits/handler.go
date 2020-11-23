package visits

import (
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
)

// SetRoutes sets routes for http.ListenAndServe
func SetRoutes(prod *kafka.Producer) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/visited", visitedHandler)
	return r
}

func visitedHandler(w http.ResponseWriter, r *http.Request) {
}
