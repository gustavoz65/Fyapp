package config

import (
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary       Primary              `koanf:"primary" validate:"required"`
	Server        ServerConfig         `koanf:"server" validate:"required"`
	Database      DatabaseConfig       `koanf:"database" validate:"required"`
	Auth          AuthConfig           `koanf:"auth" validate:"required"`
	Redis         RedisConfig          `koanf:"redis" validate:"required"`
	Observability *ObservabilityConfig `koanf:"observability"`
	Integration   IntegrationConfig    `koanf:"integration"`
}

type Primary struct {
	Env string `koanf:"env" validate:"required"`
}

type ServerConfig struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout"`
	WriteTimeout       int      `koanf:"write_timeout"`
	IdleTimeout        int      `koanf:"idle_timeout"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins"`
}

type DatabaseConfig struct {
	Host            string `koanf:"host" validate:"required"`
	Port            int    `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password"`
	Name            string `koanf:"db_name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode"`
	MaxOpenConns    int    `koanf:"max_open_conns"`
	MaxIdleConns    int    `koanf:"max_idle_conns"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time"`
}

type AuthConfig struct {
	SecretKey            string `koanf:"secret_key" validate:"required"`
	AccessTokenDuration  int    `koanf:"access_token_duration"`  // in minutes
	RefreshTokenDuration int    `koanf:"refresh_token_duration"` // in hours
	Issuer               string `koanf:"issuer"`
}

func (a *AuthConfig) GetSecretKey() string {
	return a.SecretKey
}

type IntegrationConfig struct {
	ResendAPIKey string `koanf:"resend_api_key"`
}

type RedisConfig struct {
	Address  string `koanf:"address" validate:"required"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

// transformEnvKey converte nomes de variáveis de ambiente para chaves koanf
// Mapeamentos explícitos preservam underscores nos nomes dos campos
func transformEnvKey(s string) string {
	s = strings.TrimPrefix(s, "BOILERPLATE_")
	s = strings.ToLower(s)

	// Mapeamentos explícitos para campos com underscores
	replacements := map[string]string{
		"primary_env":                                         "primary.env",
		"server_port":                                         "server.port",
		"server_read_timeout":                                 "server.read_timeout",
		"server_write_timeout":                                "server.write_timeout",
		"server_idle_timeout":                                 "server.idle_timeout",
		"server_cors_allowed_origins":                         "server.cors_allowed_origins",
		"database_host":                                       "database.host",
		"database_port":                                       "database.port",
		"database_user":                                       "database.user",
		"database_password":                                   "database.password",
		"database_db_name":                                    "database.db_name",
		"database_ssl_mode":                                   "database.ssl_mode",
		"database_max_open_conns":                             "database.max_open_conns",
		"database_max_idle_conns":                             "database.max_idle_conns",
		"database_conn_max_lifetime":                          "database.conn_max_lifetime",
		"database_conn_max_idle_time":                         "database.conn_max_idle_time",
		"auth_secret_key":                                     "auth.secret_key",
		"auth_access_token_duration":                          "auth.access_token_duration",
		"auth_refresh_token_duration":                         "auth.refresh_token_duration",
		"auth_issuer":                                         "auth.issuer",
		"redis_address":                                       "redis.address",
		"redis_password":                                      "redis.password",
		"redis_db":                                            "redis.db",
		"integration_resend_api_key":                          "integration.resend_api_key",
		"observability_service_name":                          "observability.service_name",
		"observability_environment":                           "observability.environment",
		"observability_logging_level":                         "observability.logging.level",
		"observability_logging_format":                        "observability.logging.format",
		"observability_logging_slow_query_threshold":          "observability.logging.slow_query_threshold",
		"observability_new_relic_license_key":                 "observability.new_relic.license_key",
		"observability_new_relic_app_log_forwarding_enabled":  "observability.new_relic.app_log_forwarding_enabled",
		"observability_new_relic_distributed_tracing_enabled": "observability.new_relic.distributed_tracing_enabled",
		"observability_new_relic_debug_logging":               "observability.new_relic.debug_logging",
		"observability_health_checks_enabled":                 "observability.health_checks.enabled",
		"observability_health_checks_interval":                "observability.health_checks.interval",
		"observability_health_checks_timeout":                 "observability.health_checks.timeout",
		"observability_health_checks_checks":                  "observability.health_checks.checks",
	}

	if mapped, ok := replacements[s]; ok {
		return mapped
	}

	// Padrão: substitui todos os underscores por pontos
	return strings.ReplaceAll(s, "_", ".")
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")

	// Transforma variáveis de ambiente com mapeamentos explícitos para preservar underscores nos nomes dos campos
	err := k.Load(env.Provider("BOILERPLATE_", ".", transformEnvKey), nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig := &Config{}

	err = k.Unmarshal("", mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal config into struct")
	}

	setDefaults(mainConfig)

	validate := validator.New()

	err = validate.Struct(mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("config validation failed")
	}

	if mainConfig.Observability == nil {
		mainConfig.Observability = DefaultObservabilityConfig()
	}

	mainConfig.Observability.ServiceName = "cashing-api"
	mainConfig.Observability.Environment = mainConfig.Primary.Env

	if err := mainConfig.Observability.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("invalid observability config")
	}

	return mainConfig, nil
}

func setDefaults(cfg *Config) {
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30
	}
	if cfg.Server.IdleTimeout == 0 {
		cfg.Server.IdleTimeout = 60
	}

	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 25
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 5
	}
	if cfg.Database.ConnMaxLifetime == 0 {
		cfg.Database.ConnMaxLifetime = 300
	}
	if cfg.Database.ConnMaxIdleTime == 0 {
		cfg.Database.ConnMaxIdleTime = 60
	}

	if cfg.Auth.AccessTokenDuration == 0 {
		cfg.Auth.AccessTokenDuration = 15
	}
	if cfg.Auth.RefreshTokenDuration == 0 {
		cfg.Auth.RefreshTokenDuration = 168
	}
	if cfg.Auth.Issuer == "" {
		cfg.Auth.Issuer = "cashing-api"
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
}
