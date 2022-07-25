package main

import (
	"github.com/chbm/toca/internal/server"
	"github.com/chbm/toca/internal/storage"
)

func main() {
	clerkCh := storage.Start()
	httpserver.Start(clerkCh)
}
	
