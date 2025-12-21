package DAL

import (
	"database/sql"
)

type SQLDatabase interface {
	Connection() *sql.DB
	Close() error
}
