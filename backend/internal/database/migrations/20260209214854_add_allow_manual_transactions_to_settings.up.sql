-- Add allow_manual_transactions column to user_settings table
ALTER TABLE user_settings
ADD COLUMN allow_manual_transactions BOOLEAN NOT NULL DEFAULT TRUE
AFTER low_balance_threshold;
