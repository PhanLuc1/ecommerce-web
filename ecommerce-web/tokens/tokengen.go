package tokens

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/Phanluc1/ecommerce-web/database"
	jwt "github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	User_ID    int
	Email      string
	First_Name string
	Last_Name  string
	jwt.StandardClaims
}

var UserData *sql.Rows = database.UserData(database.Client, "USERS")
var SECRET_KEY = os.Getenv("SECRET_LOVE")

func TokenGenerator(email string, firstname string, lastname string, userid int) (signedtoken string, signedrefreshtoken string, err error) {
	claims := &SignedDetails{
		User_ID:    userid,
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}
	refreshclaims := &SignedDetails{
		User_ID: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panicln(err)
		return
	}
	return token, refreshtoken, err
}
func ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The Token is invalid"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}
	return claims, msg
}
