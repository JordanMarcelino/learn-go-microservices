package repository

import (
	"context"

	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/entity"
)

type VerificationRepository interface {
	Save(ctx context.Context, verification *entity.Verification) error
}

type verificationRepositoryImpl struct {
	DB DBTX
}

func NewVerificationRepository(db DBTX) VerificationRepository {
	return &verificationRepositoryImpl{
		DB: db,
	}
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
