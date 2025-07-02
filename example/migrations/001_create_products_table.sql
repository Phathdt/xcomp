-- +goose Up
-- Create products table with all necessary indexes and triggers
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    category VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_is_active ON products(is_active);

-- Note: Function and trigger creation moved to separate steps for Goose compatibility

INSERT INTO products (name, description, price, stock_quantity, category) VALUES
('Laptop Pro', 'High-performance laptop for professionals', 1299.99, 10, 'Electronics'),
('Wireless Mouse', 'Ergonomic wireless mouse with precision tracking', 29.99, 50, 'Electronics'),
('Coffee Mug', 'Ceramic coffee mug with thermal insulation', 15.99, 100, 'Home & Kitchen'),
('Notebook Set', 'Premium notebook set for writing and sketching', 24.99, 25, 'Stationery');

-- +goose Down
-- Drop products table
DROP TABLE IF EXISTS products;
