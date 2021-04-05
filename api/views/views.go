package views

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/splinter0/api/database"
	"github.com/splinter0/api/security"
)

var validate = validator.New()

func bad(c *gin.Context) {
	c.JSON(400, gin.H{
		"message": "Bad request",
	})
}

type login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"`
}

func Login(c *gin.Context) {
	var l login
	err := c.ShouldBindJSON(&l)
	if err != nil {
		bad(c)
		return
	}
	validationErr := validate.Struct(l)
	if validationErr != nil {
		bad(c)
		return
	}
	user := database.FindUser(l.Username)
	fmt.Println(user.Password)
	if user.Username == l.Username && security.VerifyPassword(user.Password, l.Password) {
		token := security.GenerateToken(user.Username, user.Role)
		database.AddUserToken(user.Username, token)
		c.JSON(200, gin.H{
			"message":  "success",
			"token":    token,
			"username": user.Username,
		})
	} else {
		security.NotAuth(c)
	}
}

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "success",
		"role":    c.GetString("role"),
		"user":    c.GetString("username"),
	})
}
