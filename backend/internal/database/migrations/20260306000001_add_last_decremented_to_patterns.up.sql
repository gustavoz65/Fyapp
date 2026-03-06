ALTER TABLE category_patterns
ADD COLUMN last_decremented_at DATETIME NULL AFTER confidence,
ADD INDEX idx_last_decremented (last_decremented_at);
