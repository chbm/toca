package httpserver

import (
	"net/http"
	"io/ioutil"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/chbm/toca/internal/storage"
)

func statusRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	return r
}


func Start(clerk chan storage.Command) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Mount("/_status", statusRouter())		

	router.Post("/{ns}", func(w http.ResponseWriter, r *http.Request) {
		ns := chi.URLParam(r, "ns")
		
		rc := make(chan storage.Value)
		clerk <- storage.Command{
			Op: storage.CreateNs,
			Ns: ns,
			Key: "", 
			Value: "",
			R: rc,
		}
		res := <-rc
		if res.Exists {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(201)
		}
	})

	// XXX so much boilerplate ...
	router.Get("/{ns}/{key}", func(w http.ResponseWriter, r *http.Request) {
		ns := chi.URLParam(r, "ns")
		key := chi.URLParam(r, "key")

		rc := make(chan storage.Value)
		clerk <- storage.Command{
			Op: storage.Get,
			Ns: ns,
			Key: key, 
			Value: "",
			R: rc,
		}
		res := <-rc
		if res.Exists {
			w.Write([]byte(res.V)) 
		} else {
			w.WriteHeader(404)
		}
	})

	router.Put("/{ns}/{key}", func(w http.ResponseWriter, r *http.Request) {
		ns := chi.URLParam(r, "ns")
		key := chi.URLParam(r, "key")
		bodyB, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rc := make(chan storage.Value)
		clerk <-storage.Command{
			Op: storage.Put,
			Ns: ns,
			Key: key,
			Value: string(bodyB),
			R: rc,
		}
		res := <-rc
		if res.Exists {
			w.WriteHeader(204)
		} else {
			w.WriteHeader(404)
		}
	})

	router.Delete("/{ns}/{key}", func(w http.ResponseWriter, r *http.Request) {
		ns := chi.URLParam(r, "ns")
		key := chi.URLParam(r, "key")
		rc := make(chan storage.Value)
		clerk <-storage.Command{
			Op: storage.Delete,
			Ns: ns,
			Key: key,
			Value: "",
			R: rc,
		}
		ret := <-rc
		if ret.Exists {
			w.WriteHeader(204)
		} else {
			w.WriteHeader(404)
		}
		
	})

	return router 
}


