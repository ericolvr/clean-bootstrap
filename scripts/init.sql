-- Criação do banco de dados (se necessário)
-- CREATE DATABASE XYZ;

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mobile VARCHAR(20) NOT NULL UNIQUE,
    user_type INTEGER NOT NULL CHECK (user_type >= 1 AND user_type <= 4),
    password VARCHAR(255),
    status BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for mobile lookups (already unique, but explicit for performance)
CREATE INDEX IF NOT EXISTS idx_users_mobile ON users(mobile);

-- Index for user_type filtering
CREATE INDEX IF NOT EXISTS idx_users_user_type ON users(user_type);

-- Index for status filtering (active/inactive users)
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- Composite index for common queries (status + user_type)
CREATE INDEX IF NOT EXISTS idx_users_status_user_type ON users(status, user_type);

-- Index for name searches (case-insensitive)
CREATE INDEX IF NOT EXISTS idx_users_name_lower ON users(LOWER(name));

-- Index for created_at for sorting/filtering by creation date
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);