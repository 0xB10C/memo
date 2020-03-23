package database

import (
	"database/sql"

	"github.com/0xb10c/memo/config"
	"github.com/0xb10c/memo/logger"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDB holds a connection to a SQLite database
var SQLiteDB *sql.DB

// SetupSQLite sets up the SQLite Database used for persistently saving mempool entries
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
		entryTime	INT	NOT NULL,
		txid		TEXT	NOT NULL,
		fee			INT		NOT NULL,
		size		INT		NOT NULL,
		version 	INT 	NOT NULL,
		inputs		INT		NOT NULL,
		outputs		INT		NOT NULL,
		locktime  	INT		NOT NULL,
		outSum		INT		NOT NULL,
		spendsSegWit  	BOOLEAN	NOT NULL,
		spendsMultisig  BOOLEAN	NOT NULL,
		bip69compliant	BOOLEAN	NOT NULL,
		signalsRBF 	BOOLEAN	NOT NULL,
		spends		TEXT	NOT NULL,
		paysto		TEXT	NOT NULL,
		multisigs	TEXT,
		opreturndata	TEXT,
		PRIMARY KEY (entryTime, txid)
	)`

	_, err := SQLiteDB.Exec(statement)
	if err != nil {
		logger.Error.Printf("Cound not create table mempoolEntries: %v.\n", err)
		return
	}
}
