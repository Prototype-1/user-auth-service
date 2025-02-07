package handlers

import (
	"context"
	"github.com/Prototype-1/user-auth-service/internal/usecase"
	proto "github.com/Prototype-1/user-auth-service/proto"
)

type UserHandler struct {
	proto.UnimplementedUserServiceServer
	userUsecase usecase.UserUsecase
	routeUsecase usecase.RouteUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase, routeUsecase usecase.RouteUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		routeUsecase: routeUsecase
	}
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

func (h *UserHandler) BlockUser(ctx context.Context, req *proto.UserRequest) (*proto.StatusResponse, error) {
	err := h.userUsecase.BlockUser(uint(req.UserId))
	if err != nil {
		return nil, err
	}

	return &proto.StatusResponse{
		Message: "User blocked successfully",
	}, nil
}

func (h *UserHandler) UnblockUser(ctx context.Context, req *proto.UserRequest) (*proto.StatusResponse, error) {
	err := h.userUsecase.UnblockUser(uint(req.UserId))
	if err != nil {
		return nil, err
	}

	return &proto.StatusResponse{
		Message: "User unblocked successfully",
	}, nil
}

func (h *UserHandler) SuspendUser(ctx context.Context, req *proto.UserRequest) (*proto.StatusResponse, error) {
	err := h.userUsecase.SuspendUser(uint(req.UserId))
	if err != nil {
		return nil, err
	}

	return &proto.StatusResponse{
		Message: "User suspended successfully",
	}, nil
}

func (h *UserHandler) GetAllUsers(ctx context.Context, req *proto.Empty) (*proto.UserList, error) {
	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var userList []*proto.User
	for _, u := range users {
		userList = append(userList, &proto.User{
			Id:             uint32(u.ID),
			Email:          u.Email,
			Name:           u.Name,
			BlockedStatus:  u.BlockedStatus,
			InactiveStatus: u.InactiveStatus,
		})
	}

	return &proto.UserList{
		Users: userList,
	}, nil
}

func (h *UserHandler) GetAllRoutes(ctx context.Context, req *proto.Empty) (*proto.RouteList, error) {
	routes, err := h.routeUsecase.GetAllRoutes()
	if err != nil {
		return nil, err
	}

	var routeList []*proto.Route
	for _, r := range routes {
		routeList = append(routeList, &proto.Route{
			Id:          uint32(r.ID),
			Name:        r.Name,
			Source:      r.Source,
			Destination: r.Destination,
		})
	}

	return &proto.RouteList{Routes: routeList}, nil
}