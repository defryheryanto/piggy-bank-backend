CREATE TABLE IF NOT EXISTS accounts(
    account_id serial PRIMARY KEY,
    account_name VARCHAR (50) NOT NULL,
    account_type_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);