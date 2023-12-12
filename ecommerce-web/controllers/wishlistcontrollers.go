package controllers

import (
	"net/http"

	"github.com/Phanluc1/ecommerce-web/database"
	"github.com/Phanluc1/ecommerce-web/models"
	generate "github.com/Phanluc1/ecommerce-web/tokens"
	"github.com/gin-gonic/gin"
)

func GetWishListCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "" {
			var wishList []models.WishList
			claims, _ := generate.ValidateToken(token)
			if claims == nil {
				c.JSON(401, nil)
				return
			}
			query := "SELECT product.id,product.name,product.currentPrice,product.lastPrice FROM wishlistproduct JOIN product ON product.id = wishlistproduct.idProduct	WHERE wishlistproduct.idUser = ?"
			result, err := database.Client.Query(query, claims.User_ID)
			if err != nil {
				c.JSON(500, gin.H{"error": "Syntax SQL"})
				return
			}
			for result.Next() {
				var productWishList models.WishList
				var imageProduct []models.Image
				result.Scan(&productWishList.Id, &productWishList.Name, &productWishList.CurrentPrice, &productWishList.LastPrice)
				result1, err := database.Client.Query("SELECT productimage.url,productimage.description FROM productimage WHERE idProduct = ?", *productWishList.Id)
				if err != nil {
					c.JSON(500, gin.H{"error": err})
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
				productWishList.ProductImage = imageProduct
				wishList = append(wishList, productWishList)
			}
			c.IndentedJSON(http.StatusOK, wishList)
			return
		}
		c.JSON(401, nil)
	}
}

func AddCartIntoWishList() gin.HandlerFunc {
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
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			query := "INSERT INTO wishlistproduct (wishlistproduct.idUser,wishlistproduct.idProduct) VALUES (?, ?)"
			_, err := database.Client.Query(query, claims.User_ID, product.ProductId)
			if err != nil {
				c.JSON(200, gin.H{"Message": "The product is already in the wishlist"})
				return
			}
			var wishList []models.WishList
			query = "SELECT product.id,product.name,product.currentPrice,product.lastPrice FROM wishlistproduct JOIN product ON product.id = wishlistproduct.idProduct	WHERE wishlistproduct.idUser = ? AND wishlistproduct.idProduct = ?"
			result, err := database.Client.Query(query, claims.User_ID, product.ProductId)
			if err != nil {
				c.JSON(500, gin.H{"error": "Syntax SQL"})
				return
			}
			for result.Next() {
				var productWishList models.WishList
				var imageProduct []models.Image
				result.Scan(&productWishList.Id, &productWishList.Name, &productWishList.CurrentPrice, &productWishList.LastPrice)
				result1, err := database.Client.Query("SELECT productimage.url,productimage.description FROM productimage WHERE idProduct = ?", *productWishList.Id)
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
				productWishList.ProductImage = imageProduct
				wishList = append(wishList, productWishList)
			}
			c.IndentedJSON(http.StatusOK, wishList)
			return
		}
		c.JSON(201, nil)
	}
}
func DeleteAllCartWishList(deleteAll string) gin.HandlerFunc {
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
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			query := "DELETE FROM wishlistproduct WHERE wishlistproduct.idUser = ?"
			_, _ = database.Client.Query(query, claims.User_ID)
			c.JSON(http.StatusOK, gin.H{"message": "Done!"})
			return
		}
		c.JSON(401, nil)
	}
}

func DeleteProductIntoWishList(deleteStatus string) gin.HandlerFunc {
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
			query := "DELETE FROM wishlistproduct WHERE wishlistproduct.idUser = ? AND wishlistproduct.idProduct = ?"
			_, _ = database.Client.Query(query, claims.User_ID, product.ProductId)
			c.JSON(http.StatusOK, gin.H{"message": "Done!"})
			return
		}
		c.JSON(401, nil)
	}
}
func ChangeCartWishList(c *gin.Context) {
	deleteStatus := c.Query("delete")
	deleteAll := c.Query("deleteAll")
	if deleteStatus == "true" {
		DeleteProductIntoWishList(deleteStatus)(c)
		return
	}
	if deleteAll == "true" {
		DeleteAllCartWishList(deleteAll)(c)
		return
	}
	AddCartIntoWishList()(c)
}
