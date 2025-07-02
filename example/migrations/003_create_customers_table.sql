-- +goose Up
-- Create customers table with all necessary indexes
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_customers_username ON customers(username);
CREATE INDEX idx_customers_email ON customers(email);

-- Insert sample data
INSERT INTO customers (username, email) VALUES
('john_doe', 'john.doe@example.com'),
('jane_smith', 'jane.smith@example.com'),
('bob_wilson', 'bob.wilson@example.com'),
('alice_johnson', 'alice.johnson@example.com');

-- +goose Down
-- Drop customers table
DROP TABLE IF EXISTS customers;
