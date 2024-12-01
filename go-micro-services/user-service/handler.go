package main

import (
	"context"
	"errors"
	"go-micro-services/common/models"
	"go-micro-services/common/utils"
	userpb "go-micro-services/user-service/userpb"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
	db *mongo.Collection
}

func NewUserServiceServer(db *mongo.Collection) *UserServiceServer {
	return &UserServiceServer{db: db}
}

func (s *UserServiceServer) CheckUserExists(ctx context.Context, req *userpb.UserRequest) (*userpb.UserResponse, error) {

	userId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {

		return nil, errors.New("invalid user ID format")
	}

	filter := bson.M{"_id": userId}
	var result struct{}
	err = s.db.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(err, filter)

			return &userpb.UserResponse{Exists: false}, nil
		}
		return nil, err
	}

	return &userpb.UserResponse{Exists: true}, nil
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.JSONResponse(c.Writer, http.StatusBadRequest, map[string]string{"error": "Invalid data"})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.JSONResponse(c.Writer, http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	_, err = UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		utils.JSONResponse(c.Writer, http.StatusInternalServerError, map[string]string{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func LoginUser(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user struct {
		ID       string `json:"_id" bson:"_id"`
		Email    string `json:"email" bson:"email"`
		Password string `json:"password" bson:"password"`
	}
	err := UserCollection.FindOne(context.TODO(), bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
