package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
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
	r.HandleFunc("/api/load-image/{filename}", loadImageHandler).Methods("GET")
	return r
}

// Set max file size to 5MB.
const maxFileSize = 1024 * 1024 * 5

func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// Set limit for file size.
	if r.ContentLength > maxFileSize {
		http.Error(w, "file size too big", http.StatusBadRequest)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	// Get image from request.
	fr, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}
	defer fr.Close()

	// Create dir images in case one does not exists.
	// Could move to SetRoutes() to create and not check every request.
	// Leaving here to not spread the code.
	if err := os.MkdirAll("images", os.ModeDir); err != nil {
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}

	// Peek first 512 bytes to detect content type.
	sniff, err := bufio.NewReader(fr).Peek(512)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}
	var ext string
	switch mime := http.DetectContentType(sniff); mime {
	case "image/jpeg":
		ext = ".jpeg"
	case "image/png":
		ext = ".png"
	default:
		http.Error(w, "wrong file type", http.StatusBadRequest)
		return
	}

	// Create file name by generating unique name and adding extension
	// determined by http.DetectContentType.
	fname := xid.New().String() + ext

	// Required to read after peeking.
	// TODO: find out why this is needed.
	fr.Seek(0, 0)

	// Create new file using name generated above.
	fs, err := os.Create("images/" + fname)
	if err != nil {
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}
	defer fs.Close()

	// Copy content from request file to system file.
	if _, err := io.Copy(fs, fr); err != nil {
		fmt.Println(err)
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}

	// Respond with JSON containing filename.
	resp := map[string]string{"filename": fname}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}

}

func loadImageHandler(w http.ResponseWriter, r *http.Request) {
	fn := mux.Vars(r)["filename"]

	// Check if name is valid.
	// Name without ext is 20 chars, containing all lowercase sequence of a to v letters and 0 to 9 numbers ([0-9a-v]{20}).
	// https://github.com/rs/xid
	// Valid extensions are jpeg and png.
	if ok, err := regexp.Match(
		"^[0-9a-v]{20}\\.(?:jpeg|png)$",
		[]byte(fn),
	); !ok || err != nil {
		http.Error(w, "unexpected error occured", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, "./images/"+fn)
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
