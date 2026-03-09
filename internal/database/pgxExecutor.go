package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGExecutor struct {
	config *pgxpool.Config
	pool   *pgxpool.Pool
}

func NewPGExecutor(config *pgxpool.Config) *PGExecutor {
	return &PGExecutor{
		config: config,
	}
}

func (pg *PGExecutor) Execute(ctx context.Context, statement string, args ...any) (int64, error) {

	tag, err := pg.pool.Exec(ctx, statement, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (pg *PGExecutor) QueryRow(ctx context.Context, dest []any, statement string, args ...any) error {
	err := pg.pool.QueryRow(ctx, statement, args...).Scan(dest...)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PGExecutor) Query(ctx context.Context, statement string, args ...any) ([]map[string]any, error) {
	allRows := []map[string]any{}
	rows, err := pg.pool.Query(ctx, statement, args...)
	if err != nil {
		return allRows, err
	}
	defer rows.Close()
	fds := rows.FieldDescriptions()
	for rows.Next() {
		values, _ := rows.Values()
		row := map[string]any{}
		for i, fd := range fds {
			row[string(fd.Name)] = values[i]
		}
		allRows = append(allRows, row)
	}
	return allRows, nil
}

func (pg *PGExecutor) Connect() error {
	poolCon, err := pgxpool.NewWithConfig(context.Background(), pg.config)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	pg.pool = poolCon
	return nil
}

func (pg *PGExecutor) Disconnect() {
	if pg.pool != nil {
		pg.pool.Close()
	}
}
