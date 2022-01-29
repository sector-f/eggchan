package postgres

import (
	"database/sql"
	"fmt"
)

type EggchanService struct {
	db *sql.DB
}

func New(options Options) (*EggchanService, error) {
	// TODO: figure out why I apparently wasn't passing a username or password before
	connectionString := fmt.Sprintf(
		"host=%s dbname=%s sslmode=disable",
		options.Hostname,
		options.Database,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error establishing database connection: %w", err)
	}

	return &EggchanService{db}, nil
}

type Options struct {
	Hostname string
	Database string
	Username string
	Password string
}
