
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
    wallet_name varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS linked_accounts (
    id varchar(255) PRIMARY KEY,
    user_id varchar(255) NOT NULL,
    account_name varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id varchar(255) PRIMARY KEY,
    wallet_id varchar(255) NOT NULL,
    account_id varchar(255) NOT NULL,
    amount decimal(10, 2) NOT NULL,
    currency varchar(10) NOT NULL DEFAULT 'VND',
    transaction_kind varchar(100) NOT NULL,
    status varchar(100) NOT NULL,
    note text,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE wallets ADD CONSTRAINT fk_wallet_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE linked_accounts ADD CONSTRAINT fk_linked_account_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE transactions ADD CONSTRAINT fk_trans_wallet_id FOREIGN KEY (wallet_id) REFERENCES wallets(id);
ALTER TABLE transactions ADD CONSTRAINT fk_trans_account_id FOREIGN KEY (account_id) REFERENCES linked_accounts(id);

-- +migrate Down
