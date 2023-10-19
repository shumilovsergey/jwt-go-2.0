package controllers

import (
	"jwt/initializers"
	"jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// что такуое * ? - operator
// что такое & ? - operator
// что такое с ?  -  code
// что такое gin.H{} - message

func Signup(c *gin.Context) {
	// Get email + pass
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// проверка на пустые поля

	if body.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Miss email",
		})
		return
	}

	if body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Miss password",
		})
		return
	}

	// Hash password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user (возможно что почта уже занята но это не точно)",
		})
		return
	}
	// Respond 200

	c.JSON(http.StatusOK, gin.H{})

}

func Login(c *gin.Context) {
	// Get correct body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	// Find user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email",
		})
		return
	}

	// Compare passwords hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	// Generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token generate ERROR",
		})
		return
	}

	// Send it back
	c.SetSameSite(http.SameSiteDefaultMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)

	c.JSON(http.StatusBadRequest, gin.H{})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	// id :=user.(models.User).ID

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
