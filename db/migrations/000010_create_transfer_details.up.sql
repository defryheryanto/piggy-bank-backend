CREATE TABLE transfer_details (
    transfer_detail_id SERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    target_account_id BIGINT NOT NULL
);