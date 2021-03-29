package db

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/vectorman1/analysis/analysis-api/common"
	"github.com/vectorman1/analysis/analysis-api/generated/proto_models"
	"github.com/vectorman1/analysis/analysis-api/model/db/entities"
)

type UserRepositoryContract interface {
	GetByUsername(context.Context, string) (*entities.User, error)
	GetPaged(context.Context, *proto_models.PagedFilter) (*[]entities.User, uint, error)
	Create(context.Context, *entities.User) error
	Update(context.Context, *entities.User) error
	Delete(context.Context, string) error
}

type UserRepository struct {
	db *pgx.ConnPool
}

func NewUserRepository(db *pgx.ConnPool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	// Find user with matching username
	query, args, err := squirrel.
		Select("*").
		From("\"user\".users").
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	var res entities.User
	row := r.db.QueryRow(query, args...)
	err = row.Scan(&res.ID, &res.Uuid, &res.PrivateRole, &res.Username, &res.Password, &res.CreatedAt, &res.UpdatedAt, &res.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *UserRepository) GetPaged(ctx context.Context, filter *proto_models.PagedFilter) (*[]entities.User, uint, error) {
	// generate query
	order := common.FormatOrderQuery(filter.Order, filter.Ascending)
	query, args, err := squirrel.
		Select("*, count(*) OVER() AS total_count").
		From("\"user\".users").
		OrderBy(order).
		Offset((filter.PageNumber - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Where("deleted_at is NULL").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryEx(ctx, query, &pgx.QueryExOptions{}, args...)
	if err != nil {
		return nil, 0, err
	}

	var result []entities.User
	var totalItems uint
	for rows.Next() {
		user := entities.User{}
		if err = rows.Scan(
			&user.ID,
			&user.Uuid,
			&user.PrivateRole,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&totalItems); err != nil {
			return nil, 0, err
		}
		result = append(result, user)
	}

	return &result, totalItems, nil
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query, args, err := squirrel.
		Insert("\"user\".users").
		Columns("uuid, privateRole, username, password, createdAt, updatedAt").
		Values(&user.Uuid, &user.PrivateRole, &user.Username, &user.Password, time.Now(), time.Now()).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecEx(ctx, query, &pgx.QueryExOptions{}, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	query, args, err := squirrel.
		Update("\"user\".users").
		Where(squirrel.Eq{"uuid": user.Uuid}).
		Set("username", user.Username).
		Set("password", user.Password).
		Set("updatedAt", time.Now()).
		Set("privateRole", user.PrivateRole).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecEx(ctx, query, &pgx.QueryExOptions{}, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, uuid string) error {
	query, args, err := squirrel.
		Update("\"user\".users").
		Where(squirrel.Eq{"uuid": uuid}).
		Set("deletedAt", time.Now()).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecEx(ctx, query, &pgx.QueryExOptions{}, args...)
	if err != nil {
		return err
	}

	return nil
}
