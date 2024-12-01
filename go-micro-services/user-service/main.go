package main

import (
	"log"
	"net"

	userpb "go-micro-services/user-service/userpb"

	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

func main() {
	InitUserDB()

	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	m := cmux.New(listener)

	httpListener := m.Match(cmux.HTTP1Fast())

	grpcListener := m.Match(cmux.HTTP2())

	router := gin.Default()

	router.POST("/login", LoginUser)
	router.POST("/users", CreateUser)

	go func() {
		if err := router.RunListener(httpListener); err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	grpcServer := grpc.NewServer()
	userServiceServer := NewUserServiceServer(UserCollection)
	userpb.RegisterUserServiceServer(grpcServer, userServiceServer)

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()

	log.Println("UserService running on port 8001")

	if err := m.Serve(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
