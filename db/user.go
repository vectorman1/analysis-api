package db

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/vectorman1/analysis/analysis-api/model/db/entities"
	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	Create(user *entities.User) (*entities.User, error)
	Get(username string, password string) (*entities.User, error)
}

type UserRepository struct {
	userRepository
	db *pgx.ConnPool
}

func NewUserRepository(db *pgx.ConnPool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *entities.User) error {
	query, args, err := squirrel.
		Insert("\"user\".users").
		Columns("uuid, private_role, username, password, created_at, updated_at").
		Values(&user.Uuid, &user.PrivateRole, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Get(username string, password string) (*entities.User, error) {
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

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &res, nil
}
