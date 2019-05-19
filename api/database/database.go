package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/0xb10c/memo/api/config"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func SetupDatabase() error {

	dbUser := config.GetString("database.user")
	dbPasswd := config.GetString("database.passwd")
	dbHost := config.GetString("database.host")
	dbName := config.GetString("database.name")
	dbConnection := config.GetString("database.connection")
	connectionString := dbUser + ":" + dbPasswd + "@" + dbConnection + "(" + dbHost + ")" + "/" + dbName

	db, err := sql.Open("mysql", connectionString+"?parseTime=true")
	if err != nil {
		return err
	}

	// Ping the database once since Open() doesn't open a connection
	err = db.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Setup database connection")
	DB = db
	return nil
}

// GetMempool gets the current mempool from the database
func GetMempool() (timestamp time.Time, byCount string, megabyteMarkersJSON string, mempoolSize int, err error) {

	sqlStatement := "SELECT timestamp, byCount, positionsInGreedyBlocks, mempoolSize FROM current_mempool WHERE id = 1"
	row := DB.QueryRow(sqlStatement)

	err = row.Scan(&timestamp, &byCount, &megabyteMarkersJSON, &mempoolSize)
	if err != nil {
		if err == sql.ErrNoRows {
			return timestamp, byCount, megabyteMarkersJSON, mempoolSize, err
		} else {
			panic(err)
		}
	}

	return
}

type RecentBlock struct {
	time      time.Time
	Timestamp int64 `json:"timestamp"`
	Height    int   `json:"height"`
	TxCount   int   `json:"txCount"`
	Size      int   `json:"size"`
	Weight    int   `json:"weight"`
}

func GetRecentBlocks() (blocks []RecentBlock, err error) {
	sqlStatement := "SELECT height, timestamp, txCount, size, weight FROM pastBlocks ORDER BY height DESC LIMIT 10;"
	rows, err := DB.Query(sqlStatement)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		block := RecentBlock{}
		err = rows.Scan(&block.Height, &block.time, &block.TxCount, &block.Size, &block.Weight)
		if err != nil {
			// handle this error
			panic(err)
		}
		block.Timestamp = block.time.Unix()
		blocks = append(blocks, block)
	}
	return
}
