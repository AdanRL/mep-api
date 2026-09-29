-- ============================================================
-- REVERSION INITIAL MIGRATION (DOWN)
-- Author: Adán de la Rosa Lugo
-- Date: 2024-06-10
-- ============================================================
USE mep_db;

DROP TABLE IF EXISTS inventory_movements;
DROP TABLE IF EXISTS ingredient_allergens;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS dish_ingredients;
DROP TABLE IF EXISTS allergens;
DROP TABLE IF EXISTS dishes;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS dining_tables;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
