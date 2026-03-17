package main

import (
	"fmt"
	"gamebook-backend/config"
	"log"
	"os"
	"path/filepath"

	"gamebook-backend/middlewares"
	"gamebook-backend/modules/auth"
	"gamebook-backend/modules/game"
	"gamebook-backend/modules/user"
	"gamebook-backend/providers"
	"gamebook-backend/script"

	"github.com/samber/do"

	"github.com/gin-gonic/gin"
)

func args(injector *do.Injector) bool {
	if len(os.Args) > 1 {
		flag := script.Commands(injector)
		return flag
	}

	return true
}

func run(cfg *config.Config, server *gin.Engine) {
	server.Static("/assets", "./assets")

	serve := fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

func getRootPath() string {
	return filepath.Dir("..")
}

func main() {
	var (
		injector = do.New()
	)

	config.LoadEnv(getRootPath())
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("❌ Configuration error: %v", err)
	}

	// Validate JWT configuration
	if err := config.ValidateJWTConfig(&cfg.JWT); err != nil {
		log.Fatalf("❌ JWT configuration error: %v", err)
	}

	// Setup error logger
	if err := config.SetupErrorLogger(); err != nil {
		log.Printf("failed to setup error logger: %v", err)
	}

	providers.RegisterDependencies(cfg, injector)

	if !args(injector) {
		return
	}

	server := gin.Default()
	server.Use(middlewares.ErrorLoggerMiddleware())
	server.Use(middlewares.CORSMiddleware())

	// Register module routes
	auth.RegisterRoutes(server, injector)
	user.RegisterRoutes(server, injector)
	game.RegisterRoutes(server, injector)

	run(cfg, server)
}
