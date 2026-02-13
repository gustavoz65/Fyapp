-- =====================================================
-- USERS TABLE
-- =====================================================
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NULL,
    avatar_url VARCHAR(500) NULL,
    preferred_currency VARCHAR(3) DEFAULT 'BRL',
    preferred_language VARCHAR(5) DEFAULT 'pt-BR',
    timezone VARCHAR(50) DEFAULT 'America/Sao_Paulo',
    role ENUM('user', 'admin', 'premium') DEFAULT 'user',
    email_verified BOOLEAN DEFAULT FALSE,
    email_verified_at DATETIME NULL,
    last_login_at DATETIME NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_email (email),
    INDEX idx_users_role (role),
    INDEX idx_users_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- USER SESSIONS TABLE (for JWT refresh tokens)
-- =====================================================
CREATE TABLE user_sessions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    device_info VARCHAR(500) NULL,
    ip_address VARCHAR(45) NULL,
    user_agent VARCHAR(500) NULL,
    expires_at DATETIME NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_sessions_user_id (user_id),
    INDEX idx_user_sessions_expires_at (expires_at),
    INDEX idx_user_sessions_refresh_token (refresh_token_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- CATEGORIES TABLE
-- =====================================================
CREATE TABLE categories (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255) NULL,
    type ENUM('income', 'expense') NOT NULL,
    color VARCHAR(7) DEFAULT '#6366F1',
    icon VARCHAR(50) DEFAULT 'wallet',
    is_system BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_category (user_id, name, type),
    INDEX idx_categories_user_id (user_id),
    INDEX idx_categories_type (type),
    INDEX idx_categories_is_system (is_system)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- BANK ACCOUNTS TABLE
-- =====================================================
CREATE TABLE bank_accounts (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    bank_name VARCHAR(100) NULL,
    bank_code VARCHAR(10) NULL,
    account_type ENUM('checking', 'savings', 'credit_card', 'investment', 'cash', 'other') NOT NULL,
    account_number VARCHAR(50) NULL,
    agency VARCHAR(20) NULL,
    initial_balance DECIMAL(15,2) DEFAULT 0.00,
    current_balance DECIMAL(15,2) DEFAULT 0.00,
    credit_limit DECIMAL(15,2) NULL,
    closing_day TINYINT UNSIGNED NULL CHECK (closing_day >= 1 AND closing_day <= 31),
    due_day TINYINT UNSIGNED NULL CHECK (due_day >= 1 AND due_day <= 31),
    currency VARCHAR(3) DEFAULT 'BRL',
    color VARCHAR(7) DEFAULT '#10B981',
    icon VARCHAR(50) DEFAULT 'bank',
    is_active BOOLEAN DEFAULT TRUE,
    include_in_total BOOLEAN DEFAULT TRUE,
    last_sync_at DATETIME NULL,
    external_id VARCHAR(255) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_bank_accounts_user_id (user_id),
    INDEX idx_bank_accounts_account_type (account_type),
    INDEX idx_bank_accounts_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- TRANSACTIONS TABLE
-- =====================================================
CREATE TABLE transactions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    bank_account_id CHAR(36) NOT NULL,
    category_id CHAR(36) NULL,
    type ENUM('income', 'expense', 'transfer') NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    description VARCHAR(255) NOT NULL,
    notes TEXT NULL,
    transaction_date DATE NOT NULL,
    due_date DATE NULL,
    payment_date DATE NULL,
    is_paid BOOLEAN DEFAULT FALSE,
    is_recurring BOOLEAN DEFAULT FALSE,
    recurring_id CHAR(36) NULL,
    installment_number INT UNSIGNED NULL,
    total_installments INT UNSIGNED NULL,
    installment_group_id CHAR(36) NULL,
    tags JSON NULL,
    attachment_url VARCHAR(500) NULL,
    external_id VARCHAR(255) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    INDEX idx_transactions_user_id (user_id),
    INDEX idx_transactions_bank_account_id (bank_account_id),
    INDEX idx_transactions_category_id (category_id),
    INDEX idx_transactions_type (type),
    INDEX idx_transactions_transaction_date (transaction_date),
    INDEX idx_transactions_due_date (due_date),
    INDEX idx_transactions_is_paid (is_paid),
    INDEX idx_transactions_recurring_id (recurring_id),
    INDEX idx_transactions_installment_group_id (installment_group_id),
    INDEX idx_transactions_created_at (created_at),
    INDEX idx_transactions_user_date (user_id, transaction_date),
    INDEX idx_transactions_user_type_date (user_id, type, transaction_date),
    INDEX idx_transactions_account_date (bank_account_id, transaction_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- RECURRING TRANSACTIONS TABLE
-- =====================================================
CREATE TABLE recurring_transactions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    bank_account_id CHAR(36) NOT NULL,
    category_id CHAR(36) NULL,
    type ENUM('income', 'expense') NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    description VARCHAR(255) NOT NULL,
    frequency ENUM('daily', 'weekly', 'biweekly', 'monthly', 'quarterly', 'yearly') NOT NULL,
    day_of_month TINYINT UNSIGNED NULL CHECK (day_of_month >= 1 AND day_of_month <= 31),
    day_of_week TINYINT UNSIGNED NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    start_date DATE NOT NULL,
    end_date DATE NULL,
    next_occurrence DATE NOT NULL,
    last_generated_at DATE NULL,
    is_active BOOLEAN DEFAULT TRUE,
    auto_confirm BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    INDEX idx_recurring_transactions_user_id (user_id),
    INDEX idx_recurring_transactions_next_occurrence (next_occurrence),
    INDEX idx_recurring_transactions_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- BUDGETS TABLE
-- =====================================================
CREATE TABLE budgets (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    category_id CHAR(36) NULL,
    name VARCHAR(100) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    spent_amount DECIMAL(15,2) DEFAULT 0.00,
    period_type ENUM('monthly', 'quarterly', 'yearly', 'custom') NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    alert_threshold DECIMAL(5,2) DEFAULT 80.00 CHECK (alert_threshold > 0 AND alert_threshold <= 100),
    alert_sent BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    CHECK (end_date > start_date),
    INDEX idx_budgets_user_id (user_id),
    INDEX idx_budgets_category_id (category_id),
    INDEX idx_budgets_period (start_date, end_date),
    INDEX idx_budgets_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- GOALS TABLE (Financial Goals)
-- =====================================================
CREATE TABLE goals (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    target_amount DECIMAL(15,2) NOT NULL CHECK (target_amount > 0),
    current_amount DECIMAL(15,2) DEFAULT 0.00,
    target_date DATE NULL,
    icon VARCHAR(50) DEFAULT 'target',
    color VARCHAR(7) DEFAULT '#F59E0B',
    priority TINYINT UNSIGNED DEFAULT 1 CHECK (priority >= 1 AND priority <= 5),
    status ENUM('in_progress', 'completed', 'cancelled') DEFAULT 'in_progress',
    completed_at DATETIME NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_goals_user_id (user_id),
    INDEX idx_goals_status (status),
    INDEX idx_goals_target_date (target_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- GOAL CONTRIBUTIONS TABLE
-- =====================================================
CREATE TABLE goal_contributions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    goal_id CHAR(36) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    note VARCHAR(255) NULL,
    contribution_date DATE NOT NULL DEFAULT (CURRENT_DATE),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (goal_id) REFERENCES goals(id) ON DELETE CASCADE,
    INDEX idx_goal_contributions_goal_id (goal_id),
    INDEX idx_goal_contributions_date (contribution_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- NOTIFICATIONS TABLE
-- =====================================================
CREATE TABLE notifications (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    type ENUM('budget_alert', 'bill_reminder', 'goal_achieved', 'low_balance', 'transaction_alert', 'system') NOT NULL,
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    data JSON NULL,
    is_read BOOLEAN DEFAULT FALSE,
    read_at DATETIME NULL,
    sent_via JSON DEFAULT ('["in_app"]'),
    scheduled_for DATETIME NULL,
    sent_at DATETIME NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_notifications_user_id (user_id),
    INDEX idx_notifications_is_read (is_read),
    INDEX idx_notifications_type (type),
    INDEX idx_notifications_created_at (created_at),
    INDEX idx_notifications_scheduled (scheduled_for)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- USER SETTINGS TABLE
-- =====================================================
CREATE TABLE user_settings (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL UNIQUE,
    notification_email BOOLEAN DEFAULT TRUE,
    notification_push BOOLEAN DEFAULT TRUE,
    notification_sms BOOLEAN DEFAULT FALSE,
    budget_alerts BOOLEAN DEFAULT TRUE,
    bill_reminders BOOLEAN DEFAULT TRUE,
    bill_reminder_days INT DEFAULT 3,
    weekly_summary BOOLEAN DEFAULT TRUE,
    monthly_report BOOLEAN DEFAULT TRUE,
    low_balance_alert BOOLEAN DEFAULT TRUE,
    low_balance_threshold DECIMAL(15,2) DEFAULT 100.00,
    theme ENUM('light', 'dark', 'system') DEFAULT 'system',
    dashboard_layout JSON NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_settings_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- BANK INTEGRATIONS TABLE (for Open Banking)
-- =====================================================
CREATE TABLE bank_integrations (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    bank_account_id CHAR(36) NULL,
    provider VARCHAR(50) NOT NULL,
    provider_account_id VARCHAR(255) NULL,
    access_token_encrypted TEXT NULL,
    refresh_token_encrypted TEXT NULL,
    token_expires_at DATETIME NULL,
    consent_expires_at DATETIME NULL,
    status ENUM('pending', 'active', 'expired', 'revoked', 'error') DEFAULT 'pending',
    last_sync_at DATETIME NULL,
    last_error TEXT NULL,
    metadata JSON NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE SET NULL,
    INDEX idx_bank_integrations_user_id (user_id),
    INDEX idx_bank_integrations_status (status),
    INDEX idx_bank_integrations_provider (provider)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- AUDIT LOG TABLE
-- =====================================================
CREATE TABLE audit_logs (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NULL,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id CHAR(36) NULL,
    old_values JSON NULL,
    new_values JSON NULL,
    ip_address VARCHAR(45) NULL,
    user_agent VARCHAR(500) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_audit_logs_user_id (user_id),
    INDEX idx_audit_logs_action (action),
    INDEX idx_audit_logs_entity (entity_type, entity_id),
    INDEX idx_audit_logs_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- REPORTS TABLE (for scheduled/saved reports)
-- =====================================================
CREATE TABLE reports (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    type ENUM('cash_flow', 'expense_by_category', 'income_vs_expense', 'budget_analysis', 'custom') NOT NULL,
    parameters JSON NOT NULL,
    schedule ENUM('daily', 'weekly', 'monthly') NULL,
    last_generated_at DATETIME NULL,
    file_url VARCHAR(500) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_reports_user_id (user_id),
    INDEX idx_reports_type (type),
    INDEX idx_reports_schedule (schedule)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- TRANSFERS TABLE (for transfers between accounts)
-- =====================================================
CREATE TABLE transfers (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    from_account_id CHAR(36) NOT NULL,
    to_account_id CHAR(36) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    description VARCHAR(255) NULL,
    transfer_date DATE NOT NULL,
    from_transaction_id CHAR(36) NULL,
    to_transaction_id CHAR(36) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (from_account_id) REFERENCES bank_accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (to_account_id) REFERENCES bank_accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (from_transaction_id) REFERENCES transactions(id) ON DELETE SET NULL,
    FOREIGN KEY (to_transaction_id) REFERENCES transactions(id) ON DELETE SET NULL,
    CHECK (from_account_id != to_account_id),
    INDEX idx_transfers_user_id (user_id),
    INDEX idx_transfers_from_account (from_account_id),
    INDEX idx_transfers_to_account (to_account_id),
    INDEX idx_transfers_date (transfer_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- INSTALLMENTS TABLE (for installment tracking)
-- =====================================================
CREATE TABLE installments (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    bank_account_id CHAR(36) NOT NULL,
    category_id CHAR(36) NULL,
    description VARCHAR(255) NOT NULL,
    total_amount DECIMAL(15,2) NOT NULL CHECK (total_amount > 0),
    installment_amount DECIMAL(15,2) NOT NULL CHECK (installment_amount > 0),
    total_installments INT UNSIGNED NOT NULL CHECK (total_installments > 0),
    paid_installments INT UNSIGNED DEFAULT 0,
    first_due_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    INDEX idx_installments_user_id (user_id),
    INDEX idx_installments_bank_account_id (bank_account_id),
    INDEX idx_installments_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- PASSWORD RESET TOKENS TABLE
-- =====================================================
CREATE TABLE password_reset_tokens (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at DATETIME NOT NULL,
    used_at DATETIME NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_password_reset_user_id (user_id),
    INDEX idx_password_reset_token (token_hash),
    INDEX idx_password_reset_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- EMAIL VERIFICATION TOKENS TABLE
-- =====================================================
CREATE TABLE email_verification_tokens (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at DATETIME NOT NULL,
    used_at DATETIME NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_email_verification_user_id (user_id),
    INDEX idx_email_verification_token (token_hash),
    INDEX idx_email_verification_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================================
-- DEFAULT SYSTEM CATEGORIES
-- =====================================================
INSERT INTO categories (id, user_id, name, description, type, color, icon, is_system) VALUES
-- Expense categories
(UUID(), NULL, 'Alimentacao', 'Gastos com alimentacao e restaurantes', 'expense', '#EF4444', '🍴', TRUE),
(UUID(), NULL, 'Transporte', 'Gastos com transporte e combustivel', 'expense', '#F59E0B', '🚗', TRUE),
(UUID(), NULL, 'Moradia', 'Aluguel, condominio e contas da casa', 'expense', '#10B981', '🏠', TRUE),
(UUID(), NULL, 'Saude', 'Consultas, remedios e plano de saude', 'expense', '#EC4899', '❤️', TRUE),
(UUID(), NULL, 'Educacao', 'Cursos, livros e materiais', 'expense', '#8B5CF6', '📚', TRUE),
(UUID(), NULL, 'Lazer', 'Entretenimento e diversao', 'expense', '#06B6D4', '🎮', TRUE),
(UUID(), NULL, 'Compras', 'Compras diversas', 'expense', '#F97316', '🛍️', TRUE),
(UUID(), NULL, 'Servicos', 'Assinaturas e servicos', 'expense', '#6366F1', '⚙️', TRUE),
(UUID(), NULL, 'Outros Despesas', 'Despesas diversas', 'expense', '#6B7280', '➕', TRUE),
-- Income categories
(UUID(), NULL, 'Salario', 'Salario e proventos', 'income', '#22C55E', '💼', TRUE),
(UUID(), NULL, 'Freelance', 'Trabalhos freelance', 'income', '#14B8A6', '💻', TRUE),
(UUID(), NULL, 'Investimentos', 'Rendimentos de investimentos', 'income', '#3B82F6', '📈', TRUE),
(UUID(), NULL, 'Presente', 'Dinheiro recebido de presente', 'income', '#A855F7', '🎁', TRUE),
(UUID(), NULL, 'Reembolso', 'Reembolsos recebidos', 'income', '#84CC16', '🔄', TRUE),
(UUID(), NULL, 'Outras Receitas', 'Outras receitas', 'income', '#6B7280', '➕', TRUE);

