CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login TEXT NOT NULL UNIQUE,                                                 
    password_hash BYTEA NOT NULL   
);

CREATE TABLE IF NOT EXISTS ads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(100) NOT NULL,
    description TEXT,
    price INT NOT NULL,
    image_url TEXT DEFAULT 'default_product.jpg',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);