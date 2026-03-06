ALTER TABLE category_patterns
DROP INDEX idx_last_decremented,
DROP COLUMN last_decremented_at;
