package config

import "database/sql"

type RuntimeConfig struct {
	DBhandler *sql.DB
}
