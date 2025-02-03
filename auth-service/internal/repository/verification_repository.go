package repository

import (
	"context"

	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/entity"
)

type VerificationRepository interface {
	FindByUserID(ctx context.Context, userID int64) (*entity.Verification, error)
	Save(ctx context.Context, verification *entity.Verification) error
	DeleteByUserID(ctx context.Context, userID int64) error
}

type verificationRepositoryImpl struct {
	DB DBTX
}

func NewVerificationRepository(db DBTX) VerificationRepository {
	return &verificationRepositoryImpl{
		DB: db,
	}
}

func (r *verificationRepositoryImpl) FindByUserID(ctx context.Context, userID int64) (*entity.Verification, error) {
	query := `
		SELECT
			id, token, expire_at
		FROM
			user_verifications
		WHERE
			user_id = $1
	`

	verification := &entity.Verification{UserID: userID}
	if err := r.DB.QueryRowContext(ctx, query, userID).Scan(&verification.ID, &verification.Token, &verification.ExpireAt); err != nil {
		return nil, err
	}
	return verification, nil
}

func (r *verificationRepositoryImpl) Save(ctx context.Context, verification *entity.Verification) error {
	query := `
		INSERT INTO
			user_verifications(user_id, token, expire_at)
		VALUES
			($1, $2, $3)
		RETURNING
			id
	`

	return r.DB.QueryRowContext(ctx, query, verification.UserID, verification.Token, verification.ExpireAt).Scan(&verification.ID)
}

func (r *verificationRepositoryImpl) DeleteByUserID(ctx context.Context, userID int64) error {
	query := `
		DELETE FROM
			user_verifications
		WHERE
			user_id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, userID)
	return err
}
