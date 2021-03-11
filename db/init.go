package db

import (
	"time"

	"github.com/jackc/pgx"
	"github.com/vectorman1/analysis/analysis-api/common"
)

func GetConnPool(config *common.Config) (*pgx.ConnPool, error) {
	cfg := pgx.ConnConfig{
		User:     config.PostgreSQLConfig.DatastoreDBUser,
		Password: config.PostgreSQLConfig.DatastoreDBPassword,
		Database: config.PostgreSQLConfig.DatastoreDBSchema,
		Host:     config.PostgreSQLConfig.DatastoreDBHost,
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     cfg,
		MaxConnections: config.PostgreSQLConfig.DatabaseMaxConnections,
		AcquireTimeout: 5 * time.Second,
	}

	return pgx.NewConnPool(poolConfig)
}
