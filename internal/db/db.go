package db

import (
	"database/sql"
	"interview/internal/models"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewDB(connString string) *DB {
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return &DB{conn: conn}
}
func (db *DB) InsertStructuredLog(log *models.StructuredLog) error {
	query := `INSERT INTO structured_logs (status_code, api_endpoint, message, timestamp, ip_address)
              VALUES ($1, $2, $3, $4, $5)`
	_, err := db.conn.Exec(query, log.StatusCode, log.API, log.Message, log.Timestamp, log.IPAddress)
	return err
}
