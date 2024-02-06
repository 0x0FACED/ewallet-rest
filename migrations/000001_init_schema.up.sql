CREATE TABLE wallets
(
    id      uuid PRIMARY KEY,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 100.0
);

CREATE TABLE transactions
(
    id           SERIAL PRIMARY KEY,
    time         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sender_id    uuid            NOT NULL,
    recipient_id uuid            NOT NULL,
    amount       DECIMAL(10, 2) NOT NULL
);