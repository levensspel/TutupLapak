package userRepository

import (
	"context"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/TutupLapak/user/src/model/dtos/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
)

type UserRepositoryInterface interface {
	CreateUserByEmail(ctx context.Context, pool *pgxpool.Pool, email, passwordHash string) (userId string, err error)
	CreateUserByPhone(ctx context.Context, pool *pgxpool.Pool, phone, passwordHash string) (userId string, err error)
	GetAuthByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (auth repository.AuthByEmail, err error)
	GetAuthByPhone(ctx context.Context, pool *pgxpool.Pool, phone string) (auth repository.AuthByPhone, err error)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}

func NewUserRepositoryInject(i do.Injector) (UserRepositoryInterface, error) {
	return NewUserRepository(
		do.MustInvoke[*pgxpool.Pool](i),
	), nil
}

func (ur *UserRepository) CreateUserByEmail(ctx context.Context, pool *pgxpool.Pool, email, passwordHash string) (userId string, err error) {
	query := `INSERT INTO users(email, password_hash) VALUES($1, $2) RETURNING id`

	row := pool.QueryRow(ctx, query, email, passwordHash)
	err = row.Scan(&userId)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (ur *UserRepository) CreateUserByPhone(ctx context.Context, pool *pgxpool.Pool, phone, passwordHash string) (userId string, err error) {
	query := `INSERT INTO users(phone, password_hash) VALUES($1, $2) RETURNING id`

	row := pool.QueryRow(ctx, query, phone, passwordHash)
	err = row.Scan(&userId)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (ur *UserRepository) GetAuthByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (auth repository.AuthByEmail, err error) {
	query := `SELECT id, password_hash, phone FROM users WHERE email = $1 ;`

	row := pool.QueryRow(ctx, query, email)

	var phone *string
	err = row.Scan(&auth.UserId, &auth.HashPassword, &phone)
	if err != nil {
		return repository.AuthByEmail{}, err
	}

	if phone != nil {
		auth.Phone = *phone
	}

	return auth, nil
}

func (ur *UserRepository) GetAuthByPhone(ctx context.Context, pool *pgxpool.Pool, phone string) (auth repository.AuthByPhone, err error) {
	query := `SELECT id, password_hash, email FROM users WHERE phone = $1 ;`

	row := pool.QueryRow(ctx, query, phone)

	var email *string
	err = row.Scan(&auth.UserId, &auth.HashPassword, &email)
	if err != nil {
		return repository.AuthByPhone{}, err
	}

	if email != nil {
		auth.Email = *email
	}

	return auth, nil
}
