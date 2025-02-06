package usecase

import (
	"errors"
	"github.com/Prototype-1/user-auth-service/internal/models"
	"github.com/Prototype-1/user-auth-service/internal/repository"
	"github.com/Prototype-1/user-auth-service/config"
	"github.com/Prototype-1/user-auth-service/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Signup(name, email, password string) (*models.User, error)
	Login(email, password string) (string, error)
}

type userUsecaseImpl struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecaseImpl{userRepo: userRepo}
}

func (u *userUsecaseImpl) Signup(name, email, password string) (*models.User, error) {
	existingUser, _ := u.userRepo.GetUserByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	err = u.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecaseImpl) Login(email, password string) (string, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(int(user.ID), string(config.JWTSecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}
