package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mubarok-ridho/misi-paket.backend/controller"
	handlers "github.com/mubarok-ridho/misi-paket.backend/handler"
	"github.com/mubarok-ridho/misi-paket.backend/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// ✅ Root
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "FaiExpress API is running!"})
	})

	// ✅ WebSocket Chat (per Order ID)
	r.POST("/chat/send", handlers.SendChatMessage)
	r.GET("/centrifugo/token", handlers.GenerateCentrifugoToken)
	r.GET("/chat/load/:order_id", handlers.GetMessagesByOrderID)
	r.POST("/schedule/delete-chat/:id", handlers.ScheduleDeleteChat)

	// ✅ Auth
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.PUT("/users/:id/password", controller.ChangePassword)

	// ✅ Tracking
	r.POST("/kurir/track", controller.UpdateKurirLocation)
	r.GET("/kurir/track/:id", controller.GetKurirLocation)
	r.GET("/kurir/:id/location", controller.GetKurirLocation)
	r.GET("/kurir/available", controller.GetAvailableKurir)
	r.PUT("/api/orders/tagihan", middleware.AuthMiddleware(), controller.UpdateTagihan)
	r.PUT("/api/orders/payment-validasi", middleware.AuthMiddleware(), controller.ValidasiPembayaran)
	r.PUT("/api/orders/:id/metode_bayar", middleware.AuthMiddleware(), controller.UpdatePaymentMethod)
	r.GET("/pendapatan/total-today", controller.GetTotalPendapatanToday)

	// ✅ Protected with JWT
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuthMiddleware())

	// Kurir
	auth.GET("/kurir/:id/orders", middleware.RoleMiddleware("kurir"), controller.GetOrdersForKurir)
	auth.PUT("/kurir/status", middleware.RoleMiddleware("admin"), controller.UpdateKurirStatus)
	auth.PUT("/kurir/location", middleware.RoleMiddleware("kurir"), controller.UpdateLocation)
	auth.GET("/kurir/:id", middleware.RoleMiddleware("kurir"), controller.GetKurirByID)
	auth.PUT("/kurir/up/:id", middleware.RoleMiddleware("kurir"), controller.UpdateKurirByID)
	auth.GET("/kurir/:id/orders/proses", controller.GetOrdersProses)
	auth.GET("/kurir/:id/orders/selesai/today", controller.GetOrdersSelesaiToday)
	auth.GET("/pendapatan/kurir/:id/today", controller.GetPendapatanKurirToday)

	// Customer - Orders
	auth.POST("/orders", middleware.RoleMiddleware("customer"), controller.CreateOrder)
	auth.GET("/my-orders", middleware.RoleMiddleware("customer"), controller.GetMyOrders)
	auth.PUT("/update-profile", controller.UpdateProfile) // ⬅️ Ini baru

	// Admin - Orders
	auth.GET("/orders", middleware.RoleMiddleware("admin"), controller.GetAllOrders)
	auth.GET("/orders/:id", controller.GetOrderByID)
	auth.PUT("/orders/:id", controller.UpdateOrder)
	auth.DELETE("/orders/:id", middleware.RoleMiddleware("admin"), controller.DeleteOrder)
	auth.PUT("/orders/status", middleware.RoleMiddleware("kurir"), controller.UpdateOrderStatus)
	auth.GET("/orders/total-selesai-today", controller.GetTotalOrdersSelesaiToday)
	auth.GET("/pendapatan/total-all-today", controller.GetAllTotalPendapatanToday)

	// Chat via REST API (opsional)
	auth.POST("/chat", middleware.RoleMiddleware("customer", "kurir"), controller.SendChat)
	auth.GET("/chat", middleware.RoleMiddleware("customer", "kurir"), controller.GetChat)

	// Admin - User CRUD
	auth.GET("/users", middleware.RoleMiddleware("admin"), controller.GetAllUsers)
	auth.GET("/users/:id", middleware.RoleMiddleware("admin"), controller.GetUserByID)
	auth.PUT("/users/:id", middleware.RoleMiddleware("admin"), controller.UpdateUser)
	auth.DELETE("/users/:id", middleware.RoleMiddleware("admin"), controller.DeleteUser)
	auth.GET("/users/profile", middleware.RoleMiddleware("customer"), controller.GetUserProfile)

}
