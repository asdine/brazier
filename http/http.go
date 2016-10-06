package http

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	graceful "gopkg.in/tylerb/graceful.v1"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/json"
	"github.com/asdine/brazier/store"
)

// NewServer returns a configured HTTP server
func NewServer(r brazier.Registry, s brazier.Store) brazier.Server {
	http.Handle("/", &Handler{Registry: r, Store: s})
	srv := graceful.Server{
		Server: &http.Server{},
	}

	return &srv
}

// Handler is the main http handler
type Handler struct {
	Registry brazier.Registry
	Store    brazier.Store
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var bucketName, key string

	p := strings.Trim(r.URL.EscapedPath(), "/")
	parts := strings.Split(p, "/")

	if len(parts) > 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bucketName = parts[0]
	if len(parts) > 1 {
		key = parts[1]
	}

	switch r.Method {
	case "PUT":
		h.saveItem(w, r, bucketName, key)
	case "GET":
		if key != "" {
			h.getItem(w, r, bucketName, key)
		} else {
			h.listBucket(w, r, bucketName)
		}
	case "DELETE":
		h.deleteItem(w, r, bucketName, key)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) saveItem(w http.ResponseWriter, r *http.Request, bucketName string, key string) {
	bucket, err := store.GetBucketOrCreate(h.Registry, h.Store, bucketName)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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

	data := json.ToValidJSON(buffer.Bytes())
	_, err = bucket.Save(key, data)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getItem(w http.ResponseWriter, r *http.Request, bucketName string, key string) {
	info, err := h.Registry.BucketInfo(bucketName)
	if err != nil {
		if err != store.ErrNotFound {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bucket, err := h.Store.Bucket(info.Name)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
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
	info, err := h.Registry.BucketInfo(bucketName)
	if err != nil {
		if err != store.ErrNotFound {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bucket, err := h.Store.Bucket(info.Name)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
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

func (h *Handler) listBucket(w http.ResponseWriter, r *http.Request, bucketName string) {
	info, err := h.Registry.BucketInfo(bucketName)
	if err != nil {
		if err != store.ErrNotFound {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bucket, err := h.Store.Bucket(info.Name)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer bucket.Close()

	items, err := bucket.Page(1, -1)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	raw, err := json.MarshalList(items)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(raw)
}
