package controllers

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/smtp"
	"regexp"

	"github.com/Phanluc1/ecommerce-web/database"
	"github.com/Phanluc1/ecommerce-web/models"
	"github.com/gin-gonic/gin"
)

var codeMap = make(map[string]string)

func GetCodeChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var count int
		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err})
			return
		}

		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		re := regexp.MustCompile(emailRegex)
		if !re.MatchString(*user.Email) {
			c.JSON(422, gin.H{"error": "Invalid email format"})
			return
		}

		err := database.Client.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", user.Email).Scan(&count)
		if err != nil {
			c.JSON(500, gin.H{"message": "Error checking email existence", "error": err.Error()})
			return
		}
		if count == 0 {
			c.JSON(404, gin.H{"message": "Email does not exist "})
			return
		}
		err = sendEmail(*user.Email)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, gin.H{"message": "Check your email to get the code"})
	}
}
func UpdateNewPassWord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonData map[string]interface{}
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		email := jsonData["email"]
		newPasswod := jsonData["password"]
		code := jsonData["code"]

		emailStr := email.(string)
		newPasswodStr := newPasswod.(string)
		codeStr := code.(string)

		if codeMap[emailStr] != codeStr {
			c.JSON(400, gin.H{"message": "Code is incorrect"})
			return
		}
		newPasswodHashed := HashPassword(newPasswodStr)
		query := "UPDATE user SET password = ? WHERE email = ?"
		_, _ = database.Client.Query(query, newPasswodHashed, emailStr)
		delete(codeMap, emailStr)
		c.JSON(200, gin.H{"message": "Your password has been changed"})
	}
}
func generateOTP() string {
	const otpLength = 6
	otpChars := "0123456789"
	otp := make([]byte, otpLength)

	rand.Read(otp)

	for i := range otp {
		otp[i] = otpChars[int(otp[i])%len(otpChars)]
	}

	return string(otp)
}
func sendEmail(email string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	from := "lucphan1602@gmail.com"
	password := "ynex thsd sfwp gref"

	otp := generateOTP()
	codeMap[email] = otp

	subject := "Mã Xác Thực"
	body := fmt.Sprintf("Mã xác thực của bạn là: %s", otp)
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	_ = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(message))

	// if err != nil {
	// 	return err
	// }

	return nil
}
