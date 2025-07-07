package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mubarok-ridho/misi-paket.backend/controller"
	"github.com/mubarok-ridho/misi-paket.backend/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// Root
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "FaiExpress API is running!"})
	})

	// Auth
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.POST("/kurir/track", controller.UpdateKurirLocation)
	r.GET("/kurir/track/:id", controller.GetKurirLocation)
	r.PUT("/users/:id/password", controller.ChangePassword)
	r.RegisterChatr(r)
	r.RegisterWebSocketRoutes(r)
	// Public
	r.GET("/kurir/available", controller.GetAvailableKurir)
	r.GET("/kurir/:id/location", controller.GetKurirLocation)

	// 🔐 JWT Protected Group
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuthMiddleware())
	auth.GET("/kurir/:id/orders", middleware.RoleMiddleware("kurir"), controller.GetOrdersForKurir)

	// 📦 Orders (customer only for create + my-orders)
	auth.POST("/orders", middleware.RoleMiddleware("customer"), controller.CreateOrder)
	auth.GET("/my-orders", middleware.RoleMiddleware("customer"), controller.GetMyOrders)

	// 🔧 Order management (admin or kurir, as needed)
	auth.GET("/orders", middleware.RoleMiddleware("admin"), controller.GetAllOrders)
	auth.GET("/orders/:id", controller.GetOrderByID)
	auth.PUT("/orders/:id", controller.UpdateOrder)
	auth.DELETE("/orders/:id", middleware.RoleMiddleware("admin"), controller.DeleteOrder)

	// 🛵 Kurir only
	auth.PUT("/kurir/status", middleware.RoleMiddleware("kurir"), controller.UpdateKurirStatus)
	auth.PUT("/kurir/location", middleware.RoleMiddleware("kurir"), controller.UpdateLocation)

	// 🗨️ Chat (kurir/customer)
	auth.POST("/chat", middleware.RoleMiddleware("customer", "kurir"), controller.SendChat)
	auth.GET("/chat", middleware.RoleMiddleware("customer", "kurir"), controller.GetChat)

	// 👤 User CRUD (admin only)
	auth.GET("/users", middleware.RoleMiddleware("admin"), controller.GetAllUsers)
	auth.GET("/users/:id", middleware.RoleMiddleware("admin"), controller.GetUserByID)
	auth.PUT("/users/:id", middleware.RoleMiddleware("admin"), controller.UpdateUser)
	auth.DELETE("/users/:id", middleware.RoleMiddleware("admin"), controller.DeleteUser)
}
