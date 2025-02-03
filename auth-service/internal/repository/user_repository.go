package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/entity"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	VerifyByUserID(ctx context.Context, userID int64) error
}

type userRepositoryImpl struct {
	DB DBTX
}

func NewUserRepository(db DBTX) UserRepository {
	return &userRepositoryImpl{
		DB: db,
	}
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT
			id, hash_password, is_verified
		FROM
			users
		WHERE
			email = $1
	`

	user := &entity.User{Email: email}
	if err := r.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.HashPassword, &user.IsVerified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *entity.User) error {
	query := `
	INSERT INTO
		users(email, hash_password)
	VALUES
		($1, $2)
	RETURNING
		id, created_at, updated_at
	`

	return r.DB.QueryRowContext(ctx, query, user.Email, user.HashPassword).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepositoryImpl) VerifyByUserID(ctx context.Context, userID int64) error {
	query := `
		UPDATE
			users
		SET
			is_verified = true
		WHERE
			id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, userID)
	return err
}
