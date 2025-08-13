package service

import (
	"context"
	"database/sql"
	"os"
	"pisondev/markdown-notes-api/exception"
	"pisondev/markdown-notes-api/helper"
	"pisondev/markdown-notes-api/model/domain"
	"pisondev/markdown-notes-api/model/web"
	"pisondev/markdown-notes-api/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
	Log            *logrus.Logger
}

func NewUserService(userRepository repository.UserRepository, DB *sql.DB, validate *validator.Validate, log *logrus.Logger) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
		Log:            log,
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, req web.UserAuthRequest) (web.UserRegisterResponse, error) {
	s.Log.Info("validating req struct...")
	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Errorf("validation error: %v", err)
		return web.UserRegisterResponse{}, err
	}
	s.Log.Info("begin transaction...")
	tx, err := s.DB.Begin()
	if err != nil {
		s.Log.Errorf("failed to begin tx: %v", err)
		return web.UserRegisterResponse{}, err
	}

	s.Log.Info("generate hashed password...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Errorf("failed to hash password: %v", err)
		return web.UserRegisterResponse{}, err
	}

	_, err = s.UserRepository.FindByUsername(ctx, tx, req.Username)
	if err != sql.ErrNoRows {
		s.Log.Errorf("username already exist")
		return web.UserRegisterResponse{}, exception.ErrConflictUser
	}

	user := domain.User{
		Username:       req.Username,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now().UTC().Truncate(time.Second),
	}

	s.Log.Info("calling register repository...")
	savedUser, err := s.UserRepository.Register(ctx, tx, user)
	if err != nil {
		s.Log.Errorf("failed to use register repository in service layer: %v", err)
		err := tx.Rollback()
		if err != nil {
			s.Log.Errorf("failed to rollback: %v", err)
		}
		return web.UserRegisterResponse{}, err
	}

	s.Log.Info("commit transaction...")
	err = tx.Commit()
	if err != nil {
		s.Log.Errorf("failed to commit tx: %v", err)
		return web.UserRegisterResponse{}, err
	}

	return helper.ToUserRegisterResponse(savedUser), nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req web.UserAuthRequest) (web.UserLoginResponse, error) {
	err := s.Validate.Struct(req)
	if err != nil {
		s.Log.Errorf("validation error: %v", err)
		return web.UserLoginResponse{}, err
	}

	s.Log.Info("begin transaction...")
	tx, err := s.DB.Begin()
	if err != nil {
		s.Log.Errorf("failed to begin tx: %v", err)
		return web.UserLoginResponse{}, err
	}

	s.Log.Info("calling findbyusername repository...")
	selectedUser, err := s.UserRepository.FindByUsername(ctx, tx, req.Username)
	if err != nil {
		s.Log.Errorf("failed to use find repository in service layer: %v", err)
		err := tx.Rollback()
		if err != nil {
			s.Log.Errorf("failed to rollback tx: %v", err)
		}
		return web.UserLoginResponse{}, nil
	}

	s.Log.Info("compare hashed password...")
	err = bcrypt.CompareHashAndPassword([]byte(selectedUser.HashedPassword), []byte(req.Password))
	if err != nil {
		s.Log.Errorf("password doesn't match: %v", err)
		err := tx.Rollback()
		if err != nil {
			return web.UserLoginResponse{}, err
		}
		return web.UserLoginResponse{}, err
	}

	claims := web.CustomClaims{
		UserID: selectedUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		s.Log.Errorf("failed sign token: %v", err)
		return web.UserLoginResponse{}, err
	}

	s.Log.Info("commit transaction...")
	errCommit := tx.Commit()
	if errCommit != nil {
		s.Log.Errorf("failed to commit tx: %v", errCommit)
		return web.UserLoginResponse{}, errCommit
	}

	return web.UserLoginResponse{
		Token: tokenString,
	}, nil
}
