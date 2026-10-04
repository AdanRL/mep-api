package repository

import (
	"context"
	"database/sql"
	"fmt"
	"mep-api/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlIngredientRepo struct {
	db *sql.DB
}

func NewMysqlIngredientRepo(db *sql.DB) *MysqlIngredientRepo {
	return &MysqlIngredientRepo{db: db}
}

func (r *MysqlIngredientRepo) GetIngredients(ctx context.Context) ([]models.Ingredient, error) {
	query := `
        SELECT id, name, unit_of_measure, current_stock, minimum_stock
        FROM ingredients
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying ingredients: %w", err)
	}
	defer rows.Close()

	ingredients := make([]models.Ingredient, 0)
	for rows.Next() {
		ingredient := models.Ingredient{}
		if err := rows.Scan(
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.UnitOfMeasure,
			&ingredient.CurrentStock,
			&ingredient.MinimumStock,
		); err != nil {
			return nil, fmt.Errorf("scanning ingredient: %w", err)
		}

		ingredients = append(ingredients, ingredient)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating ingredients: %w", err)
	}

	return ingredients, nil
}

func (r *MysqlIngredientRepo) GetIngredientByID(ctx context.Context, id int) (*models.Ingredient, error) {

	query := `
        SELECT id, name, unit_of_measure, current_stock, minimum_stock
        FROM ingredients
        WHERE id = ?
    `
	ingredient := models.Ingredient{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ingredient.ID,
		&ingredient.Name,
		&ingredient.UnitOfMeasure,
		&ingredient.CurrentStock,
		&ingredient.MinimumStock,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning ingredient %d: %w", id, err)
	}

	return &ingredient, nil
}

func (r *MysqlIngredientRepo) CreateIngredient(ctx context.Context, ing *models.Ingredient) (*models.Ingredient, error) {
	query := `
        INSERT INTO ingredients (name, unit_of_measure, current_stock, minimum_stock)
        VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		ing.Name, ing.UnitOfMeasure, ing.CurrentStock, ing.MinimumStock,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	ing.ID = int(id)
	return ing, nil
}

func (r *MysqlIngredientRepo) DeleteIngredient(ctx context.Context, id int) error {
	query := `
        DELETE FROM ingredients
        WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrIngredientNotFound
	}

	return nil
}

func (r *MysqlIngredientRepo) UpdateIngredient(
	ctx context.Context,
	id int,
	ing *models.Ingredient) (*models.Ingredient, error) {

	query := `
          UPDATE ingredients
          SET name = ?, unit_of_measure = ?, current_stock = ?, minimum_stock = ?
          WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		ing.Name,
		ing.UnitOfMeasure,
		ing.CurrentStock,
		ing.MinimumStock,
		id,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, models.ErrIngredientNotFound
	}

	return ing, nil
}
