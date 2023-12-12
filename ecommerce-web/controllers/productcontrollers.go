package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Phanluc1/ecommerce-web/database"
	"github.com/Phanluc1/ecommerce-web/models"
	"github.com/gin-gonic/gin"
)

func ProductDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		var product models.Product
		var idCategory int
		var imageProduct []models.Image
		var productGroups []models.Group
		var productTechnicals []models.Technical
		err := database.Client.QueryRow("SELECT * FROM product WHERE id = ?", productID).Scan(
			&product.ProductId,
			&product.ProductName,
			&product.Rate,
			&product.Sold,
			&product.CurrentPrice,
			&product.LastPrice,
			&idCategory,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product is not available"})
			return
		}

		result1, err := database.Client.Query("SELECT productimage.url,productimage.description FROM productimage WHERE idProduct = ?", *product.ProductId)
		if err != nil {
			c.JSON(500, gin.H{"error 2": err})
			return
		}
		for result1.Next() {
			var imageTemp models.Image
			err := result1.Scan(&imageTemp.Url, &imageTemp.Description)
			if err != nil {
				c.JSON(404, gin.H{"error": "No matching information found"})
				return
			}
			imageProduct = append(imageProduct, imageTemp)
		}

		result5, err := database.Client.Query("SELECT producttechnical.title,producttechnical.description FROM producttechnical WHERE idProduct = ?", *product.ProductId)
		if err != nil {
			c.JSON(500, gin.H{"error": "Syntax SQL"})
		}
		for result5.Next() {
			var producttechnical models.Technical
			err := result5.Scan(&producttechnical.Title, &producttechnical.Description)
			if err != nil {
				c.JSON(404, gin.H{"error": "No matching information found"})
				return
			}
			productTechnicals = append(productTechnicals, producttechnical)
		}

		query := "SELECT DISTINCT productgroup.id,productgroup.tittle FROM productgroup	JOIN producttype ON producttype.idGroup = productgroup.id JOIN product_producttype ON product_producttype.idType = producttype.id JOIN product ON product_producttype.idProduct = product.id WHERE product.id = " + productID
		query += " ORDER BY productgroup.id ASC "
		result2, err := database.Client.Query(query)
		if err != nil {
			c.JSON(500, gin.H{"error": "SQL syntax"})
			return
		}
		for result2.Next() {
			var productGroup models.Group
			var types []models.Type
			result2.Scan(&productGroup.Id, &productGroup.Title)
			result3, err := database.Client.Query("SELECT producttype.id,producttype.title,producttype.description FROM producttype JOIN product_producttype ON product_producttype.idType = producttype.id JOIN product ON product.id = product_producttype.idProduct WHERE idGroup = ? AND idProduct = ?", productGroup.Id, productID)
			if err != nil {
				c.JSON(500, gin.H{"error": "syntax SQL"})
				return
			}
			for result3.Next() {
				var typeTemp models.Type
				err := result3.Scan(&typeTemp.Id, &typeTemp.Title, &typeTemp.Description)
				if err != nil {
					c.JSON(200, gin.H{"Erorr": err})
					return
				}
				types = append(types, typeTemp)
			}

			productGroup.Type = types
			productGroups = append(productGroups, productGroup)
		}

		var productComments []models.Comment
		query = "SELECT user.firstName,user.lastName,productcomment.description FROM user JOIN productcomment ON productcomment.idUser = user.id JOIN product ON product.id = productcomment.idProduct WHERE product.id = ?"
		result4, err := database.Client.Query(query, productID)
		if err != nil {
			c.JSON(500, gin.H{"error": "syntax SQL"})
		}
		for result4.Next() {
			var productComment models.Comment
			err := result4.Scan(&productComment.FirstName, &productComment.LastName, &productComment.Description)
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
			}
			productComments = append(productComments, productComment)
		}

		product.ProductImages = imageProduct
		product.ProuctTechnical = productTechnicals
		product.ProductGroup = productGroups
		product.ProductComment = productComments
		c.IndentedJSON(http.StatusOK, product)
	}
}
func SearchProduct() gin.HandlerFunc {
	return nil
}
func GetGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productGroups []models.Group
		categoryId := c.Query("categoryId")
		query := "SELECT DISTINCT productgroup.id,productgroup.tittle FROM productgroup JOIN producttype ON producttype.idGroup = productgroup.id JOIN product_producttype ON product_producttype.idType = producttype.id JOIN product ON product_producttype.idProduct = product.id JOIN category ON category.id = product.idCategory WHERE category.id = " + categoryId
		query += " ORDER BY productgroup.id ASC "
		result, err := database.Client.Query(query)
		if err != nil {
			c.JSON(200, gin.H{"message": "Không tìm thấy sản phẩm phù hợp"})
			return
		}
		for result.Next() {
			var productGroup models.Group
			var types []models.Type
			result.Scan(&productGroup.Id, &productGroup.Title)
			result2, err := database.Client.Query("SELECT producttype.id,producttype.title,producttype.description FROM producttype WHERE idGroup = ?", productGroup.Id)
			if err != nil {
				c.JSON(200, gin.H{"message": "Không tìm thấy sản phẩm phù hợp"})
				return
			}
			for result2.Next() {
				var typeTemp models.Type
				err := result2.Scan(&typeTemp.Id, &typeTemp.Title, &typeTemp.Description)
				if err != nil {
					c.JSON(200, gin.H{"Erorr": err})
					return
				}
				types = append(types, typeTemp)
			}
			productGroup.Type = types
			productGroups = append(productGroups, productGroup)
		}
		c.IndentedJSON(http.StatusOK, productGroups)
	}
}
func SearchProductByQuery() gin.HandlerFunc {
	return nil
}
func ProductView() gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []models.Product
		categoryID := c.Query("categoryId")
		productTypeId := c.Query("productTypeId")
		sort := c.Query("sort")
		var params []interface{}
		query := "SELECT p.* FROM product AS p JOIN product_producttype AS ppt ON p.id = ppt.idproduct JOIN producttype AS pt ON ppt.idtype = pt.id JOIN productgroup AS pg ON pt.idGroup = pg.id WHERE "
		if categoryID != "" {
			query += fmt.Sprintf("idCategory = %s ", categoryID)
		}
		if productTypeId != "" {
			typeIds := strings.Split(productTypeId, "-")
			if categoryID != "" {
				query += "AND pt.id IN (?"
			} else {
				query += "pt.id IN (?"
			}
			for i := 1; i < len(typeIds); i++ {
				query += ", ?"
			}
			query += ") "
			for _, typeId := range typeIds {
				params = append(params, typeId)
			}
		}
		if sort == "ascending" {
			query += "ORDER BY p.currentPrice ASC"
		} else if sort == "descending" {
			query += "ORDER BY p.currentPrice DESC"
		}

		if categoryID == "" && productTypeId == "" {
			query = "SELECT * FROM product"
		}
		result, err := database.Client.Query(query, params...)
		if err != nil {
			c.JSON(200, gin.H{"message": "Không tìm thấy sản phẩm phù hợp"})
			return
		}
		for result.Next() {
			var product models.Product
			var imageProduct []models.Image
			var idCategory int
			err := result.Scan(&product.ProductId, &product.ProductName, &product.Rate, &product.Sold, &product.CurrentPrice, &product.LastPrice, &idCategory)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			result1, err := database.Client.Query("SELECT productimage.url,productimage.description FROM productimage WHERE idProduct = ?", *product.ProductId)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			for result1.Next() {
				var imageTemp models.Image
				err := result1.Scan(&imageTemp.Url, &imageTemp.Description)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err})
					return
				}
				imageProduct = append(imageProduct, imageTemp)
			}
			product.ProductImages = imageProduct
			products = append(products, product)
		}
		c.IndentedJSON(http.StatusOK, products)
	}
}
