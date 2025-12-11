-- Remove index
DROP INDEX IF EXISTS idx_users_role_id;

-- Remove foreign key constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_role;

-- Remove role_id column from users
ALTER TABLE users DROP COLUMN IF EXISTS role_id;

-- Drop roles table
DROP TABLE IF EXISTS roles;
