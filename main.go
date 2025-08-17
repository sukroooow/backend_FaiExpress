package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mubarok-ridho/misi-paket.backend/config"
	"github.com/mubarok-ridho/misi-paket.backend/route"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // hanya kalau butuh cookie/session
	}))

	route.SetupRoutes(r)

	r.Run(":8080")
}
