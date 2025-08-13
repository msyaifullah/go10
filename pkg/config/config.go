// pkg/config/config.go
package config

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	App      AppConfig      `toml:"app"`
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
	Redis    RedisConfig    `toml:"redis"`
	Logger   LoggerConfig   `toml:"logger"`
	Email    EmailConfig    `toml:"email"`
	Payment  PaymentConfig  `toml:"payment"`
	Cron     CronConfig     `toml:"cron"`
}

type AppConfig struct {
	Name        string `toml:"name"`
	Environment string `toml:"environment"`
	Debug       bool   `toml:"debug"`
}

type ServerConfig struct {
	Port           string        `toml:"port"`
	Host           string        `toml:"host"`
	ReadTimeout    time.Duration `toml:"read_timeout"`
	WriteTimeout   time.Duration `toml:"write_timeout"`
	TrustedProxies []string      `toml:"trusted_proxies"`
}

type DatabaseConfig struct {
	Host            string        `toml:"host"`
	Port            int           `toml:"port"`
	User            string        `toml:"user"`
	Password        string        `toml:"password"`
	DBName          string        `toml:"dbname"`
	SSLMode         string        `toml:"sslmode"`
	MaxOpenConns    int           `toml:"max_open_conns"`
	MaxIdleConns    int           `toml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `toml:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host         string        `toml:"host"`
	Port         int           `toml:"port"`
	Password     string        `toml:"password"`
	DB           int           `toml:"db"`
	PoolSize     int           `toml:"pool_size"`
	MinIdleConns int           `toml:"min_idle_conns"`
	MaxRetries   int           `toml:"max_retries"`
	DialTimeout  time.Duration `toml:"dial_timeout"`
	ReadTimeout  time.Duration `toml:"read_timeout"`
	WriteTimeout time.Duration `toml:"write_timeout"`
	PoolTimeout  time.Duration `toml:"pool_timeout"`
	IdleTimeout  time.Duration `toml:"idle_timeout"`
}

type LoggerConfig struct {
	Level         string `toml:"level"`
	Format        string `toml:"format"`
	Output        string `toml:"output"`
	MaskSensitive bool   `toml:"mask_sensitive"`
}

type EmailConfig struct {
	Provider     string `toml:"provider"`
	APIKey       string `toml:"api_key"`
	SMTPHost     string `toml:"smtp_host"`
	SMTPPort     int    `toml:"smtp_port"`
	SMTPUsername string `toml:"smtp_username"`
	SMTPPassword string `toml:"smtp_password"`
	FromAddress  string `toml:"from_address"`
}

type PaymentConfig struct {
	Provider      string `toml:"provider"`
	APIKey        string `toml:"api_key"`
	SecretKey     string `toml:"secret_key"`
	WebhookSecret string `toml:"webhook_secret"`
}

type CronConfig struct {
	InvestmentAgreementSchedule string `toml:"investment_agreement_schedule"`
}

func Load(configPath, environment string) (*Config, error) {
	var config Config

	configFile := fmt.Sprintf("%s/%s.toml", configPath, environment)

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, fmt.Errorf("failed to decode config file %s: %w", configFile, err)
	}

	return &config, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
