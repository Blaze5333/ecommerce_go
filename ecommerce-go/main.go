package main

import (
	"log"
	"os"

	"github.com/Blaze5333/ecommerce_go/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
		return
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	// app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	routes.AddressRoutes(router)
	routes.CartRoutes(router)
	routes.ProductRoutes(router)
	log.Fatal(router.Run(":" + port))
	log.Println("Server running on port:", port)
}
