package main

import (
	"log"
	"net"

	// "net/http"

	productpb "go-micro-services/product-service/productpb"

	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

func main() {

	InitProductDB()

	listener, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	m := cmux.New(listener)

	httpListener := m.Match(cmux.HTTP1Fast())

	grpcListener := m.Match(cmux.HTTP2())

	router := gin.Default()
	router.POST("/products", CreateProduct)

	grpcServer := grpc.NewServer()
	productServiceServer := NewProductServiceServer(ProductCollection)
	productpb.RegisterProductServiceServer(grpcServer, productServiceServer)

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()

	go func() {
		if err := router.RunListener(httpListener); err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	log.Println("ProductService running on port 8002")
	if err := m.Serve(); err != nil {
		log.Fatalf("cmux server failed: %v", err)
	}
}
