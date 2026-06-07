package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	LLM      LLMConfig
}

type AppConfig struct {
	Port string
}

type DatabaseConfig struct {
	Dsn          string
	MaxIdleConns int
	MaxOpenConns int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
}

type LLMConfig struct {
	BaseURL   string
	APIKey    string
	Model     string
	MaxTokens int
}

func Load() *Config {
	// Try config.yml as fallback
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")
	_ = viper.ReadInConfig() // ignore error, env vars take priority

	cfg := &Config{
		App: AppConfig{
			Port: envOr("APP_PORT", viper.GetString("app.port")),
		},
		Database: DatabaseConfig{
			Dsn:          envOr("DB_DSN", viper.GetString("database.dsn")),
			MaxIdleConns: viper.GetInt("database.maxidleconns"),
			MaxOpenConns: viper.GetInt("database.maxopenconns"),
		},
		Redis: RedisConfig{
			Addr:     envOr("REDIS_ADDR", viper.GetString("redis.addr")),
			Password: envOr("REDIS_PASSWORD", viper.GetString("redis.password")),
			DB:       viper.GetInt("redis.db"),
		},
		JWT: JWTConfig{
			Secret: envOr("JWT_SECRET", ""),
		},
		LLM: LLMConfig{
			BaseURL:   envOr("LLM_BASE_URL", viper.GetString("llm.base_url")),
			APIKey:    envOr("LLM_API_KEY", viper.GetString("llm.api_key")),
			Model:     envOr("LLM_MODEL", viper.GetString("llm.model")),
			MaxTokens: viper.GetInt("llm.max_tokens"),
		},
	}

	if cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required (set via env or config)")
	}
	if cfg.Database.Dsn == "" {
		log.Fatal("DB_DSN is required (set via env or config)")
	}
	if cfg.Redis.Addr == "" {
		cfg.Redis.Addr = "localhost:6379"
	}
	if cfg.App.Port == "" {
		cfg.App.Port = ":3000"
	}
	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 10
	}
	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 100
	}
	if cfg.LLM.BaseURL == "" {
		cfg.LLM.BaseURL = "https://api.anthropic.com"
	}
	if cfg.LLM.Model == "" {
		cfg.LLM.Model = "claude-sonnet-4-20250514"
	}
	if cfg.LLM.MaxTokens == 0 {
		cfg.LLM.MaxTokens = 2048
	}

	return cfg
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
