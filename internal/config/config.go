package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

// Config представляет полную конфигурацию приложения
type Config struct {
	Logs     LogsConfig     `toml:"logs"`
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
}

// LogsConfig содержит настройки логирования
type LogsConfig struct {
	Level string `toml:"level"`
	File  string `toml:"file"`
}

// ServerConfig содержит настройки HTTP сервера
type ServerConfig struct {
	HTTPPort        int `toml:"http_port"`
	ReadTimeout     int `toml:"read_timeout"`
	WriteTimeout    int `toml:"write_timeout"`
	IdleTimeout     int `toml:"idle_timeout"`
	ShutdownTimeout int `toml:"shutdown_timeout"`
}

// DatabaseConfig содержит настройки подключения к PostgreSQL
type DatabaseConfig struct {
	Host            string `toml:"host"`
	Port            int    `toml:"port"`
	User            string `toml:"user"`
	Password        string `toml:"password"`
	DBName          string `toml:"dbname"`
	SSLMode         string `toml:"sslmode"`
	MaxOpenConns    int    `toml:"max_open_conns"`
	MaxIdleConns    int    `toml:"max_idle_conns"`
	ConnMaxLifetime int    `toml:"conn_max_lifetime"`
}

// DSN формирует строку подключения к PostgreSQL
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// Load загружает конфигурацию из TOML файла с поддержкой переменных окружения
func Load(path string) (*Config, error) {
	var cfg Config

	// Читаем TOML файл
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode TOML config: %w", err)
	}

	// Переопределяем значения из переменных окружения (если они установлены)
	overrideFromEnv(&cfg)

	// Валидация конфигурации
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// overrideFromEnv переопределяет значения из переменных окружения
func overrideFromEnv(cfg *Config) {
	// Database
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Database.Port = port
		}
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.DBName = v
	}
	if v := os.Getenv("DB_SSLMODE"); v != "" {
		cfg.Database.SSLMode = v
	}

	// Server
	if v := os.Getenv("HTTP_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.HTTPPort = port
		}
	}

	// Logs
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Logs.Level = v
	}
	if v := os.Getenv("LOG_FILE"); v != "" {
		cfg.Logs.File = v
	}
}

// validate проверяет корректность конфигурации
func validate(cfg *Config) error {
	// Database validation
	if cfg.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// Server validation
	if cfg.Server.HTTPPort <= 0 || cfg.Server.HTTPPort > 65535 {
		return fmt.Errorf("HTTP port must be between 1 and 65535")
	}

	// Logs validation
	if cfg.Logs.Level == "" {
		cfg.Logs.Level = "info" // default
	}
	if cfg.Logs.File == "" {
		cfg.Logs.File = "./logs/app.log" // default
	}

	// Set defaults for timeouts if not specified
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 15
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 15
	}
	if cfg.Server.IdleTimeout == 0 {
		cfg.Server.IdleTimeout = 60
	}
	if cfg.Server.ShutdownTimeout == 0 {
		cfg.Server.ShutdownTimeout = 10
	}

	// Set defaults for database connection pool
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 25
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 5
	}
	if cfg.Database.ConnMaxLifetime == 0 {
		cfg.Database.ConnMaxLifetime = 300 // 5 minutes
	}

	return nil
}
