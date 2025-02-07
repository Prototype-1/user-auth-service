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
	
	"google.golang.org/grpc"
)

func main() {
	utils.InitLogger() 
	utils.Log.Info("Logger initialized successfully")
	
	config.LoadConfig()

	utils.InitDB()

	userRepo := repository.NewUserRepository(utils.DB)

	userUsecase := usecase.NewUserUsecase(userRepo)

	userHandler := handlers.NewUserHandler(userUsecase)
	port := ":50052" 
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, userHandler)

	fmt.Printf("gRPC Server started on %s\n", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
