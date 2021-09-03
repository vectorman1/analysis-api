package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/vectorman1/analysis/analysis-api/domain/instrument/model"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vectorman1/analysis/analysis-api/common"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type HistoryRepositoryContract interface {
	InsertMany(ctx context.Context, list *[]model.History) (int, error)
	GetSymbolHistory(ctx context.Context, symbolUuid string, startDate time.Time, endDate time.Time, desc bool) ([]model.History, error)
	GetInstrumentHistoryEntries(ctx context.Context, instrumentUuid string, n int, timeDesc bool) ([]model.History, error)
	GetLastSymbolHistory(ctx context.Context, symbolUuid string) (*model.LastHistory, error)
}

type HistoryRepository struct {
	mongodb *mongo.Database
}

func NewHistoryRepository(mongodb *mongo.Database) *HistoryRepository {
	return &HistoryRepository{
		mongodb: mongodb,
	}
}

func (r *HistoryRepository) InsertMany(ctx context.Context, list *[]model.History) (int, error) {
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

func (r *HistoryRepository) GetSymbolHistory(ctx context.Context, instrumentUuid string, startDate time.Time, endDate time.Time, timeDesc bool) ([]model.History, error) {
	opts := options.Find()
	if timeDesc {
		opts.SetSort(bson.D{{"timestamp", -1}})
	} else {
		opts.SetSort(bson.D{{"timestamp", 1}})
	}
	filter := bson.M{
		"symboluuid": instrumentUuid,
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

	var result []model.History
	for filterCursor.Next(ctx) {
		history := model.History{}
		err := filterCursor.Decode(&history)
		if err != nil {
			return nil, err
		}
		result = append(result, history)
	}

	return result, nil
}

// GetSymbolHistoryEntries returns N count of history entries per instrument
func (r *HistoryRepository) GetInstrumentHistoryEntries(ctx context.Context, instrumentUuid string, n int, timeDesc bool) ([]model.History, error) {
	opts := options.Find()
	if timeDesc {
		opts.SetSort(bson.D{{"timestamp", -1}})
	} else {
		opts.SetSort(bson.D{{"timestamp", 1}})
	}

	filter := bson.M{
		"instrumentUuid": instrumentUuid,
		"$limit":         n,
	}

	filterCursor, err := r.mongodb.Collection(common.HistoriesCollection).
		Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer filterCursor.Close(ctx)

	var result []model.History
	for filterCursor.Next(ctx) {
		history := model.History{}
		err := filterCursor.Decode(&history)
		if err != nil {
			return nil, err
		}
		result = append(result, history)
	}

	return result, nil
}

func (r *HistoryRepository) GetLastSymbolHistory(ctx context.Context, instrumentUuid string) (*model.LastHistory, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"symboluuid": instrumentUuid,
			},
		},
		{
			"$sort": bson.M{
				"timestamp": -1,
			},
		},
		{
			"$limit": 1,
		},
		{
			"$group": bson.M{
				"_id": "$_id",
				"close": bson.M{
					"$last": "$close",
				},
				"timestamp": bson.M{
					"$last": "$timestamp",
				},
			},
		},
	}

	var res model.LastHistory
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
