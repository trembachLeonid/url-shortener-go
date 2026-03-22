package persistence

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
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
		`INSERT INTO urls (original_url, expire_time) VALUES ($1, $2) RETURNING id`,
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

func GetURL(shortURL string) (originalURL string, expireTime *time.Time, err error) {
	conn, err := connect()
	if err != nil {
		return
	}
	defer conn.Close()

	originalURL = ""
	expireTime = nil
	err = conn.QueryRow("SELECT original_url, expire_time FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL, &expireTime)

	return
}

func ShortURLExists(shortURL string) (bool, error) {
	conn, err := connect()
	if err != nil {
		return true, err
	}
	defer conn.Close()

	exists := false
	err = conn.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE short_url = $1);", shortURL).Scan(&exists)

	return exists, err
}
