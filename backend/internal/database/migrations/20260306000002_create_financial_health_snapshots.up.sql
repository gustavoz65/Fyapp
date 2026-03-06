CREATE TABLE financial_health_snapshots (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    period DATE NOT NULL COMMENT 'YYYY-MM-01 (primeiro dia do mês)',
    score DECIMAL(3,2) NOT NULL COMMENT '0.00 a 5.00',
    score_percentage DECIMAL(5,2) NOT NULL COMMENT '0.00 a 100.00',
    economy_rate DECIMAL(3,2) COMMENT '40% do score',
    budget_compliance DECIMAL(3,2) COMMENT '20% do score',
    goals_progress DECIMAL(3,2) COMMENT '20% do score',
    spending_reduction DECIMAL(3,2) COMMENT '10% do score',
    consistency DECIMAL(3,2) COMMENT '10% do score',
    breakdown_json TEXT COMMENT 'Detalhes em JSON',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_period (user_id, period),
    INDEX idx_user_created (user_id, created_at DESC),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
