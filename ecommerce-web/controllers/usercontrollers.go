package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Phanluc1/ecommerce-web/database"
	"github.com/Phanluc1/ecommerce-web/models"
	generate "github.com/Phanluc1/ecommerce-web/tokens"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var refreshTokenMap = make(map[string]bool)
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		re := regexp.MustCompile(emailRegex)
		if !re.MatchString(*user.Email) {
			c.JSON(422, gin.H{"error": "Invalid email format"})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(422, gin.H{"error": validationErr})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		result, err := createUser(ctx, database.Client, user)
		if err != nil {
			c.JSON(500, gin.H{"error": "Syntax SQL"})
			return
		}
		result.LastInsertId()
		c.JSON(201, gin.H{"message": "Your account has been created"})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := database.Client.QueryRow("SELECT * FROM user WHERE email = ?", user.Email).Scan(
			&founduser.Id,
			&founduser.Email,
			&founduser.First_Name,
			&founduser.Last_Name,
			&founduser.Password,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "Email is not available"})
			return
		}
		PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)
		if !PasswordIsValid {
			c.JSON(401, gin.H{"error": msg})
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, *founduser.Id)
		refreshTokenMap[refreshToken] = true
		c.IndentedJSON(http.StatusOK, gin.H{"token": token, "refreshToken": refreshToken})
	}
}

func GetUserData() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "" {
			claims, _ := generate.ValidateToken(token)
			if claims == nil {
				c.JSON(401, nil)
				return
			}
			c.IndentedJSON(http.StatusOK, gin.H{
				"id":        claims.User_ID,
				"firstName": claims.First_Name,
				"lastName":  claims.Last_Name,
				"email":     claims.Email,
			})
			return
		}
		c.JSON(401,nil)
	}
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken := c.GetHeader("refreshToken")
		if refreshToken != "" {
			if !refreshTokenMap[refreshToken] {
				c.JSON(401, nil)
				return
			}
			delete(refreshTokenMap, refreshToken)
			claims, _ := generate.ValidateToken(refreshToken)
			if claims == nil {
				c.JSON(http.StatusBadRequest, nil)
				return
			}
			var firstName string
			var lastName string
			var email string
			_ = database.Client.QueryRow("SELECT user.email,user.firstName,user.lastName FROM user WHERE id = ?", claims.User_ID).Scan(
				&email,
				&firstName,
				&lastName,
			)
			newToken, newRefreshToken, _ := generate.TokenGenerator(email, firstName, lastName, claims.User_ID)
			refreshTokenMap[newRefreshToken] = true
			c.IndentedJSON(http.StatusOK, gin.H{
				"token":        newToken,
				"refreshToken": newRefreshToken,
			})
			return
		}
		c.JSON(401,nil)
	}
}
func createUser(ctx context.Context, db *sql.DB, user models.User) (sql.Result, error) {
	stmt, err := db.PrepareContext(ctx, "INSERT INTO user (email, firstName, lastName, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, user.Email, user.First_Name, user.Last_Name, user.Password)
	if err != nil {
		return nil, err
	}
	return result, nil
}
