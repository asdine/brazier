package http

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
)

// Handler is the main http handler
type Handler struct {
	Store brazier.Store
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var bucketName, key string

	p := strings.Trim(r.URL.EscapedPath(), "/")
	parts := strings.Split(p, "/")

	if len(parts) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bucketName = parts[0]
	key = parts[1]

	switch r.Method {
	case "PUT":
		h.createItem(w, r, bucketName, key)
	case "GET":
		h.getItem(w, r, bucketName, key)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createItem(w http.ResponseWriter, r *http.Request, bucketName string, key string) {
	var value interface{}
	var buffer bytes.Buffer

	bucket, err := h.Store.Bucket(bucketName)
	if err != nil {
		if err != store.ErrNotFound {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = h.Store.Create(bucketName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bucket, err = h.Store.Bucket(bucketName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	defer bucket.Close()

	_, err = buffer.ReadFrom(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := buffer.Bytes()
	err = json.Unmarshal(data, &value)
	if err != nil {
		data, err = json.Marshal(buffer.String())
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	_, err = bucket.Save(key, data)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getItem(w http.ResponseWriter, r *http.Request, bucketName string, key string) {
	bucket, err := h.Store.Bucket(bucketName)
	if err != nil {
		if err != store.ErrNotFound {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer bucket.Close()

	item, err := bucket.Get(key)
	if err != nil {
		if err != store.ErrNotFound {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(item.Data)
}
