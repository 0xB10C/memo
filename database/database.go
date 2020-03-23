package database

import (
	"database/sql"

	"github.com/gomodule/redigo/redis"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// Pool is a pool of redis connections
	Pool *redis.Pool
	// SQLiteDB is a open SQLite database
	SQLiteDB *sql.DB
)
