package main

import (
	"log"
	"mep-api/config"
	"net/http"
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
  router.SetTrustedProxies(nil)

  router.GET("/", func(c *gin.Context) {
      c.JSON(http.StatusOK, gin.H{
          "message": "Welcome to the MiseEnPlace API.",
      })
  })

  if err := router.Run(cfg.BackendPort); err != nil {
      log.Fatalf("[FATAL] Failed to start server on %s: %v", cfg.BackendPort, err)
    }

}
