package wire

import (
	"database/sql"

	db "github.com/leminhthai/train-ticket/train-service/db/generated"
)

func ProvideQueries(sqlDB *sql.DB) *db.Queries {
	return db.New(sqlDB)
}
