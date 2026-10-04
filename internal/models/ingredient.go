package models

import (
	"context"
	"errors"
)

var ErrIngredientNotFound = errors.New("ingredient not found")

type Ingredient struct {
	ID            int     `json:"id"`
	Name          string  `json:"name" binding:"required"`
	UnitOfMeasure string  `json:"unit_of_measure" binding:"required"`
	CurrentStock  float64 `json:"current_stock" binding:"gte=0"`
	MinimumStock  float64 `json:"minimum_stock" binding:"gte=0"`
}

type IngredientRepository interface {
	GetIngredients(ctx context.Context) ([]Ingredient, error)
	GetIngredientByID(ctx context.Context, id int) (*Ingredient, error)
	CreateIngredient(ctx context.Context, newIngredient *Ingredient) (*Ingredient, error)
	UpdateIngredient(ctx context.Context, id int, data *Ingredient) (*Ingredient, error)
	DeleteIngredient(ctx context.Context, id int) error
}
