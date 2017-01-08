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
func NewServer(r *store.Store) brazier.Server {
	http.Handle("/", &Handler{Store: r})
	srv := graceful.Server{
		Server: &http.Server{},
	}

	return &srv
}

// Handler is the main http handler
type Handler struct {
	Store *store.Store
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rawPath := r.URL.EscapedPath()

	switch r.Method {
	case "PUT":
		h.putItem(w, r, rawPath)
	case "GET":
		h.getNode(w, r, rawPath)
	case "DELETE":
		h.deleteItem(w, r, rawPath)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) putItem(w http.ResponseWriter, r *http.Request, rawPath string) {
	if r.ContentLength == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r.Body)
	r.Body.Close()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = h.Store.Put(rawPath, json.ToValidJSON(buffer.Bytes()))
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getNode(w http.ResponseWriter, r *http.Request, rawPath string) {
	if !strings.HasSuffix(rawPath, "/") {
		item, err := h.Store.Get(rawPath)
		if err != nil {
			if err == store.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(item.Data)
		return
	}

	var items []brazier.Item
	var err error
	if r.URL.Query().Get("recursive") != "" {
		items, err = h.Store.Tree(rawPath)
	} else {
		items, err = h.Store.List(rawPath, 1, -1)
	}
	if err != nil {
		if err == store.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.MarshalList(items)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *Handler) deleteItem(w http.ResponseWriter, r *http.Request, rawPath string) {
	err := h.Store.Delete(rawPath)
	if err != nil {
		if err != store.ErrNotFound {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
