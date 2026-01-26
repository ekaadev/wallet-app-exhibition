CREATE TABLE wallet_mutations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    wallet_id BIGINT UNSIGNED NOT NULL,
    transaction_id BIGINT UNSIGNED NOT NULL,
    type ENUM('debit', 'credit') NOT NULL,
    amount DECIMAL(20, 2) NOT NULL,
    balance_before DECIMAL(20, 2) NOT NULL,
    balance_after DECIMAL(20, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_wallet_mutations_wallet_id (wallet_id),
    INDEX idx_wallet_mutations_transaction_id (transaction_id),
    INDEX idx_wallet_mutations_type (type),
    CONSTRAINT fk_wallet_mutations_wallet_id FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_wallet_mutations_transaction_id FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
