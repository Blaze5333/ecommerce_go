package routes

import (
	"github.com/Blaze5333/ecommerce_go/controllers"
	"github.com/Blaze5333/ecommerce_go/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/addproduct", controllers.AddProduct())
	incomingRoutes.GET("/users/productview", controllers.SearchProduct())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
}
func AddressRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/address", middleware.Authentication(), controllers.AddAddress())
	incomingRoutes.GET("/users/address", middleware.Authentication(), controllers.GetAddress())
	incomingRoutes.DELETE("/users/address", middleware.Authentication(), controllers.DeleteAddress())
	incomingRoutes.DELETE("/users/address/:address_id", middleware.Authentication(), controllers.DeleteAddressById())
	incomingRoutes.PUT("/users/address/:address_id", middleware.Authentication(), controllers.EditWorkAddress())
}
func CartRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users/cart/", middleware.Authentication(), controllers.GetItemFromCart())
	incomingRoutes.POST("/users/cart/:product_id", middleware.Authentication(), controllers.AddToCart())
	incomingRoutes.DELETE("/users/cart/:product_id", middleware.Authentication(), controllers.RemoveCartItem())
	incomingRoutes.POST("/users/cart/buy", middleware.Authentication(), controllers.BuyFromCart())
	incomingRoutes.POST("/users/cart/instantbuy/:product_id", middleware.Authentication(), controllers.InstantBuy())
}
func ProductRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/product", middleware.Authentication(), controllers.AddProductInDB())
	incomingRoutes.GET("/product", middleware.Authentication(), controllers.SearchProduct())
	incomingRoutes.GET("/productByQuery", middleware.Authentication(), controllers.SearchProductByQuery())
}
