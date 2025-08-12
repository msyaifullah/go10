// pkg/database/database.go
package database

import (
	"database/sql"
	"fmt"

	"loan-service/pkg/config"
	"loan-service/pkg/logger"

	_ "github.com/lib/pq"
)

func NewConnection(cfg config.DatabaseConfig, logger *logger.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	logger.Info("Connecting to database", map[string]interface{}{
		"host":    cfg.Host,
		"port":    cfg.Port,
		"dbname":  cfg.DBName,
		"user":    cfg.User,
		"sslmode": cfg.SSLMode,
	})

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to database", map[string]interface{}{
		"max_open_conns":    cfg.MaxOpenConns,
		"max_idle_conns":    cfg.MaxIdleConns,
		"conn_max_lifetime": cfg.ConnMaxLifetime.String(),
	})

	return db, nil
}
