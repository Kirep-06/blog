package main

import (
	"fmt"
	"log"

	"blog/config"
	"blog/internal/database"
	"blog/internal/router"
	"blog/internal/seed"
	"blog/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load config
	if err := config.Load(); err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 2. Connect database and migrate
	if err := database.Connect(); err != nil {
		log.Fatalf("database: %v", err)
	}

	// 3. Seed initial user
	if err := seed.Run(); err != nil {
		log.Fatalf("seed: %v", err)
	}

	// 4. Init storage provider
	provider, err := initStorage()
	if err != nil {
		log.Fatalf("storage: %v", err)
	}

	// 5. Setup router
	gin.SetMode(config.C.Server.Mode)
	engine := gin.Default()

	// Serve local uploads if using local storage
	if config.C.Storage.Driver == "local" {
		engine.Static(config.C.Storage.Local.URLPrefix, config.C.Storage.Local.UploadDir)
	}

	router.Setup(engine, provider)

	// Serve frontend static files
	engine.Static("/frontend", "./frontend")
	engine.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/frontend/index.html")
	})

	// 6. Run
	addr := fmt.Sprintf(":%d", config.C.Server.Port)
	log.Printf("server: listening on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func initStorage() (storage.StorageProvider, error) {
	cfg := config.C.Storage
	switch cfg.Driver {
	case "s3":
		return storage.NewS3Storage(
			cfg.S3.Region,
			cfg.S3.Endpoint,
			cfg.S3.AccessKeyID,
			cfg.S3.SecretAccessKey,
			cfg.S3.Bucket,
			cfg.S3.PublicURLBase,
			cfg.S3.ForcePathStyle,
		)
	default: // "local"
		baseURL := config.C.Server.BaseURL + cfg.Local.URLPrefix
		return storage.NewLocalStorage(cfg.Local.UploadDir, baseURL)
	}
}
