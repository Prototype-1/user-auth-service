package handlers

import (
	"context"
	"github.com/Prototype-1/user-auth-service/internal/usecase"
	proto "github.com/Prototype-1/user-auth-service/proto"
)

type UserHandler struct {
	proto.UnimplementedUserServiceServer
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) Signup(ctx context.Context, req *proto.SignupRequest) (*proto.AuthResponse, error) {
	_, err := h.userUsecase.Signup(req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &proto.AuthResponse{
		AccessToken:  "",  
		RefreshToken: "",
		Message:      "Signup successful",
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	token, err := h.userUsecase.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &proto.AuthResponse{
		AccessToken: token,
		RefreshToken: "",
		Message:      "Login successful",
	}, nil
}
