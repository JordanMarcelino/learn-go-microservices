package repository

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
)

type LockRepository interface {
	Get(ctx context.Context, key string, ttl time.Duration, opt *redislock.Options) (*redislock.Lock, error)
}

type LockRepositoryImpl struct {
	rdl *redislock.Client
}

func NewLockRepository(rdl *redislock.Client) LockRepository {
	return &LockRepositoryImpl{
		rdl: rdl,
	}
}

func (r *LockRepositoryImpl) Get(ctx context.Context, key string, ttl time.Duration, opt *redislock.Options) (*redislock.Lock, error) {
	lock, err := r.rdl.Obtain(ctx, key, ttl, opt)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			return nil, httperror.NewRequestDuplicateError()
		}
		return nil, err
	}

	return lock, nil
}
