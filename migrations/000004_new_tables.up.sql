-- 4_up.sql

-- Удаление существующих таблиц, если они существуют
-- DROP TABLE IF EXISTS transactions;
-- DROP TABLE IF EXISTS wallets;

-- Создание новых таблиц с id type VARCHAR
CREATE TABLE wallets
(
    id      VARCHAR PRIMARY KEY,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 100.0
);

CREATE TABLE transactions
(
    id           SERIAL PRIMARY KEY,
    time         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sender_id    VARCHAR        NOT NULL,
    recipient_id VARCHAR        NOT NULL,
    amount       DECIMAL(10, 2) NOT NULL
);
