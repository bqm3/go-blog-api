package main

import (
	"blog-api/config"
	"blog-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.ConnectDatabase()
	r.StaticFS("/uploads", gin.Dir("public/uploads", true))

	routes.SetupRoutes(r)
	r.Run(":8080")
}
