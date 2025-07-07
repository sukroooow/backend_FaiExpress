package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mubarok-ridho/misi-paket.backend/controller"
	handlers "github.com/mubarok-ridho/misi-paket.backend/handler"
	"github.com/mubarok-ridho/misi-paket.backend/middleware"
)

func RegisterChatRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(middleware.JWTAuthMiddleware())

	api.POST("/chat", middleware.RoleMiddleware("customer", "kurir"), controller.SendChat)
	api.GET("/chat", middleware.RoleMiddleware("customer", "kurir"), controller.GetChat)
}

func RegisterWebSocketRoutes(r *gin.Engine) {
	r.GET("/ws/chat/:orderId", handlers.ChatWebSocket)
}
