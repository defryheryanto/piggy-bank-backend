CREATE TABLE IF NOT EXISTS user_configs(
    config_id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    monthly_start_date INT DEFAULT 1
);