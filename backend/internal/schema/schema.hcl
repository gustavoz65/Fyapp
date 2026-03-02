// Schema definition for Cashing-go database
// This file defines the database structure in Atlas HCL format
// When you change this file, run: task migrations:new name=your_change

// Define the database schema (uses the database name from atlas.hcl)
schema "cashing" {
}

 ====================================
// USERS TABLE
 ====================================
table "users" {
  schema = schema.cashing
  
  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "email" {
    type = varchar(255)
    null = false
  }
  column "password_hash" {
    type = varchar(255)
    null = false
  }
  column "first_name" {
    type = varchar(100)
    null = false
  }
  column "last_name" {
    type = varchar(100)
    null = false
  }
  column "phone" {
    type = varchar(20)
    null = true
  }
  column "avatar_url" {
    type = varchar(500)
    null = true
  }
  column "preferred_currency" {
    type    = varchar(3)
    default = "BRL"
  }
  column "preferred_language" {
    type    = varchar(5)
    default = "pt-BR"
  }
  column "timezone" {
    type    = varchar(50)
    default = "America/Sao_Paulo"
  }
  column "role" {
    type    = enum("user", "admin", "premium")
    default = "user"
  }
  column "email_verified" {
    type    = bool
    default = false
  }
  column "email_verified_at" {
    type = datetime
    null = true
  }
  column "last_login_at" {
    type = datetime
    null = true
  }
  column "is_active" {
    type    = bool
    default = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  index "idx_users_email" {
    columns = [column.email]
  }
  index "idx_users_role" {
    columns = [column.role]
  }
  index "idx_users_is_active" {
    columns = [column.is_active]
  }
}

 ====================================
// USER SESSIONS TABLE (for JWT refresh tokens)
 ====================================
table "user_sessions" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "refresh_token_hash" {
    type = varchar(255)
    null = false
  }
  column "device_info" {
    type = varchar(500)
    null = true
  }
  column "ip_address" {
    type = varchar(45)
    null = true
  }
  column "user_agent" {
    type = varchar(500)
    null = true
  }
  column "expires_at" {
    type = datetime
    null = false
  }
  column "is_revoked" {
    type    = bool
    default = false
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_user_sessions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  index "idx_user_sessions_user_id" {
    columns = [column.user_id]
  }
  index "idx_user_sessions_expires_at" {
    columns = [column.expires_at]
  }
  index "idx_user_sessions_refresh_token" {
    columns = [column.refresh_token_hash]
  }
}

 ====================================
// CATEGORIES TABLE
 ====================================
table "categories" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = true
  }
  column "name" {
    type = varchar(100)
    null = false
  }
  column "description" {
    type = varchar(255)
    null = true
  }
  column "type" {
    type = enum("income", "expense")
    null = false
  }
  column "color" {
    type    = varchar(7)
    default = "#6366F1"
  }
  column "icon" {
    type    = varchar(50)
    default = "wallet"
  }
  column "is_system" {
    type    = bool
    default = false
  }
  column "is_active" {
    type    = bool
    default = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_categories_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  index "unique_user_category" {
    unique  = true
    columns = [column.user_id, column.name, column.type]
  }

  index "idx_categories_user_id" {
    columns = [column.user_id]
  }
  index "idx_categories_type" {
    columns = [column.type]
  }
  index "idx_categories_is_system" {
    columns = [column.is_system]
  }
}

 ====================================
// BANK ACCOUNTS TABLE
 ====================================
table "bank_accounts" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "name" {
    type = varchar(100)
    null = false
  }
  column "bank_name" {
    type = varchar(100)
    null = true
  }
  column "bank_code" {
    type = varchar(10)
    null = true
  }
  column "account_type" {
    type = enum("checking", "savings", "credit_card", "investment", "cash", "other")
    null = false
  }
  column "account_number" {
    type = varchar(50)
    null = true
  }
  column "agency" {
    type = varchar(20)
    null = true
  }
  column "initial_balance" {
    type    = decimal(15, 2)
    default = 0.00
  }
  column "current_balance" {
    type    = decimal(15, 2)
    default = 0.00
  }
  column "credit_limit" {
    type = decimal(15, 2)
    null = true
  }
  column "closing_day" {
    type = tinyint
    null = true
  }
  column "due_day" {
    type = tinyint
    null = true
  }
  column "currency" {
    type    = varchar(3)
    default = "BRL"
  }
  column "color" {
    type    = varchar(7)
    default = "#10B981"
  }
  column "icon" {
    type    = varchar(50)
    default = "bank"
  }
  column "is_active" {
    type    = bool
    default = true
  }
  column "include_in_total" {
    type    = bool
    default = true
  }
  column "last_sync_at" {
    type = datetime
    null = true
  }
  column "external_id" {
    type = varchar(255)
    null = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_bank_accounts_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  index "idx_bank_accounts_user_id" {
    columns = [column.user_id]
  }
  index "idx_bank_accounts_account_type" {
    columns = [column.account_type]
  }
  index "idx_bank_accounts_is_active" {
    columns = [column.is_active]
  }
}

 ====================================
// TRANSACTIONS TABLE
 ====================================
table "transactions" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "bank_account_id" {
    type = char(36)
    null = false
  }
  column "category_id" {
    type = char(36)
    null = true
  }
  column "type" {
    type = enum("income", "expense", "transfer")
    null = false
  }
  column "amount" {
    type = decimal(15, 2)
    null = false
  }
  column "description" {
    type = varchar(255)
    null = false
  }
  column "notes" {
    type = text
    null = true
  }
  column "source" {
    type    = enum("manual", "bank_sync", "recurring")
    default = "manual"
  }
  column "transaction_date" {
    type = date
    null = false
  }
  column "due_date" {
    type = date
    null = true
  }
  column "payment_date" {
    type = date
    null = true
  }
  column "is_paid" {
    type    = bool
    default = false
  }
  column "is_recurring" {
    type    = bool
    default = false
  }
  column "recurring_id" {
    type = char(36)
    null = true
  }
  column "installment_number" {
    type = int
    null = true
  }
  column "total_installments" {
    type = int
    null = true
  }
  column "installment_group_id" {
    type = char(36)
    null = true
  }
  column "tags" {
    type = json
    null = true
  }
  column "attachment_url" {
    type = varchar(500)
    null = true
  }
  column "external_id" {
    type = varchar(255)
    null = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_transactions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_transactions_bank_account" {
    columns     = [column.bank_account_id]
    ref_columns = [table.bank_accounts.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_transactions_category" {
    columns     = [column.category_id]
    ref_columns = [table.categories.column.id]
    on_delete   = SET_NULL
  }

  index "idx_transactions_user_id" {
    columns = [column.user_id]
  }
  index "idx_transactions_bank_account_id" {
    columns = [column.bank_account_id]
  }
  index "idx_transactions_category_id" {
    columns = [column.category_id]
  }
  index "idx_transactions_type" {
    columns = [column.type]
  }
  index "idx_transactions_transaction_date" {
    columns = [column.transaction_date]
  }
  index "idx_transactions_due_date" {
    columns = [column.due_date]
  }
  index "idx_transactions_is_paid" {
    columns = [column.is_paid]
  }
  index "idx_transactions_recurring_id" {
    columns = [column.recurring_id]
  }
  index "idx_transactions_installment_group_id" {
    columns = [column.installment_group_id]
  }
  index "idx_transactions_created_at" {
    columns = [column.created_at]
  }
  index "idx_transactions_user_date" {
    columns = [column.user_id, column.transaction_date]
  }
  index "idx_transactions_user_type_date" {
    columns = [column.user_id, column.type, column.transaction_date]
  }
  index "idx_transactions_account_date" {
    columns = [column.bank_account_id, column.transaction_date]
  }
}

 ====================================
// RECURRING TRANSACTIONS TABLE
 ====================================
table "recurring_transactions" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "bank_account_id" {
    type = char(36)
    null = false
  }
  column "category_id" {
    type = char(36)
    null = true
  }
  column "type" {
    type = enum("income", "expense")
    null = false
  }
  column "amount" {
    type = decimal(15, 2)
    null = false
  }
  column "description" {
    type = varchar(255)
    null = false
  }
  column "frequency" {
    type = enum("daily", "weekly", "biweekly", "monthly", "quarterly", "yearly")
    null = false
  }
  column "day_of_month" {
    type = tinyint
    null = true
  }
  column "day_of_week" {
    type = tinyint
    null = true
  }
  column "start_date" {
    type = date
    null = false
  }
  column "end_date" {
    type = date
    null = true
  }
  column "next_occurrence" {
    type = date
    null = false
  }
  column "last_generated_at" {
    type = date
    null = true
  }
  column "is_active" {
    type    = bool
    default = true
  }
  column "auto_confirm" {
    type    = bool
    default = false
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_recurring_transactions_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_recurring_transactions_bank_account" {
    columns     = [column.bank_account_id]
    ref_columns = [table.bank_accounts.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_recurring_transactions_category" {
    columns     = [column.category_id]
    ref_columns = [table.categories.column.id]
    on_delete   = SET_NULL
  }

  index "idx_recurring_transactions_user_id" {
    columns = [column.user_id]
  }
  index "idx_recurring_transactions_next_occurrence" {
    columns = [column.next_occurrence]
  }
  index "idx_recurring_transactions_is_active" {
    columns = [column.is_active]
  }
}

 ====================================
// BUDGETS TABLE
 ====================================
table "budgets" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "category_id" {
    type = char(36)
    null = true
  }
  column "name" {
    type = varchar(100)
    null = false
  }
  column "amount" {
    type = decimal(15, 2)
    null = false
  }
  column "spent_amount" {
    type    = decimal(15, 2)
    default = 0.00
  }
  column "period_type" {
    type = enum("monthly", "quarterly", "yearly", "custom")
    null = false
  }
  column "start_date" {
    type = date
    null = false
  }
  column "end_date" {
    type = date
    null = false
  }
  column "alert_threshold" {
    type    = decimal(5, 2)
    default = 80.00
  }
  column "alert_sent" {
    type    = bool
    default = false
  }
  column "is_active" {
    type    = bool
    default = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_budgets_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_budgets_category" {
    columns     = [column.category_id]
    ref_columns = [table.categories.column.id]
    on_delete   = CASCADE
  }

  index "idx_budgets_user_id" {
    columns = [column.user_id]
  }
  index "idx_budgets_category_id" {
    columns = [column.category_id]
  }
  index "idx_budgets_period" {
    columns = [column.start_date, column.end_date]
  }
  index "idx_budgets_is_active" {
    columns = [column.is_active]
  }
}

 ====================================
// GOALS TABLE (Metas Financeiras)
 ====================================
table "goals" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "name" {
    type = varchar(100)
    null = false
  }
  column "description" {
    type = text
    null = true
  }
  column "target_amount" {
    type = decimal(15, 2)
    null = false
  }
  column "current_amount" {
    type    = decimal(15, 2)
    default = 0.00
  }
  column "target_date" {
    type = date
    null = true
  }
  column "icon" {
    type    = varchar(50)
    default = "target"
  }
  column "color" {
    type    = varchar(7)
    default = "#F59E0B"
  }
  column "priority" {
    type    = tinyint
    default = 1
  }
  column "status" {
    type    = enum("in_progress", "completed", "cancelled")
    default = "in_progress"
  }
  column "completed_at" {
    type = datetime
    null = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_goals_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  index "idx_goals_user_id" {
    columns = [column.user_id]
  }
  index "idx_goals_status" {
    columns = [column.status]
  }
  index "idx_goals_target_date" {
    columns = [column.target_date]
  }
}

 ====================================
// GOAL CONTRIBUTIONS TABLE
 ====================================
table "goal_contributions" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "goal_id" {
    type = char(36)
    null = false
  }
  column "amount" {
    type = decimal(15, 2)
    null = false
  }
  column "note" {
    type = varchar(255)
    null = true
  }
  column "contribution_date" {
    type    = date
    null    = false
    default = sql("(CURRENT_DATE)")
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_goal_contributions_goal" {
    columns     = [column.goal_id]
    ref_columns = [table.goals.column.id]
    on_delete   = CASCADE
  }

  index "idx_goal_contributions_goal_id" {
    columns = [column.goal_id]
  }
  index "idx_goal_contributions_date" {
    columns = [column.contribution_date]
  }
}

 ====================================
// TRANSFERS TABLE (for transfers between accounts)
 ====================================
table "transfers" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "from_account_id" {
    type = char(36)
    null = false
  }
  column "to_account_id" {
    type = char(36)
    null = false
  }
  column "amount" {
    type = decimal(15, 2)
    null = false
  }
  column "description" {
    type = varchar(255)
    null = true
  }
  column "transfer_date" {
    type = date
    null = false
  }
  column "from_transaction_id" {
    type = char(36)
    null = true
  }
  column "to_transaction_id" {
    type = char(36)
    null = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_transfers_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_transfers_from_account" {
    columns     = [column.from_account_id]
    ref_columns = [table.bank_accounts.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_transfers_to_account" {
    columns     = [column.to_account_id]
    ref_columns = [table.bank_accounts.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_transfers_from_transaction" {
    columns     = [column.from_transaction_id]
    ref_columns = [table.transactions.column.id]
    on_delete   = SET_NULL
  }
  foreign_key "fk_transfers_to_transaction" {
    columns     = [column.to_transaction_id]
    ref_columns = [table.transactions.column.id]
    on_delete   = SET_NULL
  }

  index "idx_transfers_user_id" {
    columns = [column.user_id]
  }
  index "idx_transfers_from_account" {
    columns = [column.from_account_id]
  }
  index "idx_transfers_to_account" {
    columns = [column.to_account_id]
  }
  index "idx_transfers_date" {
    columns = [column.transfer_date]
  }
}

 ====================================
// INSTALLMENTS TABLE (for installment tracking)
 ====================================
table "installments" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "bank_account_id" {
    type = char(36)
    null = false
  }
  column "category_id" {
    type = char(36)
    null = true
  }
  column "description" {
    type = varchar(255)
    null = false
  }
  column "total_amount" {
    type = decimal(15, 2)
    null = false
  }
  column "installment_amount" {
    type = decimal(15, 2)
    null = false
  }
  column "total_installments" {
    type = int
    null = false
  }
  column "paid_installments" {
    type    = int
    default = 0
  }
  column "first_due_date" {
    type = date
    null = false
  }
  column "is_active" {
    type    = bool
    default = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_installments_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_installments_bank_account" {
    columns     = [column.bank_account_id]
    ref_columns = [table.bank_accounts.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_installments_category" {
    columns     = [column.category_id]
    ref_columns = [table.categories.column.id]
    on_delete   = SET_NULL
  }

  index "idx_installments_user_id" {
    columns = [column.user_id]
  }
  index "idx_installments_bank_account_id" {
    columns = [column.bank_account_id]
  }
  index "idx_installments_is_active" {
    columns = [column.is_active]
  }
}

 ====================================
// NOTIFICATIONS TABLE
 ====================================
table "notifications" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "type" {
    type = enum("budget_alert", "bill_reminder", "goal_achieved", "low_balance", "transaction_alert", "system")
    null = false
  }
  column "title" {
    type = varchar(200)
    null = false
  }
  column "message" {
    type = text
    null = false
  }
  column "data" {
    type = json
    null = true
  }
  column "is_read" {
    type    = bool
    default = false
  }
  column "read_at" {
    type = datetime
    null = true
  }
  column "sent_via" {
    type    = json
    default = sql("('[\"in_app\"]')")
  }
  column "scheduled_for" {
    type = datetime
    null = true
  }
  column "sent_at" {
    type = datetime
    null = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_notifications_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  index "idx_notifications_user_id" {
    columns = [column.user_id]
  }
  index "idx_notifications_is_read" {
    columns = [column.is_read]
  }
  index "idx_notifications_type" {
    columns = [column.type]
  }
  index "idx_notifications_created_at" {
    columns = [column.created_at]
  }
  index "idx_notifications_scheduled" {
    columns = [column.scheduled_for]
  }
}

 ====================================
// USER SETTINGS TABLE
 ====================================
table "user_settings" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "notification_email" {
    type    = bool
    default = true
  }
  column "notification_push" {
    type    = bool
    default = true
  }
  column "notification_sms" {
    type    = bool
    default = false
  }
  column "budget_alerts" {
    type    = bool
    default = true
  }
  column "bill_reminders" {
    type    = bool
    default = true
  }
  column "bill_reminder_days" {
    type    = int
    default = 3
  }
  column "weekly_summary" {
    type    = bool
    default = true
  }
  column "monthly_report" {
    type    = bool
    default = true
  }
  column "low_balance_alert" {
    type    = bool
    default = true
  }
  column "low_balance_threshold" {
    type    = decimal(15, 2)
    default = 100.00
  }
  column "allow_manual_transactions" {
    type    = bool
    default = true
  }
  column "theme" {
    type    = enum("light", "dark", "system")
    default = "system"
  }
  column "dashboard_layout" {
    type = json
    null = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_user_settings_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  index "unique_user_settings" {
    unique  = true
    columns = [column.user_id]
  }

  index "idx_user_settings_user_id" {
    columns = [column.user_id]
  }
}

 ====================================
// BANK INTEGRATIONS TABLE (for Open Banking)
 ====================================
table "bank_integrations" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "bank_account_id" {
    type = char(36)
    null = true
  }
  column "provider" {
    type = varchar(50)
    null = false
  }
  column "provider_account_id" {
    type = varchar(255)
    null = true
  }
  column "access_token_encrypted" {
    type = text
    null = true
  }
  column "refresh_token_encrypted" {
    type = text
    null = true
  }
  column "token_expires_at" {
    type = datetime
    null = true
  }
  column "consent_expires_at" {
    type = datetime
    null = true
  }
  column "status" {
    type    = enum("pending", "active", "expired", "revoked", "error")
    default = "pending"
  }
  column "last_sync_at" {
    type = datetime
    null = true
  }
  column "last_error" {
    type = text
    null = true
  }
  column "metadata" {
    type = json
    null = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_bank_integrations_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
  foreign_key "fk_bank_integrations_bank_account" {
    columns     = [column.bank_account_id]
    ref_columns = [table.bank_accounts.column.id]
    on_delete   = SET_NULL
  }

  index "idx_bank_integrations_user_id" {
    columns = [column.user_id]
  }
  index "idx_bank_integrations_status" {
    columns = [column.status]
  }
  index "idx_bank_integrations_provider" {
    columns = [column.provider]
  }
}

 ====================================
// AUDIT LOG TABLE
 ====================================
table "audit_logs" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = true
  }
  column "action" {
    type = varchar(50)
    null = false
  }
  column "entity_type" {
    type = varchar(50)
    null = false
  }
  column "entity_id" {
    type = char(36)
    null = true
  }
  column "old_values" {
    type = json
    null = true
  }
  column "new_values" {
    type = json
    null = true
  }
  column "ip_address" {
    type = varchar(45)
    null = true
  }
  column "user_agent" {
    type = varchar(500)
    null = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_audit_logs_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = SET_NULL
  }

  index "idx_audit_logs_user_id" {
    columns = [column.user_id]
  }
  index "idx_audit_logs_action" {
    columns = [column.action]
  }
  index "idx_audit_logs_entity" {
    columns = [column.entity_type, column.entity_id]
  }
  index "idx_audit_logs_created_at" {
    columns = [column.created_at]
  }
}

 ====================================
// REPORTS TABLE (for scheduled/saved reports)
 ====================================
table "reports" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "name" {
    type = varchar(100)
    null = false
  }
  column "type" {
    type = enum("cash_flow", "expense_by_category", "income_vs_expense", "budget_analysis", "custom")
    null = false
  }
  column "parameters" {
    type = json
    null = false
  }
  column "schedule" {
    type = enum("daily", "weekly", "monthly")
    null = true
  }
  column "last_generated_at" {
    type = datetime
    null = true
  }
  column "file_url" {
    type = varchar(500)
    null = true
  }
  column "is_active" {
    type    = bool
    default = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_reports_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  index "idx_reports_user_id" {
    columns = [column.user_id]
  }
  index "idx_reports_type" {
    columns = [column.type]
  }
  index "idx_reports_schedule" {
    columns = [column.schedule]
  }
}

 ====================================
// PASSWORD RESET TOKENS TABLE
 ====================================
table "password_reset_tokens" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "token_hash" {
    type = varchar(255)
    null = false
  }
  column "expires_at" {
    type = datetime
    null = false
  }
  column "used_at" {
    type = datetime
    null = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_password_reset_tokens_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  index "idx_password_reset_user_id" {
    columns = [column.user_id]
  }
  index "idx_password_reset_token" {
    columns = [column.token_hash]
  }
  index "idx_password_reset_expires" {
    columns = [column.expires_at]
  }
}

 ====================================
// EMAIL VERIFICATION TOKENS TABLE
 ====================================
table "email_verification_tokens" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "token_hash" {
    type = varchar(255)
    null = false
  }
  column "expires_at" {
    type = datetime
    null = false
  }
  column "used_at" {
    type = datetime
    null = true
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_email_verification_tokens_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  index "idx_email_verification_user_id" {
    columns = [column.user_id]
  }
  index "idx_email_verification_token" {
    columns = [column.token_hash]
  }
  index "idx_email_verification_expires" {
    columns = [column.expires_at]
  }
}

 ====================================
// USER OAUTH PROVIDERS TABLE (for social login - Google)
 ====================================
table "user_oauth_providers" {
  schema = schema.cashing

  column "id" {
    type    = char(36)
    null    = false
    default = sql("(UUID())")
  }
  column "user_id" {
    type = char(36)
    null = false
  }
  column "provider" {
    type = varchar(50)
    null = false
    comment = "Provider name: google, facebook, github, etc"
  }
  column "provider_user_id" {
    type = varchar(255)
    null = false
    comment = "Unique user ID from provider (Firebase UID)"
  }
  column "provider_email" {
    type = varchar(255)
    null = true
    comment = "Email from provider (may differ from user.email)"
  }
  column "provider_name" {
    type = varchar(255)
    null = true
    comment = "Display name from provider"
  }
  column "provider_avatar_url" {
    type = varchar(500)
    null = true
    comment = "Profile picture URL from provider"
  }
  column "is_primary" {
    type    = bool
    default = false
    comment = "True if this provider was used to create the account"
  }
  column "metadata" {
    type = json
    null = true
    comment = "Additional provider-specific data"
  }
  column "last_login_at" {
    type = datetime
    null = true
    comment = "Last time user logged in via this provider"
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_user_oauth_providers_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }

  // Garante que um usuário só pode vincular cada provider uma vez
  index "unique_user_provider" {
    unique  = true
    columns = [column.user_id, column.provider]
  }

  // Garante que um UID do provider não pode ser usado em múltiplas contas
  index "unique_provider_user" {
    unique  = true
    columns = [column.provider, column.provider_user_id]
  }

  index "idx_user_oauth_providers_user_id" {
    columns = [column.user_id]
  }

  index "idx_user_oauth_providers_provider" {
    columns = [column.provider]
  }

  index "idx_user_oauth_providers_last_login" {
    columns = [column.last_login_at]
  }
}
