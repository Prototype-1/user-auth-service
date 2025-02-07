package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Prototype-1/user-auth-service/config"
	"github.com/Prototype-1/user-auth-service/internal/repository"
	"github.com/Prototype-1/user-auth-service/internal/usecase"
	"github.com/Prototype-1/user-auth-service/internal/handlers"
	proto "github.com/Prototype-1/user-auth-service/proto"
	"github.com/Prototype-1/user-auth-service/utils"
	routePB "github.com/Prototype-1/user-auth-service/proto/routes"
	"google.golang.org/grpc"
)

func main() {
	utils.InitLogger() 
	utils.Log.Info("Logger initialized successfully")

	config.LoadConfig()

	utils.InitDB()

	routeServiceConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure()) 
	if err != nil {
		log.Fatalf("Failed to connect to route service: %v", err)
	}
	defer routeServiceConn.Close()

	routeClient := routePB.NewRouteServiceClient(routeServiceConn)

	userRepo := repository.NewUserRepository(utils.DB)
	userUsecase := usecase.NewUserUsecase(userRepo, routeClient)
	userHandler := handlers.NewUserHandler(userUsecase, routeClient)


	port := ":50052" 
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, userHandler)

	fmt.Printf("gRPC Server started on %s...\n", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
