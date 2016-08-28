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

	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	h.createItem(w, r, bucketName, key)
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
