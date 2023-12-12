package routes

import (
	"github.com/Phanluc1/ecommerce-web/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user/signup", controllers.Signup())
	incomingRoutes.POST("/user/login", controllers.Login())
	incomingRoutes.GET("/user", controllers.GetUserData())
	incomingRoutes.GET("/user/cart", controllers.GetUserCart())
	incomingRoutes.POST("/user/cart", controllers.ChangeProductCart)
	incomingRoutes.GET("/user/authenticaiton", controllers.Authenticate())
	incomingRoutes.POST("/user/code", controllers.GetCodeChangePassword())
	incomingRoutes.POST("/user/new-password", controllers.UpdateNewPassWord())
	incomingRoutes.GET("/user/wishlist", controllers.GetWishListCart())
	incomingRoutes.POST("/user/wishlist", controllers.ChangeCartWishList)
}
