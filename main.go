package main

import (
	"log"

	"github.com/ngoyal16/owlvault/config"
	"github.com/ngoyal16/owlvault/routes"
)

func main() {
	// Read configurations
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	r := routes.GinEngine(cfg)

	// Start HTTP server
	addr := cfg.Server.Addr
	log.Printf("Server listening on addr %s", addr)
	err = r.Run(addr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
