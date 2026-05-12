DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP INDEX IF EXISTS idx_users_is_suspended;
DROP INDEX IF EXISTS idx_users_is_email_verified;
DROP INDEX IF EXISTS idx_users_email;

DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;
