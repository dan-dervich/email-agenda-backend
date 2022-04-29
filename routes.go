package main

import (
	"fmt"
	"time"
	
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json: "username"`
	Email string `json:"email"`
	Password string `json:"password"`
}
func LoginHandler(c *gin.Context) {
	Users, err := db()
	if err != nil {
		fmt.Println(err)
		return
	}
	email := c.PostForm("email")
	password := c.PostForm("password")
	userData, err := CheckUserExistenceByQuery(bson.M{"email": email}, Users)
	if err != nil {
		c.JSON(200, gin.H{
			"status": "userNotFound",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(password))
	if err != nil {
		c.JSON(200, gin.H{
			"status": "errorHashingPassword",
		})
		return
	}
	token, err := CreateJWT(userData.Email, 300)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"status": "errorCreatingJWTtoken",
		})
		return
	}
	expirationTime := time.Now().Add(300 * time.Minute).Unix()
	c.SetCookie("email", token, int(expirationTime), "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"status": "everythingIsFine",
	})
}

func SignUpHandler(c *gin.Context) {
	Users, err := db()
	if err != nil {
		fmt.Println(err)
		return
	}
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	username := req.Username
	password := req.Password
	email := req.Email

	checkUser, err := CheckUserExistenceByQuery(bson.M{"email": email}, Users)

	if err != nil && fmt.Sprint(err) != "mongo: no documents in result" || len(checkUser.Password) > 0 {
		c.JSON(200, gin.H{
			"status": "userWithThatUsernameAlreadyExists",
		})
		return
	} else {
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"status": "errorHashingPassword",
			})
			return
		}
		newUserData := User{
			Username: username,
			Password: string(encryptedPassword),
			Email:    email,
		}
		newUser := CreateUser(Users, newUserData)
		token, err := CreateJWT(email, 300)
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"status": "errorCreatingJWTtoken",
			})
			return
		}

		// expirationTime := time.Now().Add(300 * time.Minute).Unix()
		// c.SetCookie("email", token, int(expirationTime), "/", "localhost", false, true)

		c.JSON(200, gin.H{
			"ID":     newUser,
			"jwtToken": token,
			"status": "userCreatedSuccessfully",
		})
		return
	}
}
