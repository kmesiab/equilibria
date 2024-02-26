-- +goose Up
-- This section is executed when the migration is applied.

CREATE TABLE transactions
(
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    -- 'id' is a unique identifier for each transaction record.

    user_id BIGINT NOT NULL,
    -- 'user_id' links the transaction to a user. It references the 'id' column of the 'users' table.

    conversation_id BIGINT NOT NULL,
    -- 'conversation_id' links the transaction to a specific message.
    -- This is only set for a debit transaction.
    -- The 'cost' of a conversation is the sum of the cost of the messages in a conversation
    -- For instance, a conversation (id:1) consisting of two sms messages billed at a rate of 5 credits per message
    -- would create a transaction record of -10 credits with conversation id:1

    amount DECIMAL(10, 2) NOT NULL,
    -- 'amount' represents the credit amount involved in the transaction.
    -- Decimal type is used to handle amounts with precision up to two decimal places.
    -- Debits are negative numbers so totals are a calculated sum of all transactions.

    transaction_type ENUM ('credit', 'debit') NOT NULL,
    -- 'transaction_type' indicates whether the transaction is a credit (adding funds) or a debit (using funds).

    funding_source ENUM ('stripe', 'paypal', 'bank_transfer', 'cash', 'refund', 'customer credit') NOT NULL,
    -- 'funding_source' specifies the method used to fund the transaction.
    -- It's an ENUM type, meaning the value must be one of the specified options:
    -- 'stripe' indicates the transaction was funded through Stripe.
    -- 'paypal' indicates the transaction was funded through PayPal.
    -- 'bank_transfer' indicates the transaction was funded through a direct bank transfer.
    -- 'cash' indicates the transaction was funded with cash.
    -- 'refund' indicates the transaction was a refund.
    -- 'customer credit' indicates the transaction was funded by the system.
    -- 'customer credit' indicates the transaction was funded by the system.
    -- This field is marked as NOT NULL, ensuring that every transaction record must have a funding source specified.

    description TEXT,
    -- 'description' provides additional details about the transaction.

    reference_id VARCHAR(255),
    -- 'reference_id' can be used to store a reference to an external entity or process related to the transaction.

    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- 'timestamp' records the date and time when the transaction was created.

    FOREIGN KEY (conversation_id) REFERENCES conversations (id),
    FOREIGN KEY (user_id) REFERENCES users (id),

    INDEX idx_user_id (user_id),
    INDEX idx_user_id_conversation_id (user_id, conversation_id),
    INDEX idx_funding_source (funding_source),
    INDEX idx_user_id_timestamp (user_id, timestamp),
    INDEX idx_user_id_transaction_type (user_id, transaction_type)
);

-- +goose Down
-- This section is executed when the migration is rolled back.

DROP TABLE IF EXISTS transactions;
-- This command removes the 'transactions' table if it exists.
