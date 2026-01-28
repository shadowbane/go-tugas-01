package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func InitDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	zap.S().Info("Database connected")
	return db, nil
}
