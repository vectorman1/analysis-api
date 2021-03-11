package db

import (
	"context"

	"github.com/vectorman1/analysis/analysis-api/model/db/documents"

	"github.com/vectorman1/analysis/analysis-api/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type symbolOverviewRepository interface {
	Insert(ctx context.Context, overview *documents.SymbolOverview) (bool, error)
	GetBySymbolUuid(ctx context.Context, uuid string) (*documents.SymbolOverview, error)
}

type SymbolOverviewRepository struct {
	symbolOverviewRepository
	mondodb *mongo.Database
}

func NewSymbolOverviewRepository(mongodb *mongo.Database) *SymbolOverviewRepository {
	return &SymbolOverviewRepository{
		mondodb: mongodb,
	}
}

func (r *SymbolOverviewRepository) Insert(ctx context.Context, overview *documents.SymbolOverview) (bool, error) {
	_, err := r.mondodb.Collection(common.OVERVIEWS_COLLECTION).
		InsertOne(ctx, overview)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *SymbolOverviewRepository) GetBySymbolUuid(ctx context.Context, uuid string) (*documents.SymbolOverview, error) {
	var overview documents.SymbolOverview
	err := r.mondodb.Collection(common.OVERVIEWS_COLLECTION).
		FindOne(ctx, bson.M{"symboluuid": uuid}).
		Decode(&overview)
	if err != nil {
		return nil, err
	}
	return &overview, nil
}
