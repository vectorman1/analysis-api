package db

import (
	"github.com/Masterminds/squirrel"
	"github.com/dystopia-systems/alaskalog"
	"github.com/jackc/pgx"
	"github.com/vectorman1/analysis/analysis-api/model/db"
)

type currencyRepository interface {
	GetByCode(code string) (*db.Currency, error)
	GetOrCreate(code string) (*db.Currency, error)
	Create(curr *db.Currency) (uint, error)
}

type CurrencyRepository struct {
	currencyRepository
	db *pgx.ConnPool
}

func NewCurrencyRepository(pgDb *pgx.ConnPool) *CurrencyRepository {
	return &CurrencyRepository{
		db: pgDb,
	}
}

func (r *CurrencyRepository) GetByCode(code string) (*db.Currency, error) {
	queryBuilder := squirrel.
		Select("id, code, long_name").
		From("analysis.currencies").
		Where(squirrel.Eq{"code": code}).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows := r.db.QueryRow(query, args...)

	var res db.Currency
	err = rows.Scan(&res.ID, &res.Code, &res.LongName)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *CurrencyRepository) GetOrCreate(code string) (*db.Currency, error) {
	curr, err := r.GetByCode(code)
	if err != nil {
		newCurr := &db.Currency{}
		newCurr.Code = code
		newCurr.LongName = "temp name"

		id, createErr := r.Create(newCurr)
		if createErr != nil {
			alaskalog.Logger.Warnf("failed creating currency: %v", err)
			return nil, createErr
		}

		newCurr.ID = id
		return newCurr, nil
	} else {
		return curr, nil
	}
}

func (r *CurrencyRepository) Create(curr *db.Currency) (uint, error) {
	queryBuilder := squirrel.
		Insert("analysis.currencies").
		Columns("code, long_name").
		Values(curr.Code, curr.LongName).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	conn, err := r.db.Acquire()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	id := uint(0)
	err = conn.QueryRow(query+" RETURNING id;", args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
