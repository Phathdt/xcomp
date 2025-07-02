package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseConnection struct {
	Config ConfigService `inject:"ConfigService"`
	db     *pgxpool.Pool
}

func (dc *DatabaseConnection) GetServiceName() string {
	return "DatabaseConnection"
}

func (dc *DatabaseConnection) Initialize() error {
	host := dc.Config.GetString("database.host", "localhost")
	port := dc.Config.GetInt("database.port", 5432)
	user := dc.Config.GetString("database.username", "postgres")
	password := dc.Config.GetString("database.password", "password")
	dbname := dc.Config.GetString("database.database", "productdb")
	sslmode := dc.Config.GetString("database.sslmode", "disable")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse database config: %w", err)
	}

	maxConnections := dc.Config.GetInt("database.max_connections", 25)
	maxIdleConnections := dc.Config.GetInt("database.max_idle_connections", 10)
	maxLifetimeMinutes := dc.Config.GetInt("database.max_lifetime_minutes", 30)

	config.MaxConns = int32(maxConnections)
	config.MinConns = int32(maxIdleConnections)
	config.MaxConnLifetime = time.Duration(maxLifetimeMinutes) * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	dc.db = pool
	return nil
}

func (dc *DatabaseConnection) GetDB() *pgxpool.Pool {
	return dc.db
}

func (dc *DatabaseConnection) Close() error {
	if dc.db != nil {
		dc.db.Close()
	}
	return nil
}

func (dc *DatabaseConnection) HealthCheck(ctx context.Context) error {
	if dc.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	return dc.db.Ping(ctx)
}
