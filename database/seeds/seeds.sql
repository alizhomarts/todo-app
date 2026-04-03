-- +goose Up
INSERT INTO users (email, first_name, last_name, password_hash)
-- pass: 123456
VALUES
    (
        'admin@test.com',
        'Admin',
        'Admin',
        '$2a$12$gtBak6nNlgWKFJjIJCmSpOsoSQ98Fa.L2F6WWs5rcy0AZr.stfOJ.'
    ),
    (
        'user@test.com',
        'User',
        'User',
        '$2a$12$gtBak6nNlgWKFJjIJCmSpOsoSQ98Fa.L2F6WWs5rcy0AZr.stfOJ.'
    );

-- +goose Down
-- DELETE FROM users
-- WHERE email IN ('admin@test.com', 'user@test.com');