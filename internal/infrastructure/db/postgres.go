package db

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импортируем драйвер postgres для side effects
	"time"
)

type PostgresDB struct {
	*sqlx.DB
}

func NewPostgresDB(dsn string, maxOpenConns, maxIdleConns, maxLifetime int) (*PostgresDB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)

	return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) Close() error {
	return p.DB.Close()
}

// HealthCheck проверяет работоспособность БД
func (p *PostgresDB) HealthCheck() error {
	return p.DB.Ping()
}

// BeginTx - начало транзакции с контекстом
func (p *PostgresDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return p.DB.BeginTxx(ctx, opts)
}

// GetSingle - получить одну запись с контекстом
func (p *PostgresDB) GetSingle(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.DB.GetContext(ctx, dest, query, args...)
}

// GetList - получить список записей с контекстом
func (p *PostgresDB) GetList(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return p.DB.SelectContext(ctx, dest, query, args...)
}

// ExecWithContext - выполнить запрос с контекстом
func (p *PostgresDB) ExecWithContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.DB.ExecContext(ctx, query, args...)
}
