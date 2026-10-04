package routes

import (
	"net/http"

	"mep-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

type AppHandlers struct {
	Ingredient *handlers.IngredientHandler
}

func Setup(router *gin.Engine, app AppHandlers) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the MiseEnPlace API.",
		})
	})
	v1 := router.Group("/api/v1")
	RegisterIngredientRoutes(v1, app.Ingredient)
}
