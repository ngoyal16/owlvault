package main

import (
	"github.com/ngoyal16/owlvault/controllers/ks2"
	"github.com/ngoyal16/owlvault/keyprovider"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ngoyal16/owlvault/config"
	"github.com/ngoyal16/owlvault/encrypt"
	"github.com/ngoyal16/owlvault/middleware"
	"github.com/ngoyal16/owlvault/storage"
	"github.com/ngoyal16/owlvault/vault"
)

func main() {
	// Read configurations
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	// Initialize storage based on configuration
	dbStorage, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize encryptor based on configuration
	keyProvider, err := keyprovider.NewKeyProvider(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize key provider: %v", err)
	}

	// Initialize encryptor based on configuration
	encryptor, err := encrypt.NewEncryptor(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize encryptor: %v", err)
	}

	// Initialize OwlVault with the chosen storage implementation
	owlVault := vault.NewOwlVault(dbStorage, keyProvider, encryptor)

	// Create a new Gorilla Mux router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
		return
	})

	router.Use(middleware.CORSMiddleware())

	v1 := router.Group("v1")
	{
		v1.GET("ks2", ks2.KS2(owlVault))
		v1.POST("ks2", ks2.KS2(owlVault))
		v1.PATCH("ks2", ks2.KS2(owlVault))
	}

	// Start HTTP server
	addr := cfg.Server.Addr
	log.Printf("Server listening on addr %s", addr)
	err = router.Run(addr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
