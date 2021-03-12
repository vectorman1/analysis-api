package db

import (
	"context"

	"github.com/vectorman1/analysis/analysis-api/common"
	"github.com/vectorman1/analysis/analysis-api/model/db/documents"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportRepo interface {
	InsertMany(report documents.Report) (int, error)
}

type ReportRepository struct {
	mongodb *mongo.Database
}

func (r *ReportRepository) InsertMany(ctx context.Context, list *[]documents.Report) (int, error) {
	var e []interface{}
	for _, v := range *list {
		e = append(e, v)
	}

	inserted, err := r.mongodb.Collection(common.REPORTS_COLLECTION).
		InsertMany(ctx, e)
	if err != nil {
		return 0, err
	}

	return len(inserted.InsertedIDs), nil
}
