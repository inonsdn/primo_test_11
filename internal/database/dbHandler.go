package database

import "context"

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

type DBHandler struct {
	Dbx DBExecutor
}

func NewDBHandler(db DBExecutor) *DBHandler {
	db.Connect()
	dbHandler := DBHandler{
		Dbx: db,
	}

	return &dbHandler
}
