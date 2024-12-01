package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	InitOrderDB()

	router := gin.Default()

	router.POST("/orders", CreateOrder)

	log.Println("OrderService running on port 8003")

	log.Fatal(router.Run(":8003"))
}
