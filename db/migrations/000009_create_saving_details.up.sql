CREATE TABLE saving_details (
    saving_detail_id SERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    target_account_id BIGINT NOT NULL
);