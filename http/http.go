package http

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/json"
	"github.com/asdine/brazier/store"
)

// Serve runs the HTTP server
func Serve(s brazier.Store, port int) error {
	http.Handle("/", &Handler{Store: s})
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

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
		h.saveItem(w, r, bucketName, key)
	case "GET":
		h.getItem(w, r, bucketName, key)
	case "DELETE":
		h.deleteItem(w, r, bucketName, key)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) saveItem(w http.ResponseWriter, r *http.Request, bucketName string, key string) {
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

	if r.ContentLength == 0 {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := buffer.Bytes()
	if !json.IsValid(data) {
		var b bytes.Buffer
		b.Grow(len(data) + 2)
		b.WriteByte('"')
		b.Write(data)
		b.WriteByte('"')
		data = b.Bytes()
	} else {
		data = json.Clean(data)
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

func (h *Handler) deleteItem(w http.ResponseWriter, r *http.Request, bucketName string, key string) {
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

	err = bucket.Delete(key)
	if err != nil {
		if err != store.ErrNotFound {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
