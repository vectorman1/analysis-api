package db

import (
	"context"
	"fmt"
	"time"

	"github.com/vectorman1/analysis/analysis-api/model/db/entities"

	"github.com/jackc/pgx/pgtype"

	"github.com/gofrs/uuid"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/vectorman1/analysis/analysis-api/common"
	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
)

type SymbolRepo interface {
	GetPaged(ctx context.Context, req *symbol_service.GetPagedRequest) (*[]entities.Symbol, uint, error)
	GetByUuid(ctx context.Context, uuid string) (*entities.Symbol, error)
	InsertBulk(tx *pgx.Tx, ctx context.Context, symbols []*entities.Symbol) (bool, error)
	DeleteBulk(tx *pgx.Tx, ctx context.Context, symbols []*entities.Symbol) (bool, error)
	UpdateBulk(tx *pgx.Tx, ctx context.Context, symbols []*entities.Symbol) (bool, error)

	BeginTx(ctx *context.Context, options *pgx.TxOptions) (*pgx.Tx, error)
}

type SymbolRepository struct {
	db *pgx.ConnPool
}

func NewSymbolRepository(db *pgx.ConnPool) *SymbolRepository {
	return &SymbolRepository{
		db: db,
	}
}

// GetPaged returns a paged response of symbols stored
func (r *SymbolRepository) GetPaged(ctx context.Context, req *symbol_service.GetPagedRequest) (*[]entities.Symbol, uint, error) {
	// generate query
	order := common.FormatOrderQuery(req.Filter.Order, req.Filter.Ascending)
	queryBuilder := squirrel.
		Select("*, count(*) OVER() AS total_count").
		From("analysis.symbols").
		OrderBy(order).
		Offset((req.Filter.PageNumber - 1) * req.Filter.PageSize).
		Limit(req.Filter.PageSize).
		Where("deletedAt is NULL").
		PlaceholderFormat(squirrel.Dollar)

	if req.Filter.Text != "" {
		req.Filter.Text = fmt.Sprintf("%%%s%%", req.Filter.Text)
		nameLikeText := squirrel.ILike{"name": req.Filter.Text}
		identifierLikeText := squirrel.ILike{"identifier": req.Filter.Text}
		isinLikeText := squirrel.ILike{"isin": req.Filter.Text}

		queryBuilder = queryBuilder.Where(squirrel.Or{nameLikeText, identifierLikeText, isinLikeText})
	}

	q, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryEx(ctx, q, nil, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// read all resulting rows
	var result []entities.Symbol
	var totalItems uint
	for rows.Next() {
		sym := entities.Symbol{}
		if err = rows.Scan(
			&sym.ID,
			&sym.Uuid,
			&sym.CurrencyCode,
			&sym.Isin,
			&sym.Identifier,
			&sym.Name,
			&sym.MinimumOrderQuantity,
			&sym.MarketName,
			&sym.MarketHoursGmt,
			&sym.CreatedAt,
			&sym.UpdatedAt,
			&sym.DeletedAt,
			&totalItems); err != nil {
			return nil, 0, err
		}
		result = append(result, sym)
	}

	return &result, totalItems, nil
}

func (r *SymbolRepository) GetByUuid(ctx context.Context, symbolUuid string) (*entities.Symbol, error) {
	u, err := uuid.FromString(symbolUuid)
	if err != nil {
		return nil, err
	}

	queryBuilder := squirrel.
		Select("*").
		From("analysis.symbols").
		Where(fmt.Sprintf("uuid = '%s'", u.String())).
		Limit(1)
	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	sym := entities.Symbol{}
	row := r.db.QueryRowEx(ctx, query, &pgx.QueryExOptions{})
	if err = row.Scan(
		&sym.ID,
		&sym.Uuid,
		&sym.CurrencyCode,
		&sym.Isin,
		&sym.Identifier,
		&sym.Name,
		&sym.MinimumOrderQuantity,
		&sym.MarketName,
		&sym.MarketHoursGmt,
		&sym.CreatedAt,
		&sym.UpdatedAt,
		&sym.DeletedAt); err != nil {
		return nil, err
	}

	return &sym, nil
}

// InsertBulk inserts the slice in a single transaction in batches and returns success and error
func (r *SymbolRepository) InsertBulk(tx *pgx.Tx, ctx context.Context, symbols []*entities.Symbol) (bool, error) {
	// split inserts in batches
	workList := make(chan []*entities.Symbol)
	go func() {
		defer close(workList)
		batchSize := 1000
		var stack []*entities.Symbol
		for _, sym := range symbols {
			stack = append(stack, sym)
			if len(stack) == batchSize {
				workList <- stack
				stack = nil
			}
		}
		if len(stack) > 0 {
			workList <- stack
		}
	}()

	now := time.Now() // generate query for insert from batches
	for list := range workList {
		q := squirrel.
			Insert("analysis.symbols").
			Columns("uuid, currencyCode, isin, identifier, name, minimumOrderQuantity, marketName, marketHoursGmt, createdAt, updatedAt, deletedAt").
			PlaceholderFormat(squirrel.Dollar)
		for _, sym := range list {
			q = q.Values(
				&sym.Uuid,
				&sym.CurrencyCode,
				&sym.Isin,
				&sym.Identifier,
				&sym.Name,
				&sym.MinimumOrderQuantity,
				&sym.MarketName,
				&sym.MarketHoursGmt,
				now,
				now,
				&pgtype.Timestamptz{Status: pgtype.Null})
		}

		query, args, _ := q.ToSql()
		if len(args) > 0 {
			_, err := tx.ExecEx(ctx, query, &pgx.QueryExOptions{}, args...)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

// DeleteBulk sets the Deleted At values for bulk symbols to now
func (r *SymbolRepository) DeleteBulk(tx *pgx.Tx, ctx context.Context, symbols []*entities.Symbol) (bool, error) {
	// split updates in batches
	workList := make(chan []*entities.Symbol)
	go func() {
		defer close(workList)
		batchSize := 1000
		var stack []*entities.Symbol
		for _, sym := range symbols {
			stack = append(stack, sym)
			if len(stack) == batchSize {
				workList <- stack
				stack = nil
			}
		}
		if len(stack) > 0 {
			workList <- stack
		}
	}()

	now := time.Now()
	// generate query for update from batches
	for list := range workList {
		for _, sym := range list {
			var u string
			sym.Uuid.AssignTo(&u)

			q := squirrel.Update("analysis.symbols")

			q = q.
				Set("deletedAt", now).
				PlaceholderFormat(squirrel.Dollar).
				Where(squirrel.Eq{"uuid::text": u})

			query, args, _ := q.ToSql()
			if len(args) > 0 {
				_, err := tx.ExecEx(ctx, query, &pgx.QueryExOptions{}, args...)
				if err != nil {
					return false, err
				}
			}
		}
	}

	return true, nil
}

// UpdateBulk updates all columns of the symbol with the matching uuid
// with the passed symbol values
func (r *SymbolRepository) UpdateBulk(tx *pgx.Tx, ctx context.Context, symbols []*entities.Symbol) (bool, error) {
	// split updates in batches
	workList := make(chan []*entities.Symbol)
	go func() {
		defer close(workList)
		batchSize := 1000
		var stack []*entities.Symbol
		for _, sym := range symbols {
			stack = append(stack, sym)
			if len(stack) == batchSize {
				workList <- stack
				stack = nil
			}
		}
		if len(stack) > 0 {
			workList <- stack
		}
	}()

	now := time.Now()
	for list := range workList {
		for _, sym := range list {
			var u string
			sym.Uuid.AssignTo(&u)

			q := squirrel.
				Update("analysis.symbols").
				PlaceholderFormat(squirrel.Dollar)

			q = q.
				Set("name", sym.Name).
				Set("marketHoursGmt", sym.MarketHoursGmt).
				Set("updatedAt", now).
				Where(squirrel.Eq{"uuid::text": u})

			query, args, _ := q.ToSql()
			if len(args) > 0 {
				_, err := tx.ExecEx(ctx, query, &pgx.QueryExOptions{}, args...)
				if err != nil {
					return false, err
				}
			}
		}
	}

	return true, nil
}

// BeginTx starts a new transaction on the given context
func (r *SymbolRepository) BeginTx(ctx *context.Context, options *pgx.TxOptions) (*pgx.Tx, error) {
	tx, err := r.db.BeginEx(*ctx, options)
	if err != nil {
		return nil, err
	}

	return tx, err
}
