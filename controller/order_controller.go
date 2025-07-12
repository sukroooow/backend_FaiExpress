package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mubarok-ridho/misi-paket.backend/config"
	"github.com/mubarok-ridho/misi-paket.backend/model"
)

// 🔸 Create Order
func CreateOrder(c *gin.Context) {
	var input model.Order
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi kurir_id tidak null atau 0
	if input.KurirID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kurir tidak boleh kosong"})
		return
	}

	// Set status default
	if input.Status == "" {
		input.Status = "proses"
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update status kurir jadi offline
	if input.KurirID != 0 {
		config.DB.Model(&model.User{}).Where("id = ?", input.KurirID).Update("status", "offline")
	}

	// ✅ Return order ID dan pesan
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Pesanan berhasil dibuat",
		"order_id": input.ID, // ambil ID dari input setelah di-insert
	})
}

// 🔸 Get All Orders
func GetAllOrders(c *gin.Context) {
	var orders []model.Order
	if err := config.DB.
		Preload("Customer").
		Preload("Kurir").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// 🔸 Get Order by ID
func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order model.Order

	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// 🔸 Update Order
func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order model.Order

	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if order.Status == "selesai" {
		config.DB.Model(&model.User{}).Where("id = ?", order.KurirID).Update("status", "online")
	}

	c.JSON(http.StatusOK, order)
}

// 🔸 Delete Order
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := config.DB.Delete(&model.Order{}, i).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order berhasil dihapus"})
}

// 🔸 Get My Orders (Customer)
func GetMyOrders(c *gin.Context) {
	userID := c.GetUint("userID")
	var orders []model.Order
	if err := config.DB.Where("customer_id = ?", userID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// 🔸 Update Kurir Location
// 🔸 Update Kurir Location
func UpdateLocation(c *gin.Context) {
	var loc struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	if err := c.ShouldBindJSON(&loc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format lokasi salah"})
		return
	}

	// Simpan lokasi ke cache, Redis, atau memori (dummy response dulu)
	c.JSON(http.StatusOK, gin.H{"message": "Lokasi diperbarui"})
}

// 🔸 Send Chat
func SendChat(c *gin.Context) {
	var msg struct {
		OrderID uint   `json:"order_id"`
		From    string `json:"from"` // "customer" or "kurir"
		Text    string `json:"text"`
	}

	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format chat salah"})
		return
	}

	// Dummy response, seharusnya disimpan ke DB atau cache
	c.JSON(http.StatusOK, gin.H{"message": "Pesan terkirim"})
}

// 🔸 Get Chat (Customer Side)
func GetChat(c *gin.Context) {
	orderID := c.Query("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id wajib diisi"})
		return
	}

	// Dummy response
	c.JSON(http.StatusOK, gin.H{
		"messages": []map[string]string{
			{"from": "kurir", "text": "Saya OTW!"},
			{"from": "customer", "text": "Baik kak, hati-hati"},
		},
	})
}

func UpdateOrderStatus(c *gin.Context) {
	var input struct {
		ID     uint   `json:"id"`
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	if err := config.DB.Model(&model.Order{}).
		Where("id = ?", input.ID).
		Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status berhasil diperbarui"})
}

// 🔹 GET /kurir/:id/orders
func GetOrdersForKurir(c *gin.Context) {
	kurirID := c.Param("id")
	var orders []model.Order

	if err := config.DB.Preload("Customer"). // untuk ambil nama customer
							Where("kurir_id = ? AND status = ?", kurirID, "proses").
							Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data order"})
		return
	}

	// 🔄 Format response yang lengkap
	var response []map[string]interface{}
	for _, order := range orders {
		// Ambil nama_order tergantung layanannya
		var namaOrder string
		if order.Layanan == "barang" && order.NamaBarang != nil {
			namaOrder = *order.NamaBarang
		} else if order.Layanan == "makanan" && order.NamaMakanan != nil {
			namaOrder = *order.NamaMakanan
		}

		response = append(response, map[string]interface{}{
			"id":            order.ID,
			"layanan":       order.Layanan,
			"status":        order.Status,
			"alamat_jemput": order.AlamatJemput,
			"alamat_antar":  order.AlamatAntar,
			"nama_order":    namaOrder,           // ✅ ini nama barang/makanan
			"nama_customer": order.Customer.Name, // ✅ ini nama customer
		})
	}

	c.JSON(http.StatusOK, response)
}
