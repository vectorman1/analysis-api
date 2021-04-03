package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vectorman1/analysis/analysis-api/model/db/documents"

	"github.com/vectorman1/analysis/analysis-api/common"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type HistoryRepositoryContract interface {
	InsertMany(ctx context.Context, list *[]documents.History) (int, error)
	GetSymbolHistory(ctx context.Context, symbolUuid string, startDate time.Time, endDate time.Time, desc bool) (*[]documents.History, error)
	GetLastSymbolHistory(ctx context.Context, symbolUuid string) (*documents.LastHistory, error)
}

type HistoryRepository struct {
	mongodb *mongo.Database
}

func NewHistoryRepository(mongodb *mongo.Database) *HistoryRepository {
	return &HistoryRepository{
		mongodb: mongodb,
	}
}

func (r *HistoryRepository) InsertMany(ctx context.Context, list *[]documents.History) (int, error) {
	var e []interface{}
	for _, v := range *list {
		e = append(e, v)
	}

	res, err := r.mongodb.Collection(common.HistoriesCollection).
		InsertMany(ctx, e)
	if err != nil {
		return 0, err
	}

	return len(res.InsertedIDs), nil
}

func (r *HistoryRepository) GetSymbolHistory(ctx context.Context, symbolUuid string, startDate time.Time, endDate time.Time, desc bool) (*[]documents.History, error) {
	opts := options.Find()
	if desc {
		opts.SetSort(bson.D{{"timestamp", -1}})
	} else {
		opts.SetSort(bson.D{{"timestamp", 1}})
	}
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

	filterCursor, err := r.mongodb.Collection(common.HistoriesCollection).
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

func (r *HistoryRepository) GetLastSymbolHistory(ctx context.Context, symbolUuid string) (*documents.LastHistory, error) {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"symboluuid": symbolUuid,
			},
		},
		bson.M{
			"$group": bson.M{
				"close": bson.M{
					"$last": "$close",
				},
				"timestamp": bson.M{
					"$last": "$timestamp",
				},
			},
		},
		bson.M{
			"$sort": bson.M{
				"timestamp": -1,
			},
		},
	}

	var res documents.LastHistory
	curr, err := r.mongodb.Collection(common.HistoriesCollection).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer curr.Close(ctx)

	if curr.Next(ctx) {
		err = curr.Decode(&res)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no next in cursor")
	}

	return &res, nil
}
