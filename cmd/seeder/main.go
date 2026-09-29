package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func cleanDynamicTables(db *sql.DB) {

	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	db.Exec("TRUNCATE TABLE dish_ingredients")
	db.Exec("TRUNCATE TABLE dishes")
	db.Exec("TRUNCATE TABLE ingredients")
	db.Exec("TRUNCATE TABLE user_roles")
	db.Exec("TRUNCATE TABLE users")
  db.Exec("TRUNCATE TABLE roles")

	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	log.Println("🧹 Cleaning dynamic tables...")
}

func main() {
  err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found. Using environment variables.")
	}
  user := os.Getenv("DB_USER")
  password := os.Getenv("DB_PASSWORD")
  host := os.Getenv("DB_HOST")
  port := os.Getenv("DB_PORT")
  dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	gofakeit.Seed(0)

  cleanDynamicTables(db)

	seedRoles(db)
	seedCategories(db)
	seedAllergens(db)
	seedTables(db)
	seedUsers(db, 15)
  ingredientIDs := seedIngredients(db)
	seedDishes(db, ingredientIDs)

	log.Println("✅ Database seeding completed successfully.")
}

func seedRoles(db *sql.DB) {
	roles := map[int]string{1: "Administrador", 2: "Mesero", 3: "Cocinero"}
	for id, name := range roles {
		db.Exec("INSERT IGNORE INTO roles (id, name) VALUES (?, ?)", id, name)
	}
	log.Println("🔹 Inserting static roles...")
}

func seedCategories(db *sql.DB) {
	categories := map[int]string{1: "Entrantes", 2: "Platos Principales", 3: "Postres", 4: "Bebidas"}
	for id, name := range categories {
		db.Exec("INSERT IGNORE INTO categories (id, name) VALUES (?, ?)", id, name)
	}
	log.Println("🔹 Inserting static categories...")
}

func seedAllergens(db *sql.DB) {
	allergens := []struct {
		id   int
		name string
		icon string
	}{
		{1, "Gluten", "🌾"}, {2, "Lácteos", "🥛"},
		{3, "Frutos Secos", "🥜"}, {4, "Crustáceos", "🦐"},
	}
	for _, a := range allergens {
		db.Exec("INSERT IGNORE INTO allergens (id, name, icon) VALUES (?, ?, ?)", a.id, a.name, a.icon)
	}
	log.Println("🔹 Inserting static allergens...")
}

func seedTables(db *sql.DB) {
	capacities := []int{2, 2, 4, 4, 4, 6, 6, 8, 2, 10}
	for i, cap := range capacities {
		tableNum := i + 1
		db.Exec("INSERT IGNORE INTO dining_tables (id, table_number, capacity, status) VALUES (?, ?, ?, 'AVAILABLE')", tableNum, tableNum, cap)
	}
	log.Println("🔹 Inserting static dining tables...")
}

func seedUsers(db *sql.DB, count int) {
  plainPassword := os.Getenv("MYSQL_USERS_APP_PASSWORD")

  hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error generando el hash de la contraseña: %v", err)
	}
	hashedPassword := string(hashedBytes)

	for i := 0; i < count; i++ {
		res, err := db.Exec(
			"INSERT INTO users (username, email, password_hash, is_active) VALUES (?, ?, ?, ?)",
			gofakeit.Username(),
			gofakeit.Email(),
			hashedPassword,
			gofakeit.Bool(),
		)
		if err != nil {
			continue
		}

		userID, _ := res.LastInsertId()
		roleID := gofakeit.Number(1, 3)
		db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID)
	}
	log.Printf("🔹 %d Users inserted.", count)
}

func seedIngredients(db *sql.DB) map[string]int64 {
  ingredientIDs := make(map[string]int64)

  ingredients := []struct {
		Name  string
		Unit  string
		Stock float64
		Min   float64
	}{
		{"Pan de hamburguesa", "UNITS", 200, 50},
		{"Carne de ternera", "KG", 40.0, 10.0},
		{"Queso Cheddar", "KG", 15.0, 3.0},
		{"Lechuga", "KG", 10.0, 2.0},
		{"Tomate", "KG", 20.0, 5.0},
		{"Pechuga de Pollo", "KG", 30.0, 8.0},
		{"Salsa César", "LTS", 5.0, 1.0},
		{"Picatostes", "KG", 8.0, 2.0},
		{"Patatas", "KG", 100.0, 20.0},
		{"Aceite de Oliva", "LTS", 25.0, 5.0},
		{"Lata Coca-Cola", "UNITS", 500, 100},
		{"Agua Mineral 500ml", "UNITS", 300, 50},
		{"Cerveza Artesanal", "UNITS", 200, 50},
		{"Helado de Vainilla", "KG", 15.0, 3.0},
		{"Brownie", "UNITS", 50, 10},
	}

  for _, ing := range ingredients {
		res, err := db.Exec(
			"INSERT INTO ingredients (name, unit_of_measure, current_stock, minimum_stock) VALUES (?, ?, ?, ?)",
			ing.Name, ing.Unit, ing.Stock, ing.Min,
		)
		if err == nil {
			id, _ := res.LastInsertId()
			ingredientIDs[ing.Name] = id
		}
	}
	log.Printf("🔹 %d Ingredients inserted.", len(ingredients))
  return ingredientIDs
}

func seedDishes(db *sql.DB, ingredientIDs map[string]int64) {
  dishes := []struct {
		Name       string
		CategoryID int
		Price      float64
		Recipe     map[string]float64
	}{
		{
			Name: "Hamburguesa Clásica", CategoryID: 2, Price: 12.50,
			Recipe: map[string]float64{
				"Pan de hamburguesa": 1.0, "Carne de ternera": 0.200,
				"Queso Cheddar": 0.050, "Lechuga": 0.050, "Tomate": 0.100,
			},
		},
		{
			Name: "Ensalada César", CategoryID: 1, Price: 9.90,
			Recipe: map[string]float64{
				"Lechuga": 0.150, "Pechuga de Pollo": 0.150,
				"Salsa César": 0.050, "Crutones": 0.030,
			},
		},
		{
			Name: "Patatas Bravas", CategoryID: 1, Price: 6.50,
			Recipe: map[string]float64{
				"Patatas": 0.300, "Aceite de Oliva": 0.050,
			},
		},
		{
			Name: "Brownie con Helado", CategoryID: 3, Price: 5.50,
			Recipe: map[string]float64{
				"Brownie": 1.0, "Helado de Vainilla": 0.100,
			},
		},
		{
			Name: "Coca-Cola", CategoryID: 4, Price: 2.50,
			Recipe: map[string]float64{"Lata Coca-Cola": 1.0},
		},
		{
			Name: "Agua Mineral", CategoryID: 4, Price: 1.80,
			Recipe: map[string]float64{"Agua Mineral 500ml": 1.0},
		},
		{
			Name: "Cerveza IPA", CategoryID: 4, Price: 3.50,
			Recipe: map[string]float64{"Cerveza Artesanal": 1.0},
		},
	}

	for _, dish := range dishes {
		res, err := db.Exec(
			"INSERT INTO dishes (category_id, name, price, is_active) VALUES (?, ?, ?, ?)",
			dish.CategoryID, dish.Name, dish.Price, true,
		)
		if err != nil {
			continue
		}

		dishID, _ := res.LastInsertId()

		for ingName, qty := range dish.Recipe {
			ingID, exists := ingredientIDs[ingName]
			if exists {
				db.Exec(
					"INSERT INTO dish_ingredients (dish_id, ingredient_id, required_quantity) VALUES (?, ?, ?)",
					dishID, ingID, qty,
				)
			}
		}
	}
	log.Printf("🔹 %d Dishes inserted.", len(dishes))
}
