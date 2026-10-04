package routes

import (
	"mep-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterIngredientRoutes(v1 *gin.RouterGroup, ingredientHandler *handlers.IngredientHandler) {
	{
		ingredients := v1.Group("/ingredients")

		{
			ingredients.GET("", ingredientHandler.GetAll)
			ingredients.GET("/:id", ingredientHandler.GetByID)
			ingredients.POST("", ingredientHandler.Create)
			ingredients.DELETE("/:id", ingredientHandler.Delete)
			ingredients.PUT("/:id", ingredientHandler.Update)
		}
	}
}
