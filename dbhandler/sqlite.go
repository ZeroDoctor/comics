package dbhandler

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var sqlite *SqliteHandler = nil
var once sync.Once

type SqliteHandler struct {
	DB *sqlx.DB
}

func init() {
	once.Do(func() {
		sqlite = ConnectLite()
	})
}

func ConnectLite() *SqliteHandler {
	if sqlite != nil {
		return sqlite
	}

	db, err := sqlx.Connect("sqlite3", "./_uConfig.db")
	if err != nil {
		log.Fatalln(err)
	}

	sqlite = &SqliteHandler{
		DB: db,
	}

	log.Println("setup user config...")
	sqlite.UserConfig()

	return sqlite
}

func (sql SqliteHandler) UserConfig() {
	schema := `CREATE TABLE IF NOT EXISTS user_config (
		id             INTEGER DEFAULT 0,
		user           TEXT DEFAULT 'BOB_ROSS',
		rate_limit     INTEGER DEFAULT 90,
		rate_remaining INTEGER DEFAULT 90,
		last_request   DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := sql.DB.Exec(schema)
	if err != nil {
		log.Fatalln(err)
	}
}

func (sql SqliteHandler) UpdateRates(rateLimit, rateRemaining int, lastRequest time.Time) error {
	query := `DELETE FROM user_config; INSERT OR REPLACE INTO user_config (id, rate_limit, rate_remaining, last_request) 
		VALUES ('0', '` + strconv.Itoa(rateLimit) + `', '` + strconv.Itoa(rateRemaining) + `', '` + lastRequest.Format(time.RFC3339) + `');`

	_, err := sql.DB.Exec(query)
	return err
}

func (sql SqliteHandler) FetchRates() (int, int, time.Time, error) {
	var Rate struct {
		rateLimit     int       `db:"rate_limit"`
		rateRemaining int       `db:"rate_remaining"`
		lastRequest   time.Time `db:"last_request"`
	}

	err := sql.DB.Select(&Rate, "SELECT rate_limit, rate_remaining, last_request FROM user_config")
	if err != nil {
		return 90, 90, time.Now(), err
	}

	return Rate.rateLimit, Rate.rateRemaining, Rate.lastRequest, nil
}
