INSERT INTO users (email, first_name, last_name, birth_date, password_hash, is_email_verified, is_document_verified)
VALUES (
    'root@admin.com',
    'Root',
    'Admin',
    '1990-01-01',
    '$2a$10$vLrqa5Ir8C2T8ikms3.sJ.9gsA7bK0NnFc9.bjmNKhw3IwwXQ2ScG'::bytea,
    true,
    true
);

INSERT INTO user_roles (user_id, role_name)
SELECT id, 'admin' FROM users WHERE email = 'root@admin.com';
