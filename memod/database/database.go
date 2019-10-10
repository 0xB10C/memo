package database

import (
	"database/sql"

	"github.com/gomodule/redigo/redis"
	_ "github.com/mattn/go-sqlite3"

	"github.com/0xb10c/memo/memod/config"
	"github.com/0xb10c/memo/memod/logger"
)

var (
	// Pool is a pool of redis connections
	Pool *redis.Pool
	// SQLiteDB is a open SQLite database
	SQLiteDB *sql.DB
)

func newPool() *redis.Pool {
	dbUser := config.GetString("redis.user")
	dbPasswd := config.GetString("redis.passwd")
	dbHost := config.GetString("redis.host")
	dbPort := config.GetString("redis.port")
	dbConnection := config.GetString("redis.connection")

	connectionString := dbConnection + "://" + dbUser + ":" + dbPasswd + "@" + dbHost + ":" + dbPort

	return &redis.Pool{
		MaxIdle:   40,
		MaxActive: 1200, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(connectionString)
			if err != nil {
				return c, err
			}
			return c, err
		},
	}
}

// SetupRedis sets up a new Redis Pool
func SetupRedis() (err error) {
	Pool = newPool()

	c := Pool.Get() // get a new connection
	defer c.Close()

	_, err = c.Do("PING")
	if err != nil {
		return err
	}

	logger.Info.Println("Setup redis database connection pool")

	return nil
}

// SetupSQLite sets up the SQLite Database used for persitently saving mempool entries
func SetupSQLite() (db *sql.DB, err error) {
	filePath := config.GetString("zmq.saveMempoolEntries.dbPath")
	db, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	SQLiteDB = db
	createMempoolEntryTableIfNotExists()
	logger.Info.Println("Setup SQLite database in " + filePath)

	return db, nil
}

func createMempoolEntryTableIfNotExists() {
	statement := `CREATE TABLE IF NOT EXISTS mempoolEntries (
		entryTime 			INT 			NOT NULL,
		txid 						TEXT 			NOT NULL,
		fee							INT				NOT NULL,
		size						INT				NOT NULL,
		inputs					INT				NOT NULL,
		outputs					INT				NOT NULL,
		locktime  			INT				NOT NULL,
		outSum					INT				NOT NULL,
		spendsSegWit  	BOOLEAN		NOT NULL,
		spendsMultisig  BOOLEAN   NOT NULL,
		bip69compliant	BOOLEAN 	NOT NULL,
		signalsRBF 			BOOLEAN 	NOT NULL,
		spends 					TEXT 			NOT NULL,
		paysto 					TEXT 			NOT NULL,
		multisigs 			TEXT,
		opreturndata 		TEXT,
		PRIMARY KEY (entryTime, txid)
	)`

	_, err := SQLiteDB.Exec(statement)
	if err != nil {
		logger.Error.Printf("Cound not create table mempoolEntries: %v.\n", err)
		return
	}
}
