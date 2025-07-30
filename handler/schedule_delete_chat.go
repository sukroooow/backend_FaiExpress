package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mubarok-ridho/misi-paket.backend/config"
	"github.com/mubarok-ridho/misi-paket.backend/model"
)

// ScheduleDeleteChat menjadwalkan penghapusan chat berdasarkan order_id setelah 24 jam
func ScheduleDeleteChat(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID tidak ditemukan"})
		return
	}

	// Langsung respons ke user
	c.JSON(http.StatusOK, gin.H{
		"message":   "Penghapusan chat dijadwalkan dalam 24 jam",
		"order_id":  orderID,
		"scheduled": true,
	})

	// Jalankan goroutine
	go func(orderID string) {
		log.Printf("[INFO] Menunggu 24 jam sebelum hapus chat order %s", orderID)
		time.Sleep(24 * time.Hour)

		if err := config.DB.Where("order_id = ?", orderID).Delete(&model.Message{}).Error; err != nil {
			log.Printf("[ERROR] Gagal hapus chat order %s: %v", orderID, err)
		} else {
			log.Printf("[INFO] Berhasil hapus chat order %s setelah 24 jam", orderID)
		}
	}(orderID)
}
