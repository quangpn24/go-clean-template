
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id varchar(255) PRIMARY KEY,
    full_name varchar(255) NOT NULL,
    email varchar(50) NOT NULL,
    phone_number varchar(20) NOT NULL,
    current_address text NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wallets (
    id varchar(255) PRIMARY KEY,
    user_id varchar(255) NOT NULL,
    balance decimal(10, 2) NOT NULL DEFAULT 0.00,
    currency varchar(10) NOT NULL DEFAULT 'VND',
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS accounts (
    id varchar(255) PRIMARY KEY,
    user_id varchar(255) NOT NULL,
    bank_name varchar(255) NOT NULL,
    account_number varchar(50) NOT NULL,
    is_linked boolean NOT NULL DEFAULT false,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id varchar(255) PRIMARY KEY,
    sender_wallet_id varchar(255),
    receiver_wallet_id varchar(255),
    account_id varchar(255),
    amount decimal(10, 2) NOT NULL,
    currency varchar(10) NOT NULL DEFAULT 'VND',
    category varchar(100) NOT NULL,
    transaction_kind varchar(100) NOT NULL,
    note text,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE wallets ADD CONSTRAINT fk_wallet_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE accounts ADD CONSTRAINT fk_account_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE transactions ADD CONSTRAINT fk_trans_sender_wallet_id FOREIGN KEY (sender_wallet_id) REFERENCES wallets(id);
ALTER TABLE transactions ADD CONSTRAINT fk_trans_receiver_wallet_id FOREIGN KEY (receiver_wallet_id) REFERENCES wallets(id);
ALTER TABLE transactions ADD CONSTRAINT fk_trans_account_id FOREIGN KEY (account_id) REFERENCES accounts(id);

-- +migrate Down
