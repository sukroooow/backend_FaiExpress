package controller

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// In-memory store untuk menyimpan lokasi kurir
var kurirLocations = struct {
	sync.RWMutex
	data map[uint]map[string]float64 // kurir_id -> {"lat": ..., "lng": ...}
}{
	data: make(map[uint]map[string]float64),
}

// POST /kurir/track
func UpdateKurirLocation(c *gin.Context) {
	var req struct {
		KurirID uint    `json:"kurir_id"`
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	kurirLocations.Lock()
	kurirLocations.data[req.KurirID] = map[string]float64{
		"lat": req.Lat,
		"lng": req.Lng,
	}
	kurirLocations.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Lokasi kurir diperbarui"})
}

// GET /kurir/track/:id
func GetKurirLocation(c *gin.Context) {
	var req struct {
		KurirID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kurir tidak valid"})
		return
	}

	kurirLocations.RLock()
	loc, found := kurirLocations.data[req.KurirID]
	kurirLocations.RUnlock()

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lokasi belum tersedia"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"lat": loc["lat"],
		"lng": loc["lng"],
	})
}
