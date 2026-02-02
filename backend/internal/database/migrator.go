package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"strconv"
	"strings"

	"github.com/gustavoz65/Cashing-go/backend/internal/config"
	"github.com/rs/zerolog"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed migrations/*.sql
var migrations embed.FS

// MigrationVersion represents a migration record
type MigrationVersion struct {
	Version int
	Name    string
}

func Migrate(ctx context.Context, logger *zerolog.Logger, cfg *config.Config) error {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

	// MySQL DSN format
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&multiStatements=true",
		cfg.Database.User,
		cfg.Database.Password,
		hostPort,
		cfg.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	// Create schema_version table if it doesn't exist
	if err := createSchemaVersionTable(ctx, db); err != nil {
		return fmt.Errorf("failed to create schema version table: %w", err)
	}

	// Get current version
	currentVersion, err := getCurrentVersion(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	// Load and execute migrations
	subtree, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to retrieve migrations subtree: %w", err)
	}

	entries, err := fs.ReadDir(subtree, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}

	if len(migrationFiles) == 0 {
		logger.Info().Msg("no migration files found")
		return nil
	}

	newVersion := currentVersion
	for i, fileName := range migrationFiles {
		version := i + 1
		if version <= currentVersion {
			continue
		}

		logger.Info().Msgf("applying migration %d: %s", version, fileName)

		content, err := fs.ReadFile(subtree, fileName)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		// Get only the "up" part of the migration (before the separator)
		upMigration := extractUpMigration(string(content))

		// Execute the migration
		if err := executeMigration(ctx, db, upMigration); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}

		// Update version
		if err := setVersion(ctx, db, version); err != nil {
			return fmt.Errorf("failed to update schema version: %w", err)
		}

		newVersion = version
	}

	if newVersion == currentVersion {
		logger.Info().Msgf("database schema up to date, version %d", currentVersion)
	} else {
		logger.Info().Msgf("migrated database schema from version %d to %d", currentVersion, newVersion)
	}

	return nil
}

func createSchemaVersionTable(ctx context.Context, db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_version (
			version INT NOT NULL PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`
	_, err := db.ExecContext(ctx, query)
	return err
}

func getCurrentVersion(ctx context.Context, db *sql.DB) (int, error) {
	var version int
	query := `SELECT COALESCE(MAX(version), 0) FROM schema_version`
	err := db.QueryRowContext(ctx, query).Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func setVersion(ctx context.Context, db *sql.DB, version int) error {
	query := `INSERT INTO schema_version (version) VALUES (?)`
	_, err := db.ExecContext(ctx, query, version)
	return err
}

func extractUpMigration(content string) string {
	// Look for the Tern separator pattern
	separator := "---- create above / drop below ----"
	parts := strings.Split(content, separator)
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return strings.TrimSpace(content)
}

func executeMigration(ctx context.Context, db *sql.DB, migration string) error {
	// Split by semicolon and execute each statement
	statements := splitStatements(migration)

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		_, err := db.ExecContext(ctx, stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %s\nError: %w", truncateString(stmt, 100), err)
		}
	}

	return nil
}

func splitStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inString := false
	stringChar := rune(0)

	for i, char := range sql {
		if !inString && (char == '\'' || char == '"') {
			inString = true
			stringChar = char
		} else if inString && char == stringChar {
			// Check for escaped quote
			if i+1 < len(sql) && rune(sql[i+1]) == stringChar {
				current.WriteRune(char)
				continue
			}
			inString = false
		}

		if char == ';' && !inString {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
		} else {
			current.WriteRune(char)
		}
	}

	// Add the last statement if it doesn't end with semicolon
	stmt := strings.TrimSpace(current.String())
	if stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Rollback rolls back to a specific version
func Rollback(ctx context.Context, logger *zerolog.Logger, cfg *config.Config, targetVersion int) error {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&multiStatements=true",
		cfg.Database.User,
		cfg.Database.Password,
		hostPort,
		cfg.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	currentVersion, err := getCurrentVersion(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	if targetVersion >= currentVersion {
		logger.Info().Msgf("target version %d is not less than current version %d", targetVersion, currentVersion)
		return nil
	}

	subtree, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to retrieve migrations subtree: %w", err)
	}

	entries, err := fs.ReadDir(subtree, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}

	// Rollback migrations in reverse order
	for version := currentVersion; version > targetVersion; version-- {
		if version > len(migrationFiles) {
			continue
		}

		fileName := migrationFiles[version-1]
		logger.Info().Msgf("rolling back migration %d: %s", version, fileName)

		content, err := fs.ReadFile(subtree, fileName)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		downMigration := extractDownMigration(string(content))
		if downMigration == "" {
			logger.Warn().Msgf("no down migration found for version %d", version)
			continue
		}

		if err := executeMigration(ctx, db, downMigration); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", fileName, err)
		}

		if err := removeVersion(ctx, db, version); err != nil {
			return fmt.Errorf("failed to remove schema version: %w", err)
		}
	}

	logger.Info().Msgf("rolled back database schema from version %d to %d", currentVersion, targetVersion)
	return nil
}

func extractDownMigration(content string) string {
	separator := "---- create above / drop below ----"
	parts := strings.Split(content, separator)
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func removeVersion(ctx context.Context, db *sql.DB, version int) error {
	query := `DELETE FROM schema_version WHERE version = ?`
	_, err := db.ExecContext(ctx, query, version)
	return err
}
