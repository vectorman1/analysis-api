package db

import (
	"time"

	"github.com/jackc/pgx"
	"github.com/vectorman1/analysis/analysis-api/common"
)

func GetConnPool(config *common.Config) (*pgx.ConnPool, error) {
	cfg := pgx.ConnConfig{
		User:     config.DatastoreDBUser,
		Password: config.DatastoreDBPassword,
		Database: config.DatastoreDBSchema,
		Host:     config.DatastoreDBHost,
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     cfg,
		MaxConnections: config.DatabaseMaxConnections,
		AcquireTimeout: 5 * time.Second,
	}

	return pgx.NewConnPool(poolConfig)
}
