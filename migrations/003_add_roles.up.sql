-- Create roles table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert predefined roles
INSERT INTO roles (name, description) VALUES
    ('client', 'Обычный клиент автомойки - может управлять только своими данными'),
    ('manager', 'Менеджер компании - может управлять своими данными и настройками своей автомойки'),
    ('superuser', 'Суперпользователь - полный доступ ко всем данным и настройкам');

-- Add role_id column to users table
ALTER TABLE users ADD COLUMN role_id INTEGER NOT NULL DEFAULT 1;

-- Add foreign key constraint
ALTER TABLE users ADD CONSTRAINT fk_users_role
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT;

-- Add index on role_id for faster queries
CREATE INDEX idx_users_role_id ON users(role_id);

-- Add comment
COMMENT ON COLUMN users.role_id IS 'User role reference: 1=client, 2=manager, 3=superuser';
