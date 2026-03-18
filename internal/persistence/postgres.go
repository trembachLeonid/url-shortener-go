package persistence

import (
	"database/sql"
	"os"
	"time"
)

func connect() (*sql.DB, error) {
	conn, err := sql.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))

	return conn, err
}

func InsertURL(originalURL string, expireTime *time.Time) (int32, error) {
	conn, err := connect()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	id := int32(0)
	err = conn.QueryRow(
		`INSERT INTO urls (original_url, expiration_time) VALUES ($1, $2) RETURNING id`,
		originalURL, expireTime).Scan(&id)

	return id, err
}

func SetShortURL(id int32, shortURL string) error {
	conn, err := connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.QueryRow("UPDATE urls SET short_url = $1 WHERE id = $2", shortURL, id).Err()

	return err
}
