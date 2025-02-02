package usecase

import (
	"context"
	"time"

	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/entity"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/httperror"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/repository"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/utils/tokenutils"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/encryptutils"
)

type UserUseCase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
}

type userUseCaseImpl struct {
	Hasher                   encryptutils.Hasher
	DataStore                repository.DataStore
	SendVerificationProducer mq.AMQPProducer
}

func NewUserUseCase(
	hasher encryptutils.Hasher,
	dataStore repository.DataStore,
	sendVerificationProducer mq.AMQPProducer,
) UserUseCase {
	return &userUseCaseImpl{
		Hasher:                   hasher,
		DataStore:                dataStore,
		SendVerificationProducer: sendVerificationProducer,
	}
}

func (u *userUseCaseImpl) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	res := new(dto.RegisterResponse)
	err := u.DataStore.Atomic(ctx, func(ds repository.DataStore) error {
		userRepository := ds.UserRepository()
		verificationRepository := ds.VerificationRepository()

		user, err := userRepository.FindByEmail(ctx, req.Email)
		if user != nil {
			return httperror.NewUserAlreadyExistError()
		}
		if err != nil {
			return err
		}

		hashedPassword, err := u.Hasher.Hash(req.Password)
		if err != nil {
			return err
		}

		user = &entity.User{Email: req.Email, HashPassword: hashedPassword}
		if err := userRepository.Save(ctx, user); err != nil {
			return err
		}

		verification := &entity.Verification{UserID: user.ID, Token: tokenutils.GenerateOTPCode(), ExpireAt: time.Now().Add(constant.VerificationTimeout)}
		if err := verificationRepository.Save(ctx, verification); err != nil {
			return err
		}

		if err := u.SendVerificationProducer.Send(ctx, dto.SendVerificationEvent{Email: user.Email, Token: verification.Token}); err != nil {
			return err
		}

		res = dto.ToRegisterResponse(user)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
