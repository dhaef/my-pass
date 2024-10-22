package model

import (
	"database/sql"
	"net/http"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(conn *sql.DB) Database {
	return Database{
		conn: conn,
	}
}

func GetDBFromCtx(r *http.Request) Database {
	conn := r.Context().Value("dbConn").(*sql.DB)
	return NewDatabase(conn)
}
