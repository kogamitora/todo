// internal/db/db.go
package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB(dsn string, logger *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 设置连接池参数
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// 检查连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established")
	return db, nil
}
