package db

import (
	"context"
	"time"

	"github.com/vectorman1/analysis/analysis-api/model/db"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
)

type historicalRepository interface {
	GetBySymbolUuid(ctx *context.Context, symbolUuid string, startDate time.Time, endDate time.Time) (*[]db.Historical, error)
}

type HistoricalRepository struct {
	historicalRepository
	db *pgx.ConnPool
}

func NewHistoricalRepository(db *pgx.ConnPool) *HistoricalRepository {
	return &HistoricalRepository{
		db: db,
	}
}

func (r *HistoricalRepository) GetBySymbolUuid(ctx *context.Context, symbolUuid string, startDate time.Time, endDate time.Time) (*[]db.Historical, error) {
	queryBuilder := squirrel.
		Select("*").
		From("analysis.historical").
		Where(squirrel.Eq{"symbol_uuid": symbolUuid}).
		Where(squirrel.LtOrEq{"for_date": endDate}).
		Where(squirrel.GtOrEq{"for_date": startDate}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryEx(*ctx, query, &pgx.QueryExOptions{}, args)
	if err != nil {
		return nil, err
	}

	var result []db.Historical
	for rows.Next() {
		e := db.Historical{}
		if err = rows.Scan(&e.ID, &e.SymbolUuid, &e.Values, &e.ForDate, &e.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, e)
	}

	return &result, nil
}
