package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubarok-ridho/misi-paket.backend/config"
	"github.com/mubarok-ridho/misi-paket.backend/model"
)

// GET /users
func GetAllUsers(c *gin.Context) {
	var users []model.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetAvailableKurir(c *gin.Context) {
	var kurirs []model.User

	// Ambil semua kurir online
	if err := config.DB.
		Where("role = ? AND status = ? AND status_kerja = ?", "kurir", "online", "aktif").
		Find(&kurirs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var filtered []map[string]interface{}
	for _, kurir := range kurirs {
		var count int64
		config.DB.
			Model(&model.Order{}).
			Where("kurir_id = ? AND status = ?", kurir.ID, "proses").
			Count(&count)

		if count < 5 {
			filtered = append(filtered, map[string]interface{}{
				"id":             kurir.ID,
				"name":           kurir.Name,
				"kendaraan":      kurir.Kendaraan,
				"jumlah_pesanan": count,
				"no_hp":          kurir.Phone,
			})
		}
	}

	c.JSON(http.StatusOK, filtered)
}

func GetKurirByID(c *gin.Context) {
	idParam := c.Param("id")

	var user model.User
	if err := config.DB.First(&user, "id = ? AND role = ?", idParam, "kurir").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kurir tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"phone":      user.Phone,
			"kendaraan":  user.Kendaraan,
			"plat_nomor": user.PlatNomor, // ← tambahkan ini
			"status":     user.Status,
		},
	})
}

// PUT /api/kurir/:id
func UpdateKurirByID(c *gin.Context) {
	idParam := c.Param("id")

	var user model.User
	if err := config.DB.First(&user, "id = ? AND role = ?", idParam, "kurir").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kurir tidak ditemukan"})
		return
	}

	var input struct {
		Name      string  `json:"name"`
		Phone     string  `json:"phone"`
		Kendaraan *string `json:"kendaraan"`
		Email     string  `json:"email"`
		PlatNomor *string `json:"plat_nomor"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	user.Name = input.Name
	user.Phone = input.Phone
	user.Kendaraan = input.Kendaraan
	user.Email = input.Email
	user.PlatNomor = input.PlatNomor

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update kurir"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profil kurir berhasil diperbarui"})
}

func GetUserProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID uint
	switch id := userIDInterface.(type) {
	case uint:
		userID = id
	case int:
		userID = uint(id)
	case float64:
		userID = uint(id)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"phone": user.Phone,
		"role":  user.Role,
	})
}

func UpdateKurirStatus(c *gin.Context) {
	var input struct {
		ID     uint   `json:"id"`     // sementara kita pakai ID dari body
		Status string `json:"status"` // online / offline
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if err := config.DB.Model(&model.User{}).Where("id = ? AND role = ?", input.ID, "kurir").
		Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status berhasil diperbarui"})
}

// GET /users/:id
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID uint
	switch id := userIDInterface.(type) {
	case uint:
		userID = id
	case int:
		userID = uint(id)
	case float64:
		userID = uint(id)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var input struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Update field
	user.Name = input.Name
	user.Phone = input.Phone
	user.Email = input.Email

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profil berhasil diperbarui"})
}

// PUT /users/:id
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DELETE /users/:id
func SoftDeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	user.StatusKerja = "nonaktif"
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kurir berhasil dinonaktifkan"})
}
