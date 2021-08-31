package repo

import (
	"context"

	"github.com/vectorman1/analysis/analysis-api/domain/instrument/model"

	"github.com/vectorman1/analysis/analysis-api/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type overviewRepository interface {
	Insert(ctx context.Context, overview *model.InstrumentOverviewResponse) (bool, error)
	GetBySymbolUuid(ctx context.Context, uuid string) (*model.InstrumentOverviewResponse, error)
	Delete(ctx context.Context, uuid string) error
}

type OverviewRepository struct {
	mondodb *mongo.Database
}

func NewOverviewRepository(mongodb *mongo.Database) *OverviewRepository {
	return &OverviewRepository{
		mondodb: mongodb,
	}
}

func (r *OverviewRepository) Insert(ctx context.Context, overview *model.Overview) (bool, error) {
	_, err := r.mondodb.Collection(common.OverviewsCollection).
		InsertOne(ctx, overview)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *OverviewRepository) GetByInstrumentUuid(ctx context.Context, uuid string) (*model.Overview, error) {
	var overview model.Overview
	err := r.mondodb.Collection(common.OverviewsCollection).
		FindOne(ctx, bson.M{"instrument_uuid": uuid}).
		Decode(&overview)
	if err != nil {
		return nil, err
	}

	return &overview, nil
}

func (r *OverviewRepository) Delete(ctx context.Context, uuid string) error {
	_, err := r.mondodb.Collection(common.OverviewsCollection).
		DeleteOne(ctx, bson.M{"instrument_uuid": uuid})
	if err != nil {
		return err
	}

	return nil
}
