-- Add source column to transactions table
ALTER TABLE transactions
ADD COLUMN source ENUM('manual', 'bank_sync', 'recurring') NOT NULL DEFAULT 'manual'
AFTER notes;

-- Update existing recurring transactions
UPDATE transactions
SET source = 'recurring'
WHERE is_recurring = TRUE OR recurring_id IS NOT NULL;

-- Create index for source column for better query performance
CREATE INDEX idx_transactions_source ON transactions(source);
