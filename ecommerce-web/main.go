package main

import (
	"github.com/Phanluc1/ecommerce-web/middleware"
	"github.com/Phanluc1/ecommerce-web/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	routes.UserRoutes(router)
	routes.ProductRoutes(router)
	router.Run("0.0.0.0:8080")
}
