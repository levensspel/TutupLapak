CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(255) UNIQUE,
    password_hash TEXT NOT NULL,
    fileId TEXT,
    fileUri TEXT,
    fileThumbnailUri TEXT,
    bankAccountName VARCHAR(255) CHECK (length(bankAccountName) >= 4 AND length(bankAccountName) <= 32),
    bankAccountHolder VARCHAR(255) CHECK (length(bankAccountHolder) >= 4 AND length(bankAccountHolder) <= 32),
    bankAccountNumber VARCHAR(255) CHECK (length(bankAccountNumber) >= 4 AND length(bankAccountNumber) <= 32),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
