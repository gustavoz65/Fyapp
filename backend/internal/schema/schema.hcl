// Schema definition for Cashing-go database
// This file defines the database structure in Atlas HCL format
// When you change this file, run: task migrations:new name=your_change

// Define the database schema (uses the database name from atlas.hcl)
schema "cashing" {
}

// ============================================================================
// USERS TABLE
// ============================================================================
table "users" {
  schema = schema.cashing
  
  column "id" {
    type = varchar(36)
  }
  column "name" {
    type = varchar(255)
    null = false
  }
  column "email" {
    type = varchar(255)
    null = false
  }
  column "phone" {
    type = varchar(20)
    null = true
    comment = "User phone number (optional)"
  }
  column "password_hash" {
    type = varchar(255)
    null = false
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
    unique  = true
    columns = [column.email]
  }
}

// ============================================================================
// CATEGORIES TABLE
// ============================================================================
table "categories" {
  schema = schema.cashing

  column "id" {
    type = varchar(36)
  }
  column "user_id" {
    type = varchar(36)
    null = true
  }
  column "name" {
    type = varchar(100)
    null = false
  }
  column "description" {
    type = text
    null = true
  }
  column "type" {
    type = enum("income", "expense")
    null = false
  }
  column "color" {
    type    = varchar(7)
    default = "#6B7280"
  }
  column "icon" {
    type    = varchar(50)
    default = "folder"
  }
  column "is_default" {
    type    = bool
    default = false
  }
  column "created_at" {
    type    = datetime
    default = sql("CURRENT_TIMESTAMP")
  }
  
  primary_key {
    columns = [column.id]
  }
  
  foreign_key "fk_categories_user" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
  
  index "idx_categories_user_id" {
    columns = [column.user_id]
  }
  index "idx_categories_type" {
    columns = [column.type]
  }
}

// ============================================================================
// TRANSACTIONS TABLE
// ============================================================================
table "transactions" {
  schema = schema.cashing

  column "id" {
    type = varchar(36)
  }
  column "user_id" {
    type = varchar(36)
    null = false
  }
  column "category_id" {
    type = varchar(36)
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
    null = true
  }
  column "date" {
    type = date
    null = false
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
  foreign_key "fk_transactions_category" {
    columns     = [column.category_id]
    ref_columns = [table.categories.column.id]
    on_delete   = SET_NULL
  }
  
  index "idx_transactions_user_id" {
    columns = [column.user_id]
  }
  index "idx_transactions_date" {
    columns = [column.date]
  }
}
