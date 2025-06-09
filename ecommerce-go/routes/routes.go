package routes

import (
	"github.com/Blaze5333/ecommerce_go/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/addproduct", controllers.AddProduct())
	incomingRoutes.GET("/users/productview", controllers.SearchProduct())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
}
