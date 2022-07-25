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


func Start(clerk chan storage.Command) {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Mount("/_status", statusRouter())		

	router.Get("/default/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		rc := make(chan storage.Value)
		clerk <- storage.Command{
			Op: storage.Get,
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

	router.Put("/default/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		bodyB, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rc := make(chan storage.Value)
		clerk <-storage.Command{
			Op: storage.Put,
			Key: key,
			Value: string(bodyB),
			R: rc,
		}
		<-rc
		w.WriteHeader(201)
	})

	http.ListenAndServe(":3000", router)
}


