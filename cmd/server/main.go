package main

import (
	"log"
	"github.com/chbm/toca/internal/server"
	"github.com/chbm/toca/internal/storage"
)

func main() {
	logger := log.Default()

	clerkCh := storage.Start()
	router := httpserver.Start(clerkCh)
	
	logger.Printf("http listening on :3000")
	http.ListenAndServe(":3000", router)
}
	
