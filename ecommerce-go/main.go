package main

import (
	"github.com/Blaze5333/ecommerce_go/controllers"
	"github.com/Blaze5333/ecommerce_go/database"
	"github.com/Blaze5333/ecommerce_go/middleware"
	"github.com/Blaze5333/ecommerce_go/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())
	router.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantly", app.InstantBuy())
	log.Fatal(router.Run(":" + port))
}
