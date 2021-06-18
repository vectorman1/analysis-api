package repo

import (
	"context"

	"github.com/vectorman1/analysis/analysis-api/domain/instrument/model"

	"github.com/vectorman1/analysis/analysis-api/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SymbolOverviewContract interface {
	Insert(ctx context.Context, overview *model.InstrumentOverviewResponse) (bool, error)
	GetBySymbolUuid(ctx context.Context, uuid string) (*model.InstrumentOverviewResponse, error)
	Delete(ctx context.Context, uuid string) error
}

type SymbolOverviewRepository struct {
	mondodb *mongo.Database
}

func NewSymbolOverviewRepository(mongodb *mongo.Database) *SymbolOverviewRepository {
	return &SymbolOverviewRepository{
		mondodb: mongodb,
	}
}

func (r *SymbolOverviewRepository) Insert(ctx context.Context, overview *model.InstrumentOverview) (bool, error) {
	_, err := r.mondodb.Collection(common.OverviewsCollection).
		InsertOne(ctx, overview)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *SymbolOverviewRepository) GetByInstrumentUuid(ctx context.Context, uuid string) (*model.InstrumentOverview, error) {
	var overview model.InstrumentOverview
	err := r.mondodb.Collection(common.OverviewsCollection).
		FindOne(ctx, bson.M{"symboluuid": uuid}).
		Decode(&overview)
	if err != nil {
		return nil, err
	}

	return &overview, nil
}

func (r *SymbolOverviewRepository) Delete(ctx context.Context, uuid string) error {
	_, err := r.mondodb.Collection(common.OverviewsCollection).
		DeleteOne(ctx, bson.M{"symboluuid": uuid})
	if err != nil {
		return err
	}

	return nil
}
