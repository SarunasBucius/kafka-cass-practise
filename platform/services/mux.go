package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/xid"
)

// SetRoutes sets routes for http.ListenAndServe.
func SetRoutes(h Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/visits", postVisitHandler(h)).Methods("POST")
	r.HandleFunc("/api/visits", getVisitsHandler(h)).Methods("GET")
	r.HandleFunc("/api/visits/{ip}", getVisitsByIPHandler(h)).Methods("GET")
	r.HandleFunc("/api/upload-image", uploadImageHandler).Methods("POST")
	return r
}

// Set max file size to 5MB.
const maxFileSize = 1024 * 1024 * 5

func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > maxFileSize {
		http.Error(w, "file size too big", http.StatusBadRequest)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	f, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}

	var ext string
	switch mime := http.DetectContentType(b); mime {
	case "image/jpeg":
		ext = "jpeg"
	case "image/png":
		ext = "png"
	default:
		http.Error(w, "wrong file type", http.StatusBadRequest)
		return
	}

	// Could move to SetRoutes() to create and not check every request.
	// Leaving here to not spread the code.
	if err := os.MkdirAll("images", os.ModeDir); err != nil {
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}

	fname := fmt.Sprintf("%v.%v", xid.New().String(), ext)

	if err := ioutil.WriteFile("images/"+fname, b, 0666); err != nil {
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"filename": fname}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}

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
