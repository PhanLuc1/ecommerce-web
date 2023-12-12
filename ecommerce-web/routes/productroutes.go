package routes

import (
	"github.com/Phanluc1/ecommerce-web/controllers"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/product/:id", controllers.ProductDetails())
	incomingRoutes.GET("/product", controllers.ProductView())
	incomingRoutes.GET("/productGroup", controllers.GetGroup())
}
