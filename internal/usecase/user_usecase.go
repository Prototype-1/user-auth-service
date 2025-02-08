package usecase

import (
	"errors"
	"github.com/Prototype-1/user-auth-service/internal/models"
	"github.com/Prototype-1/user-auth-service/internal/repository"
	"github.com/Prototype-1/user-auth-service/config"
	"github.com/Prototype-1/user-auth-service/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	routePB "github.com/Prototype-1/user-auth-service/proto/routes"
	"context"
)

type UserUsecase interface {
	Signup(name, email, password string) (*models.User, error)
	Login(email, password string) (string, error)
	BlockUser(userID uint) error
    UnblockUser(userID uint) error
    SuspendUser(userID uint) error
    GetAllUsers() ([]*models.User, error)

	GetAllRoutes() ([]*models.Route, error)
}

type userUsecaseImpl struct {
	userRepo repository.UserRepository
	routeService routePB.RouteServiceClient
}

func NewUserUsecase(userRepo repository.UserRepository, routeClient routePB.RouteServiceClient) UserUsecase {
	return &userUsecaseImpl{
		userRepo: userRepo,
		routeService: routeClient,
	}
}

func (u *userUsecaseImpl) Signup(name, email, password string) (*models.User, error) {
	existingUser, err := u.userRepo.GetUserByEmail(email)
	if err != nil && err != gorm.ErrRecordNotFound { 
		return nil, errors.New("error checking existing user")
	}

	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	newUser := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	err = u.userRepo.CreateUser(newUser)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return newUser, nil
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
	token, err := utils.GenerateJWT(int(user.ID), user.Role, string(config.JWTSecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}


func (u *userUsecaseImpl) BlockUser(userID uint) error {
    user, err := u.userRepo.GetUserByID(userID)
    if err != nil {
        return err
    }
    user.BlockedStatus = true
    return u.userRepo.UpdateUser(user)
}

func (u *userUsecaseImpl) UnblockUser(userID uint) error {
    user, err := u.userRepo.GetUserByID(userID)
    if err != nil {
        return err
    }
    user.BlockedStatus = false
    return u.userRepo.UpdateUser(user)
}

func (u *userUsecaseImpl) SuspendUser(userID uint) error {
    user, err := u.userRepo.GetUserByID(userID)
    if err != nil {
        return err
    }
    user.InactiveStatus = true
    return u.userRepo.UpdateUser(user)
}

func (u *userUsecaseImpl) GetAllUsers() ([]*models.User, error) {
    return u.userRepo.GetAllUsers()
}

func (u *userUsecaseImpl) GetAllRoutes() ([]*models.Route, error) {
    res, err := u.routeService.GetAllRoutes(context.Background(), &routePB.GetAllRoutesRequest{})
    if err != nil {
        return nil, err
    }

    var routes []*models.Route
    for _, r := range res.Routes {
        routes = append(routes, &models.Route{
            RouteID:    int(r.RouteId),
            RouteName:  r.RouteName,
            StartStopID: int(r.StartStopId),
            EndStopID:   int(r.EndStopId),
            CategoryID: int(r.CategoryId),
        })
    }
    return routes, nil
}