-- ============================================================
-- INITIAL MIGRATION
-- Author: Adán de la Rosa Lugo
-- Date: 2024-06-10
-- ============================================================

CREATE DATABASE IF NOT EXISTS mep_db
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE mep_db;


-- ============================================================
-- 1. USERS
-- ============================================================

CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;


-- ============================================================
-- 2. ROLES
-- ============================================================

CREATE TABLE IF NOT EXISTS roles (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
) ENGINE=InnoDB;

-- ============================================================
-- 3. USER_ROLES
-- ============================================================

CREATE TABLE IF NOT EXISTS user_roles (
    user_id INT UNSIGNED NOT NULL,
    role_id INT UNSIGNED NOT NULL,
    assigned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, role_id),

    CONSTRAINT fk_user_roles_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_user_roles_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;

-- ============================================================
-- 4. DINING_TABLES
-- ============================================================

CREATE TABLE IF NOT EXISTS dining_tables (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    table_number INT UNSIGNED NOT NULL UNIQUE,
    capacity INT UNSIGNED NOT NULL,
    status ENUM(
        'DISPONIBLE',
        'OCUPADA',
        'RESERVADA',
        'LIMPIEZA',
        'FUERA_DE_SERVICIO'
    ) NOT NULL DEFAULT 'DISPONIBLE',

    CONSTRAINT chk_dining_tables_capacity
        CHECK (capacity > 0),

    CONSTRAINT chk_dining_tables_number
        CHECK (table_number > 0)
) ENGINE=InnoDB;

-- ============================================================
-- 5. CATEGORIES
-- ============================================================

CREATE TABLE IF NOT EXISTS categories (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
) ENGINE=InnoDB;

-- ============================================================
-- 6. DISHES
-- ============================================================

CREATE TABLE IF NOT EXISTS dishes (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    category_id INT UNSIGNED NOT NULL,
    name VARCHAR(150) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT chk_dishes_price
        CHECK (price >= 0),

    CONSTRAINT fk_dishes_category
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
) ENGINE=InnoDB;

-- ============================================================
-- 7. INGREDIENTS
-- ============================================================

CREATE TABLE IF NOT EXISTS ingredients (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    unit_of_measure VARCHAR(20) NOT NULL,
    current_stock DECIMAL(12,3) NOT NULL DEFAULT 0.000,
    minimum_stock DECIMAL(12,3) NOT NULL DEFAULT 0.000,

    CONSTRAINT chk_ingredients_current_stock
        CHECK (current_stock >= 0),

    CONSTRAINT chk_ingredients_minimum_stock
        CHECK (minimum_stock >= 0)
) ENGINE=InnoDB;

-- ============================================================
-- 8. DISH_INGREDIENTS
-- ============================================================

CREATE TABLE IF NOT EXISTS dish_ingredients (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    dish_id INT UNSIGNED NOT NULL,
    ingredient_id INT UNSIGNED NOT NULL,
    required_quantity DECIMAL(12,3) NOT NULL,

    CONSTRAINT uq_dish_ingredient
        UNIQUE (dish_id, ingredient_id),

    CONSTRAINT chk_dish_ingredients_quantity
        CHECK (required_quantity > 0),

    CONSTRAINT fk_dish_ingredients_dish
        FOREIGN KEY (dish_id)
        REFERENCES dishes(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_dish_ingredients_ingredient
        FOREIGN KEY (ingredient_id)
        REFERENCES ingredients(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
) ENGINE=InnoDB;

-- ============================================================
-- 9. ORDERS
-- ============================================================

CREATE TABLE IF NOT EXISTS orders (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    dining_table_id INT UNSIGNED NOT NULL,
    opened_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    closed_at DATETIME NULL,
    status ENUM(
        'ABIERTO',
        'EN_PREPARACION',
        'LISTO',
        'CERRADO',
        'CANCELADO'
    ) NOT NULL DEFAULT 'ABIERTO',
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,

    CONSTRAINT chk_orders_total
        CHECK (total_amount >= 0),

    CONSTRAINT fk_orders_dining_table
        FOREIGN KEY (dining_table_id)
        REFERENCES dining_tables(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
) ENGINE=InnoDB;


-- ============================================================
-- 10. ORDER_ITEMS
-- ============================================================

CREATE TABLE IF NOT EXISTS order_items (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id INT UNSIGNED NOT NULL,
    dish_id INT UNSIGNED NOT NULL,
    quantity INT UNSIGNED NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    kitchen_status ENUM(
        'PENDIENTE',
        'PREPARANDO',
        'LISTO',
        'SERVIDO',
        'CANCELADO'
    ) NOT NULL DEFAULT 'PENDIENTE',
    notes TEXT NULL,

    CONSTRAINT chk_order_items_quantity
        CHECK (quantity > 0),

    CONSTRAINT chk_order_items_unit_price
        CHECK (unit_price >= 0),

    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_order_items_dish
        FOREIGN KEY (dish_id)
        REFERENCES dishes(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
) ENGINE=InnoDB;


-- ============================================================
-- 11. ALLERGENS
-- ============================================================

CREATE TABLE IF NOT EXISTS allergens (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    icon VARCHAR(100) NULL
) ENGINE=InnoDB;


-- ============================================================
-- 12. INGREDIENT_ALLERGENS
-- ============================================================

CREATE TABLE IF NOT EXISTS ingredient_allergens (
    ingredient_id INT UNSIGNED NOT NULL,
    allergen_id INT UNSIGNED NOT NULL,

    PRIMARY KEY (ingredient_id, allergen_id),

    CONSTRAINT fk_ingredient_allergens_ingredient
        FOREIGN KEY (ingredient_id)
        REFERENCES ingredients(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_ingredient_allergens_allergen
        FOREIGN KEY (allergen_id)
        REFERENCES allergens(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;


-- ============================================================
-- 13. DISH_ALLERGEN_TRACES
-- ============================================================

-- CREATE TABLE IF NOT EXISTS dish_allergen_traces (
--     dish_id INT UNSIGNED NOT NULL,
--     allergen_id INT UNSIGNED NOT NULL,
--     reason VARCHAR(255) NULL,

--     PRIMARY KEY (dish_id, allergen_id),

--     CONSTRAINT fk_dish_allergen_traces_dish
--         FOREIGN KEY (dish_id)
--         REFERENCES dishes(id)
--         ON DELETE CASCADE
--         ON UPDATE CASCADE,

--     CONSTRAINT fk_dish_allergen_traces_allergen
--         FOREIGN KEY (allergen_id)
--         REFERENCES allergens(id)
--         ON DELETE CASCADE
--         ON UPDATE CASCADE
-- ) ENGINE=InnoDB;


-- ============================================================
-- 14. INVENTORY_MOVEMENTS
-- ============================================================

CREATE TABLE IF NOT EXISTS inventory_movements (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    ingredient_id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NULL,

    movement_type ENUM(
        'COMPRA',
        'CONSUMO',
        'DESPERDICIO',
        'AJUSTE',
        'DEVOLUCION'
    ) NOT NULL,

    /*
      quantity:
      > 0  aumenta el stock
      < 0  reduce el stock
    */
    quantity DECIMAL(12,3) NOT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_inventory_movements_quantity
        CHECK (quantity <> 0),

    CONSTRAINT fk_inventory_movements_ingredient
        FOREIGN KEY (ingredient_id)
        REFERENCES ingredients(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,

    CONSTRAINT fk_inventory_movements_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE
) ENGINE=InnoDB;


-- ============================================================
-- ÍNDICES ADICIONALES
-- ============================================================

CREATE INDEX idx_orders_dining_table
    ON orders(dining_table_id);

CREATE INDEX idx_order_items_order
    ON order_items(order_id);

CREATE INDEX idx_order_items_dish
    ON order_items(dish_id);

CREATE INDEX idx_dish_ingredients_dish
    ON dish_ingredients(dish_id);

CREATE INDEX idx_dish_ingredients_ingredient
    ON dish_ingredients(ingredient_id);

CREATE INDEX idx_inventory_movements_ingredient
    ON inventory_movements(ingredient_id);

CREATE INDEX idx_inventory_movements_user
    ON inventory_movements(user_id);

CREATE INDEX idx_inventory_movements_created_at
    ON inventory_movements(created_at);

