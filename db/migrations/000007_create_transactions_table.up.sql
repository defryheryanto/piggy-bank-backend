CREATE TABLE transactions(
    transaction_id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    category_id BIGINT,
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_DATE,
    description TEXT,
    notes TEXT,
    amount DECIMAL(10, 2) NOT NULL,
    transaction_type VARCHAR(20) NOT NULL
);