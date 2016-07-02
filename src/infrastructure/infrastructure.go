package infrastructure

import (
	"database/sql"
	"fmt"
	"interfaces"
)

type SqliteHandler struct {
	Conn *sql.DB
}

func (handler *SqliteHandler) Execute(statement string) {
	handler.Conn.Exec(statement)
}

func (handler *SqliteHandler) Query(statement string) interfaces.Row {
	rows, err := handler.Conn.Query(statement)
	if err != nil {
		fmt.Println(err)
		return new(SqliteRow)
	}
	row := new(SqliteRow)
	row.Rows = rows
	return row
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

func NewSqliteHandler(dbFileName string) *SqliteHandler {
	conn, _ := sql.Open("sqlite3", dbFileName)
	sqliteHandler := new(SqliteHandler)
	sqliteHandler.Conn = conn
	return sqliteHandler
}
