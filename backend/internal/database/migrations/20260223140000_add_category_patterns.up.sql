CREATE TABLE IF NOT EXISTS category_patterns (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    keyword VARCHAR(255) NOT NULL,
    category_id VARCHAR(36) NOT NULL,
    confidence INT NOT NULL DEFAULT 1,
    source VARCHAR(20) NOT NULL DEFAULT 'manual',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_keyword_category (user_id, keyword, category_id),
    INDEX idx_user_keyword (user_id, keyword),
    INDEX idx_user_confidence (user_id, confidence DESC),
    CONSTRAINT fk_category_patterns_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_category_patterns_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
