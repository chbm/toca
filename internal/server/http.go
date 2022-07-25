package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/chbm/toca/internal/storage/commands"
)

func statusRouter() http.hander {
	r := chi.NewRouter()

	r.Get("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r
}


func Start(clerk chan) {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Mount("/_status", statusRouter())		

	router.Get("/default/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		rc := make(chan commands.Value)
		c <- commands.Command{
			op: commands.Get,
			key: key, 
			value: nil,
			r: rc
		}
			res := <-r
			w.Write([]byte(res))
	})

	router.Put("/default/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		rc := make(chan commands.Value)
		c <- commands.Command{
			op: commands.Put,
			key: key,
			value. r.Body,
			r: rc
		}
		res := <- r
		w.WriteHeader(201)
	})

	http.ListenAndServe(":3000", router)
}


