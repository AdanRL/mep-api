package handlers

import (
	"errors"
	"mep-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	repo models.IngredientRepository
}
func NewIngredientHandler(repo models.IngredientRepository) *IngredientHandler {
	return &IngredientHandler{repo: repo}
}

func (h *IngredientHandler) Create(c *gin.Context) {
	var input models.Ingredient

	if err := c.ShouldBindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Not valid data: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	createdIngredient, err := h.repo.CreateIngredient(ctx, &input)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating ingredient: " + err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, createdIngredient)
}

func (h *IngredientHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	ingredients, err := h.repo.GetIngredients(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error getting all ingredients"})
		return
	}

	c.JSON(http.StatusOK, ingredients)
}

func (h *IngredientHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "This is not a valid ID.",
		})
		return
	}

	ingredient, err := h.repo.GetIngredientByID(ctx, id)
	if err != nil {

		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "Ingredient not found.",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, ingredient)
}

func (h *IngredientHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "This is not a valid int ID.",
		})
		return
	}
	err = h.repo.DeleteIngredient(ctx, id)

	if err != nil {
		if errors.Is(err, models.ErrIngredientNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"error": "Ingredient not found.",
			})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Error deleting ingredient: " + err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)

}

func (h *IngredientHandler) Update(c *gin.Context) {
	var input models.Ingredient

	if err := c.ShouldBindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Not valid data: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "This is not a valid ID.",
		})
		return
	}
	input.ID = id
	ingredient, err := h.repo.UpdateIngredient(ctx, id, &input)
	if err != nil {

		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "Ingredient not found.",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, ingredient)
}
