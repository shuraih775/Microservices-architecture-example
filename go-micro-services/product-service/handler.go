package main

import (
	"context"
	"go-micro-services/common/models"
	"go-micro-services/common/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	productpb "go-micro-services/product-service/productpb"

	"github.com/gin-gonic/gin"
)

type ProductServiceServer struct {
	productpb.UnimplementedProductServiceServer
	db *mongo.Collection
}

func NewProductServiceServer(db *mongo.Collection) *ProductServiceServer {
	return &ProductServiceServer{db: db}
}

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.JSONResponse(c.Writer, http.StatusBadRequest, map[string]string{"error": "Invalid data"})
		return
	}

	_, err := ProductCollection.InsertOne(context.TODO(), product)
	if err != nil {
		utils.JSONResponse(c.Writer, http.StatusInternalServerError, map[string]string{"error": "Failed to save product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (s *ProductServiceServer) CheckProductsExist(ctx context.Context, req *productpb.ProductRequest) (*productpb.ProductResponse, error) {

	objectIDs := make([]primitive.ObjectID, 0, len(req.ProductIds))
	for _, id := range req.ProductIds {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {

			return nil, err
		}
		objectIDs = append(objectIDs, objID)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	count, err := s.db.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	allExist := count == int64(len(objectIDs))
	return &productpb.ProductResponse{AllExist: allExist}, nil
}
