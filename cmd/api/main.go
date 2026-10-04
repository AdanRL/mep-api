package main

import (
	"log"
	"mep-api/config"
	"mep-api/db"
	"mep-api/internal/handlers"
	"mep-api/internal/repository"
	"mep-api/internal/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("[FATAL] Please provide the path to the .env file as a command-line argument.")
	}
	log.Printf("[INFO] Using file: %s", os.Args[1])
	envFile := os.Args[1]
	cfg := config.Load(envFile)
	log.Printf("[INIT] Starting MiseEnPlace API in '%s' mode...", cfg.AppEnv)

	router := gin.Default()
	database := db.InitConnection()
	router.SetTrustedProxies(nil)

	ingredientRepo := repository.NewMysqlIngredientRepo(database)
	ingredientHandler := handlers.NewIngredientHandler(ingredientRepo)

	myHandlers := routes.AppHandlers{
		Ingredient: ingredientHandler,
	}

	routes.Setup(router, myHandlers)

	if err := router.Run(cfg.BackendPort); err != nil {
		log.Fatalf("[FATAL] Failed to start server on %s: %v", cfg.BackendPort, err)
	}

}
