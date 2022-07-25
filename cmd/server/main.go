package main

import (
	"github.com/chbm/toca/internal/server/http"
	"github.com/chbm/toca/internal/storage/clerk"
)

func main() {
	clerckCh := clerk.Start()
	httpServer := httpserver.Start(clerkCh)
}
	
