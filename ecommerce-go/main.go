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
	log.Println("Server running on port:", port)
}

// test
// func main() {
// 	fmt.Println("now time ", time.Now().Format("2006-01-02T15:04:05Z07:00"))
// 	t, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 	fmt.Println("second time", t)
// }
