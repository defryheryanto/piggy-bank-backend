CREATE TABLE IF NOT EXISTS account_types(
    account_type_id serial PRIMARY KEY,
    account_type_name VARCHAR (50) UNIQUE NOT NULL
);