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
	var kurir []model.User

	// Subquery: kurir_id yang sedang memproses order
	subQuery := config.DB.
		Model(&model.Order{}).
		Select("kurir_id").
		Where("status = ?", "proses")

	// Main query: cari kurir yang online dan tidak ada di subquery
	err := config.DB.
		Model(&model.User{}).
		Where("role = ? AND status = ?", "kurir", "online").
		Where("id NOT IN (?)", subQuery).
		Find(&kurir).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kurir)
}

func GetKurirByID(c *gin.Context) {
	idParam := c.Param("id")

	var user model.User
	if err := config.DB.Where("id = ? AND role = ?", idParam, "kurir").First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kurir tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
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
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := config.DB.Delete(&model.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}
