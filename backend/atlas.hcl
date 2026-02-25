// Atlas configuration for Cashing-go project

// Define the environment variables
variable "db_user" {
  type    = string
  default = "root"
}

variable "db_pass" {
  type    = string
  default = "senha123"
}

variable "db_host" {
  type    = string
  default = "localhost"
}

variable "db_port" {
  type    = string
  default = "3306"
}

variable "db_name" {
  type    = string
  default = "cashing"
}

// Main database connection
env "local" {
  // Schema source: HCL files that define the desired database structure
  src = "file://internal/schema/schema.hcl"

  // Target database URL (where migrations will be applied)
  url = "mysql://${var.db_user}:${var.db_pass}@${var.db_host}:${var.db_port}/${var.db_name}?parseTime=true"

  // Dev database for schema comparison (Atlas uses this to calculate diffs)
  // This database is temporary and will be created/destroyed by Atlas
  dev = "mysql://${var.db_user}:${var.db_pass}@${var.db_host}:${var.db_port}/cashing_dev?parseTime=true"

  // Migration files directory
  migration {
    dir = "file://internal/database/migrations"
  }

  // Format for generated migration files
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

// Production environment (when needed)
env "prod" {
  src = "file://internal/models"

  // Will use production credentials from environment
  url = "mysql://${var.db_user}:${var.db_pass}@${var.db_host}:${var.db_port}/${var.db_name}?parseTime=true"

  dev = "mysql://${var.db_user}:${var.db_pass}@${var.db_host}:${var.db_port}/cashing_dev?parseTime=true"

  migration {
    dir = "file://internal/database/migrations"
  }
}
