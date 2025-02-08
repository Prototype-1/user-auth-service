package handlers

import (
	"context"
	"github.com/Prototype-1/user-auth-service/internal/usecase"
	proto "github.com/Prototype-1/user-auth-service/proto"
	routePB "github.com/Prototype-1/user-auth-service/proto/routes"
	"log"
	"github.com/Prototype-1/user-auth-service/utils"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type UserHandler struct {
	proto.UnimplementedUserServiceServer
	userUsecase usecase.UserUsecase
	routeClient routePB.RouteServiceClient
}

func NewUserHandler(userUsecase usecase.UserUsecase, routeClient routePB.RouteServiceClient) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		routeClient: routeClient,
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

func (s *UserHandler) authenticateUser(ctx context.Context) (uint, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, "", status.Errorf(codes.Unauthenticated, "Missing metadata")
	}

	tokenList, exists := md["authorization"]
	if !exists || len(tokenList) == 0 {
		return 0, "", status.Errorf(codes.Unauthenticated, "Authorization token not provided")
	}

	tokenString := tokenList[0]
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	userID, role, err := utils.ParseJWT(tokenString)
	if err != nil {
		return 0, "", status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}

	return userID, role, nil
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

func (h *UserHandler) GetAllRoutes(ctx context.Context, req *routePB.GetAllRoutesRequest) (*routePB.GetAllRoutesResponse, error) {
	adminID, role, err := h.authenticateUser(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("User ID %d (Role: %s) is retrieving all routes", adminID, role)

	response, err := h.routeClient.GetAllRoutes(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch routes: %v", err)
	}

	return response, nil
}
