-- Remove allow_manual_transactions column from user_settings table
ALTER TABLE user_settings
DROP COLUMN allow_manual_transactions;
