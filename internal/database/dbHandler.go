package database

import (
	"context"
	"fmt"
	"time"
)

// Interface for execute to databse, expected method that must have
//
//	Connect to connect to database server
//	Disconnect to disconnect from database server
//	Query for query item to database and expected to return map of column name to value and error if existed
type DBExecutor interface {
	Connect() error
	Disconnect()
	Execute(context.Context, string, ...any) (int64, error)
	Query(context.Context, string, ...any) ([]map[string]any, error)
	QueryRow(context.Context, []any, string, ...any) error
}

func initTable(db DBExecutor) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement := `
	CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		sale_price NUMERIC(10, 2),
		price NUMERIC(10, 2) NOT NULL
	);
	`
	_, err := db.Execute(ctx, statement)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
