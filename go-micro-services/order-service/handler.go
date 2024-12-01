package main

import (
	"context"
	"go-micro-services/common/models"
	"go-micro-services/common/utils"
	"log"
	"net/http"
	"time"

	productpb "go-micro-services/product-service/productpb"
	userpb "go-micro-services/user-service/userpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// "google.golang.org/grpc/credentials"
)

func checkUserExists(userID string) (bool, error) {

	conn, err := grpc.NewClient("127.0.0.1:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		return false, err
	}

	client := userpb.NewUserServiceClient(conn)
	reqCtx, reqCancel := context.WithTimeout(context.Background(), time.Second*2)
	defer reqCancel()

	res, err := client.CheckUserExists(reqCtx, &userpb.UserRequest{UserId: userID})
	if err != nil {
		log.Println(err)
		return false, err
	}

	return res.Exists, nil
}

func checkProductsExist(productIDs []string) (bool, error) {
	conn, err := grpc.NewClient("127.0.0.1:8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return false, err
	}
	// defer conn.Close()

	client := productpb.NewProductServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	res, err := client.CheckProductsExist(ctx, &productpb.ProductRequest{ProductIds: productIDs})
	if err != nil {
		return false, err
	}

	return res.AllExist, nil
}

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		utils.JSONResponse(c.Writer, http.StatusBadRequest, map[string]string{"error": "Invalid data"})
		return
	}

	userExists, err := checkUserExists(order.UserID)
	if err != nil || !userExists {
		utils.JSONResponse(c.Writer, http.StatusBadRequest, map[string]string{"error": "User does not exist"})
		return
	}

	productsExist, err := checkProductsExist(order.ProductIDs)
	if err != nil || !productsExist {
		utils.JSONResponse(c.Writer, http.StatusBadRequest, map[string]string{"error": "One or more products do not exist"})
		return
	}

	_, err = OrderCollection.InsertOne(context.TODO(), order)
	if err != nil {
		utils.JSONResponse(c.Writer, http.StatusInternalServerError, map[string]string{"error": "Failed to save order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}
