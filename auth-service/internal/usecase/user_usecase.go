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
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/jwtutils"
)

type UserUseCase interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Verify(ctx context.Context, req *dto.VerificationRequest) (*dto.VerificationResponse, error)
}

type userUseCaseImpl struct {
	Hasher                   encryptutils.Hasher
	JwtUtil                  jwtutils.JwtUtil
	DataStore                repository.DataStore
	SendVerificationProducer mq.AMQPProducer
}

func NewUserUseCase(
	hasher encryptutils.Hasher,
	JwtUtil jwtutils.JwtUtil,
	dataStore repository.DataStore,
	sendVerificationProducer mq.AMQPProducer,
) UserUseCase {
	return &userUseCaseImpl{
		Hasher:                   hasher,
		JwtUtil:                  JwtUtil,
		DataStore:                dataStore,
		SendVerificationProducer: sendVerificationProducer,
	}
}

func (u *userUseCaseImpl) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	res := new(dto.LoginResponse)
	err := u.DataStore.Atomic(ctx, func(ds repository.DataStore) error {
		userRepository := ds.UserRepository()

		user, err := userRepository.FindByEmail(ctx, req.Email)
		if err != nil {
			return err
		}

		if user == nil {
			return httperror.NewInvalidCredentialError()
		}

		if ok := u.Hasher.Check(req.Password, user.HashPassword); !ok {
			return httperror.NewInvalidCredentialError()
		}
		if !user.IsVerified {
			return httperror.NewUserNotVerifiedError()
		}

		token, err := u.JwtUtil.Sign(user.ID)
		if err != nil {
			return err
		}

		res.AccessToken = token
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
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

func (u *userUseCaseImpl) Verify(ctx context.Context, req *dto.VerificationRequest) (*dto.VerificationResponse, error) {
	res := new(dto.VerificationResponse)
	err := u.DataStore.Atomic(ctx, func(ds repository.DataStore) error {
		userRepository := ds.UserRepository()
		verificationRepository := ds.VerificationRepository()

		user, err := userRepository.FindByEmail(ctx, req.Email)
		if user == nil {
			return httperror.NewUserNotFoundError()
		}
		if err != nil {
			return err
		}
		if user.IsVerified {
			return httperror.NewUserAlreadyVerifiedError()
		}

		verification, err := verificationRepository.FindByUserID(ctx, user.ID)
		if err != nil {
			return err
		}

		if verification.ExpireAt.Before(time.Now()) {
			return httperror.NewTokenExpiredError()
		}
		if verification.Token != req.VerificationToken {
			return httperror.NewTokenWrongError()
		}

		user.IsVerified = true
		if err := userRepository.VerifyByUserID(ctx, user.ID); err != nil {
			return err
		}
		if err := verificationRepository.DeleteByUserID(ctx, verification.UserID); err != nil {
			return err
		}

		res = dto.ToVerificationResponse(user)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
