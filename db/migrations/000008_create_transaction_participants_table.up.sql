CREATE TABLE transaction_participants (
    participant_id SERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL
);