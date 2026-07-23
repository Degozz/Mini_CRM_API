package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// 1. Адаптер (DB), отвечающий за пул соединений к БД Postgres.
type DB struct {
	Pool *pgxpool.Pool
}

// 2. Connect-функция для создания пула соединений с Postgres.
func Connect(ctx context.Context, connStr string) (*DB, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать пул соединений: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("База данных не отвечает на ping: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// 3. Close-функция для закрытия пула соединений при выключении сервера.
func (db *DB) Close() {
	if db != nil && db.Pool != nil {
		db.Pool.Close()
	}
}

// 4. Прокси-методы, чтобы репозитории могли делать запросы напрямую через нашу структуру DB.

// I. Exec-метод для SQL-запросов в Postgres, изменяющих состояние данных, но не возвращающих обратно строки из таблиц (INSERT, UPDATE, DELETE).
func (db *DB) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return db.Pool.Exec(ctx, query, args...)
}

// II. Query-метод для SQL-запросов в Postgres на чтение и выборку множества строк из таблиц (SELECT).
func (db *DB) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return db.Pool.Query(ctx, query, args...)
}

// III. QueryRow-метод для SQL-запросов в Postgres на чтение и выборку одной конкретной строки из таблицы (SELECT…WHERE или SELECT…LIMIT 1)
func (db *DB) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return db.Pool.QueryRow(ctx, query, args...)
}
