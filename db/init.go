package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

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

func CreateMongoIndexes(db *mongo.Database) error {
	historiesSymbolUuidTimestamp := mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "symboluuid",
				Value: -1,
			},
			{
				Key:   "timestamp",
				Value: -1,
			},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := db.
		Collection(common.HistoriesCollection).
		Indexes().
		CreateOne(context.Background(), historiesSymbolUuidTimestamp)

	return err
}
