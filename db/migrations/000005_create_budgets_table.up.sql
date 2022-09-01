CREATE TABLE IF NOT EXISTS budgets(
    budget_id SERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL,
    month INT NOT NULL,
    year INT NOT NULL,
    budget BIGINT NOT NULL
);