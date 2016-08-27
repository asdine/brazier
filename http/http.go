package http

import (
	"encoding/json"
	"net/http"

	"github.com/asdine/brazier"
)

// Handler is the main http handler
type Handler struct {
	Store brazier.Store
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	h.createBucket(w, r)
}

func (h *Handler) createBucket(w http.ResponseWriter, r *http.Request) {
	var c createBucket

	d := json.NewDecoder(r.Body)
	err := d.Decode(&c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if c.Name == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = h.Store.Create(c.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
