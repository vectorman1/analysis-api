package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vectorman1/analysis/analysis-api/model/db/documents"

	"github.com/vectorman1/analysis/analysis-api/common"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type historyRepository interface {
	InsertMany(ctx context.Context, list []documents.History) (int, error)
	GetSymbolHistory(ctx context.Context, symbolUuid string, startDate time.Time, endDate time.Time) (*[]documents.History, error)
	GetLastSymbolHistory(ctx context.Context, symbolUuid string) (*documents.History, error)
}

type HistoryRepository struct {
	historyRepository
	mongodb *mongo.Database
}

func NewHistoryRepository(mongodb *mongo.Database) *HistoryRepository {
	return &HistoryRepository{
		mongodb: mongodb,
	}
}

func (r *HistoryRepository) InsertMany(ctx context.Context, list []documents.History) (int, error) {
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

func (r *HistoryRepository) GetSymbolHistory(ctx context.Context, symbolUuid string, startDate time.Time, endDate time.Time) (*[]documents.History, error) {
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

	var result []documents.History
	for filterCursor.Next(ctx) {
		history := documents.History{}
		err := filterCursor.Decode(&history)
		if err != nil {
			return nil, err
		}
		result = append(result, history)
	}

	return &result, nil
}

func (r *HistoryRepository) GetLastSymbolHistory(ctx context.Context, symbolUuid string) (*documents.History, error) {
	opts := options.FindOne()
	opts.SetSort(bson.D{{"timestamp", -1}})
	filter := bson.M{
		"symboluuid": symbolUuid,
	}

	var res documents.History
	_ = r.mongodb.Collection(common.HISTORIES_COLLECTION).
		FindOne(ctx, filter, opts).Decode(&res)

	return &res, nil
}
