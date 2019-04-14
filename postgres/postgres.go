package postgres

import (
	"database/sql"
)

type EggchanService struct {
	DB *sql.DB
}
