-- Drop index
DROP INDEX IF EXISTS idx_transactions_source ON transactions;

-- Remove source column from transactions table
ALTER TABLE transactions
DROP COLUMN source;
