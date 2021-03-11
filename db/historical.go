package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vectorman1/analysis/analysis-api/model/db/documents"

	"github.com/vectorman1/analysis/analysis-api/common"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jackc/pgx"
)

type historicalRepository interface {
	InsertMany(ctx context.Context, list []documents.Historical) (int, error)
	GetBySymbolUuid(ctx context.Context, symbolUuid string, startDate time.Time, endDate time.Time) (*[]documents.Historical, error)
}

type HistoricalRepository struct {
	historicalRepository
	pgdb    *pgx.ConnPool
	mongodb *mongo.Database
}

func NewHistoricalRepository(db *pgx.ConnPool, mongodb *mongo.Database) *HistoricalRepository {
	return &HistoricalRepository{
		pgdb:    db,
		mongodb: mongodb,
	}
}

func (r *HistoricalRepository) InsertMany(ctx context.Context, list []documents.Historical) (int, error) {
	var e []interface{}
	for _, v := range list {
		e = append(e, v)
	}

	res, err := r.mongodb.Collection(common.HISTORIES_COLLECTION).
		InsertMany(ctx, e)
	if err != nil {
		return 0, err
	}

	return len(res.InsertedIDs), nil
}

func (r *HistoricalRepository) GetBySymbolUuid(ctx context.Context, symbolUuid string, startDate time.Time, endDate time.Time) (*[]documents.Historical, error) {
	opts := options.Find()
	opts.SetSort(bson.D{{"timestamp", -1}})
	filter := bson.M{
		"symboluuid": symbolUuid,
		"$and": []bson.M{
			{
				"timestamp": bson.M{
					"$lt": endDate,
				},
			},
			{
				"timestamp": bson.M{
					"$gt": startDate,
				},
			},
		},
	}

	filterCursor, err := r.mongodb.Collection(common.HISTORIES_COLLECTION).
		Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer filterCursor.Close(ctx)

	var result []documents.Historical
	for filterCursor.Next(ctx) {
		history := documents.Historical{}
		err := filterCursor.Decode(&history)
		if err != nil {
			return nil, err
		}
		result = append(result, history)
	}

	return &result, nil
}
