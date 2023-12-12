package controllers

import (
	"net/http"

	"github.com/Phanluc1/ecommerce-web/database"
	"github.com/Phanluc1/ecommerce-web/models"
	generate "github.com/Phanluc1/ecommerce-web/tokens"
	"github.com/gin-gonic/gin"
)

func GetUserCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "" {
			var cartItems []models.CartItem
			claims, _ := generate.ValidateToken(token)
			if claims == nil {
				c.JSON(401, nil)
				return
			}
			query := "SELECT product.id,product.name,product.currentPrice,product.lastPrice,user_product.quantity FROM user_product JOIN product ON product.id = user_product.idProduct	WHERE user_product.idUser = ?"
			result, err := database.Client.Query(query, claims.User_ID)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"message": "Shopping cart is empty"})
				return
			}
			for result.Next() {
				var cartItem models.CartItem
				var imageProduct []models.Image
				result.Scan(&cartItem.Id, &cartItem.Name, &cartItem.CurrentPrice, &cartItem.LastPrice, &cartItem.Quantity)
				result1, err := database.Client.Query("SELECT productimage.url,productimage.description FROM productimage WHERE idProduct = ?", *cartItem.Id)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				for result1.Next() {
					var imageTemp models.Image
					err := result1.Scan(&imageTemp.Url, &imageTemp.Description)
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return
					}
					imageProduct = append(imageProduct, imageTemp)
				}
				cartItem.ProductImage = imageProduct
				cartItems = append(cartItems, cartItem)
			}
			c.IndentedJSON(http.StatusOK, cartItems)
			return
		}
		c.JSON(401, nil)
	}
}

func AddProductIntoCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "" {
			claims, _ := generate.ValidateToken(token)
			if claims == nil {
				c.JSON(401, nil)
				return
			}
			var product models.Product
			if err := c.BindJSON(&product); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query := "INSERT INTO user_product (user_product.idUser, user_product.idProduct, quantity) VALUES ( ?, ?, 1) ON DUPLICATE KEY UPDATE quantity = quantity + 1"
			_, err := database.Client.Query(query, claims.User_ID, product.ProductId)
			if err != nil {
				c.JSON(500, gin.H{"error": "Syntax SQL"})
			}
			var cartItems []models.CartItem
			query = "SELECT product.id,product.name,product.currentPrice,product.lastPrice,user_product.quantity FROM user_product JOIN product ON product.id = user_product.idProduct	WHERE user_product.idUser = ? AND user_product.idProduct = ?"
			result, err := database.Client.Query(query, claims.User_ID, product.ProductId)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"message": "Shopping cart is empty"})
				return
			}
			for result.Next() {
				var cartItem models.CartItem
				var imageProduct []models.Image
				result.Scan(&cartItem.Id, &cartItem.Name, &cartItem.CurrentPrice, &cartItem.LastPrice, &cartItem.Quantity)
				result1, err := database.Client.Query("SELECT productimage.url,productimage.description FROM productimage WHERE idProduct = ?", *cartItem.Id)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				for result1.Next() {
					var imageTemp models.Image
					err := result1.Scan(&imageTemp.Url, &imageTemp.Description)
					if err != nil {
						c.JSON(500, gin.H{"error": err})
						return
					}
					imageProduct = append(imageProduct, imageTemp)
				}
				cartItem.ProductImage = imageProduct
				cartItems = append(cartItems, cartItem)
			}
			c.IndentedJSON(http.StatusOK, cartItems)
			return
		}
		c.JSON(401, nil)
	}
}
func DeleteProductCart(deleteStatus string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "" {
			claims, _ := generate.ValidateToken(token)
			if claims == nil {
				c.JSON(401, nil)
				return
			}
			var product models.Product
			if err := c.BindJSON(&product); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query := "DELETE FROM user_product WHERE user_product.idUser = ? AND user_product.idProduct = ?"
			_, _ = database.Client.Query(query, claims.User_ID, product.ProductId)
			c.JSON(http.StatusOK, gin.H{"message": "Done!"})
			return
		}
		c.JSON(401, nil)
	}
}
func DecreaseProductCart(decrease string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "" {
			claims, _ := generate.ValidateToken(token)
			if claims == nil {
				c.JSON(401, nil)
				return
			}
			var product models.Product
			if err := c.BindJSON(&product); err != nil {
				c.JSON(400, gin.H{"error": err})
				return
			}
			query := "INSERT INTO user_product (user_product.idUser, user_product.idProduct, quantity) VALUES ( ?, ?, 1) ON DUPLICATE KEY UPDATE quantity = quantity - 1"
			_, err := database.Client.Query(query, claims.User_ID, product.ProductId)
			if err != nil {
				c.JSON(500, gin.H{"error": "Syntax SQL"})
			}
			var cartItems []models.CartItem
			query = "SELECT product.id,product.name,product.currentPrice,product.lastPrice,user_product.quantity FROM user_product JOIN product ON product.id = user_product.idProduct	WHERE user_product.idUser = ? AND user_product.idProduct = ?"
			result, _ := database.Client.Query(query, claims.User_ID, product.ProductId)
			for result.Next() {
				var cartItem models.CartItem
				var imageProduct []models.Image
				result.Scan(&cartItem.Id, &cartItem.Name, &cartItem.CurrentPrice, &cartItem.LastPrice, &cartItem.Quantity)
				result1, err := database.Client.Query("SELECT productimage.url,productimage.description FROM productimage WHERE idProduct = ?", *cartItem.Id)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				for result1.Next() {
					var imageTemp models.Image
					err := result1.Scan(&imageTemp.Url, &imageTemp.Description)
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return
					}
					imageProduct = append(imageProduct, imageTemp)
				}
				cartItem.ProductImage = imageProduct
				cartItems = append(cartItems, cartItem)
			}
			c.IndentedJSON(http.StatusOK, cartItems)
			return
		}
		c.JSON(401, nil)
	}
}
func DeleteAllProductCart(deleteAll string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "" {
			claims, _ := generate.ValidateToken(token)
			if claims == nil {
				c.JSON(401, nil)
				return
			}
			var product models.Product
			if err := c.BindJSON(&product); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			query := "DELETE FROM user_product WHERE user_product.idUser = ?"
			_, _ = database.Client.Query(query, claims.User_ID)
			c.JSON(http.StatusOK, gin.H{"message": "Done!"})
			return
		}
		c.JSON(401, nil)
	}
}
func ChangeProductCart(c *gin.Context) {
	deleteStatus := c.Query("delete")
	decrease := c.Query("decrease")
	deleteAll := c.Query("deleteAll")
	if deleteStatus == "true" {
		DeleteProductCart(deleteStatus)(c)
		return
	}
	if decrease == "true" {
		DecreaseProductCart(decrease)(c)
		return
	}
	if deleteAll == "true" {
		DeleteAllProductCart(deleteAll)(c)
		return
	}
	AddProductIntoCart()(c)
}
