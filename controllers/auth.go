package controllers

import (
	"bognar.dev-backend/models"
	token "bognar.dev-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("JSON Bind error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password

	t, err := models.LoginCheck(&u)

	if err != nil {
		fmt.Println("LoginCheck error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "username or password is incorrect."})
		return
	}
	fmt.Println("LoginCheck success got token: ", t)

	c.JSON(http.StatusOK, gin.H{"message": "success", "token": t})

}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password

	_, err := u.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})

}

func Status(c *gin.Context) {
	fmt.Println("{status: You are logged in}")
	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
}

func CurrentUser(c *gin.Context) {

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}
