package database

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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

func convertValue(v any) any {
	switch val := v.(type) {
	case pgtype.Numeric:
		if val.Int == nil {
			return nil
		}
		f, _ := new(big.Float).SetInt(val.Int).Float64()
		if val.Exp != 0 {
			f = f * math.Pow10(int(val.Exp))
		}
		// convert to float64
		return f

	default:
		return v
	}
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
