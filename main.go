package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(cors.Default())
	
	// router.SetTrustedProxies([]string{"192.168.1.2", "http://localhost:3000", "localhost:3000", "localhost:8080"})
	// router.LoadHTMLGlob("templates/*")
	router1 := router.Group("/auth")
	{
		router1.POST("/login", LoginHandler)
		router1.POST("/sign-up", SignUpHandler)
	}
	// router.RunTLS(":3000", "./test/cert.pem", "./test/snakeoil.key") // this can be used for https try another day to make that work
	router.Run(":8080")
}
