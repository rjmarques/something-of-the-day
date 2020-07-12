package service

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rjmarques/something-of-the-day/store"
)

type server struct {
	st *store.Store
}

func apiServer(st *store.Store) http.Handler {
	sv := &server{
		st: st,
	}

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v0").Subrouter()
	api.HandleFunc("/something", sv.getSomethingHandler)

	r.PathPrefix("/static").Handler(nocache(http.FileServer(http.Dir("frontend/"))))    // for assets in /static/
	r.PathPrefix("/*").Handler(nocache(http.FileServer(http.Dir("frontend/"))))         // for assets at the same level as index.html
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // all unknown request are sent to index.html
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, "frontend/index.html")
	})

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	return loggedRouter
}

func (sv *server) getSomethingHandler(w http.ResponseWriter, r *http.Request) {
	something := sv.st.GetRand()
	json.NewEncoder(w).Encode(something)
}

func nocache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		h.ServeHTTP(w, r)
	})
}
