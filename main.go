package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(cors.Default())
	router1 := router.Group("/auth")
	{
		go router1.POST("/login", LoginHandler)
		go router1.POST("/sign-up", SignUpHandler)
	}
	router.Run(":8080")
}
