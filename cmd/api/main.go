package main

import (
	"log"
	"net/http"

	"mep-api/internal/config"

	"github.com/gin-gonic/gin"
)


func main() {
  cfg := config.Load()
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
