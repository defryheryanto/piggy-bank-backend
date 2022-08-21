CREATE TABLE IF NOT EXISTS categories(
    category_id SERIAL PRIMARY KEY,
    category_name VARCHAR(50) NOT NULL,
    category_type VARCHAR(20),
    user_id BIGINT NOT NULL,
    budget BIGINT DEFAULT 0
);