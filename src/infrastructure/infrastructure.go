package infrastructure

import (
	"database/sql"
	//	"fmt"
	"interfaces"
)

type SqliteHandler struct {
	Conn *sql.DB
}

func (handler *SqliteHandler) Execute(statement string) {
	handler.Conn.Exec(statement)
}

func (handler *SqliteHandler) Query(statement string) interfaces.Row {
	return nil
}

type SqliteRow struct {
	Rows *sql.Rows
}

func (row SqliteRow) Next() bool {
	return row.Rows.Next()
}

func (row SqliteRow) Scan(dest ...interface{}) {
	row.Rows.Scan(dest...)
}
