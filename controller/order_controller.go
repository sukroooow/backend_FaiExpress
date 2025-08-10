package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	// if input.KurirID != 0 {
	// 	config.DB.Model(&model.User{}).Where("id = ?", input.KurirID).Update("status", "offline")
	// }

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
	if err := config.DB.
		Preload("Kurir").
		Preload("Customer").
		First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	var activeCount int64
	config.DB.Model(&model.Order{}).
		Where("kurir_id = ? AND status = ?", order.KurirID, "proses").
		Count(&activeCount)

	kurirData := map[string]interface{}{
		"id":            order.Kurir.ID,
		"name":          order.Kurir.Name,
		"vehicle":       order.Kurir.Kendaraan,
		"active_orders": activeCount,
	}

	// return response dengan kurir info yang dilengkapi
	c.JSON(http.StatusOK, gin.H{
		"order":          order,
		"kurir":          kurirData,
		"user_id":        order.Customer.ID, // ✅ tambahkan ini
		"payment_status": order.PaymentStatus,
		"tagihan":        order.Nominal,
		"MetodeBayar":    order.MetodeBayar,
	})
}

func DeleteMessagesByOrderID(c *gin.Context) {
	orderID := c.Param("id")

	if err := config.DB.Where("order_id = ?", orderID).Delete(&model.Message{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pesan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Semua pesan berhasil dihapus"})
}

func UpdatePaymentMethod(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Method string `json:"metode_bayar"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var order model.Order
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.MetodeBayar = req.Method // ✅ langsung assign string, bukan pointer

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Metode bayar diperbarui"})
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

	if err := config.DB.
		Preload("Kurir"). // ⬅️ PENTING: preload kurir
		Where("customer_id = ?", userID).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

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

	// Pakai struct agar UpdatedAt ikut berubah otomatis
	updateData := model.Order{
		Status: input.Status,
	}

	if err := config.DB.Model(&model.Order{}).
		Where("id = ?", input.ID).
		Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status berhasil diperbarui"})
}

func CheckOrderKurirReady(c *gin.Context) {
	orderID := c.Param("id")
	var order model.Order

	if err := config.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan"})
		return
	}

	tagihanSiap := *order.Nominal > 0
	metodeBayarDiisi := order.MetodeBayar != ""

	c.JSON(http.StatusOK, gin.H{
		"tagihan_siap":       tagihanSiap,
		"metode_bayar_diisi": metodeBayarDiisi,
		"bisa_lanjut":        tagihanSiap && metodeBayarDiisi,
	})
}

func UpdateTagihan(c *gin.Context) {
	type RincianItem struct {
		Judul   string `json:"judul"`
		Nominal int    `json:"nominal"`
	}

	var req struct {
		ID      uint          `json:"id"`
		Nominal int           `json:"nominal"`
		Rincian []RincianItem `json:"rincian"` // diterima, tapi tidak disimpan
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var order model.Order
	if err := config.DB.First(&order, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	nominal := uint(req.Nominal)
	status := "pending"

	order.Nominal = &nominal
	order.PaymentStatus = &status

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	// Cetak rincian ke log (tidak disimpan ke DB)
	for _, item := range req.Rincian {
		fmt.Printf("🔹 Rincian: %s = %d\n", item.Judul, item.Nominal)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tagihan berhasil diperbarui (rincian tidak disimpan)",
	})
}

func ValidasiPembayaran(c *gin.Context) {
	var req struct {
		ID uint `json:"id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := config.DB.Model(&model.Order{}).
		Where("id = ?", req.ID).
		Update("payment_status", "done").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran divalidasi"})
}

func GetOrdersProses(c *gin.Context) {
	kurirID := c.Param("id")
	var orders []model.Order

	if err := config.DB.Preload("Customer").
		Where("kurir_id = ? AND status = ?", kurirID, "proses").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data order"})
		return
	}

	var response []map[string]interface{}
	for _, order := range orders {
		response = append(response, map[string]interface{}{
			"id":            order.ID,
			"layanan":       order.Layanan,
			"status":        order.Status,
			"nama_order":    fmt.Sprintf("Order #%d", order.ID),
			"nama_customer": order.Customer.Name,
			"customer_id":   order.CustomerID,
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetTotalPendapatanToday(c *gin.Context) {
	var totalPendapatan float64

	// Ambil waktu sekarang
	now := time.Now()
	// Tentukan awal dan akhir hari
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Second)

	// Hitung total biaya dari order yang statusnya selesai hari ini
	err := config.DB.Model(&model.Order{}).
		Where("status = ? AND updated_at BETWEEN ? AND ?", "selesai", startOfDay, endOfDay).
		Select("COALESCE(SUM(biaya), 0)"). // pakai COALESCE supaya hasilnya 0 kalau tidak ada data
		Scan(&totalPendapatan).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total pendapatan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_pendapatan": totalPendapatan})
}

func GetPendapatanKurirToday(c *gin.Context) {
	kurirID := c.Param("id")
	today := time.Now().Format("2006-01-02")

	var totalPendapatan float64

	err := config.DB.
		Model(&model.Order{}).
		Select("COALESCE(SUM(nominal), 0)").
		Where("kurir_id = ? AND status = ? AND DATE(updated_at) = ?", kurirID, "selesai", today).
		Scan(&totalPendapatan).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pendapatan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_pendapatan": totalPendapatan,
	})
}

func GetAllTotalPendapatanToday(c *gin.Context) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	today := time.Now().In(loc).Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	var totalPendapatan float64

	err := config.DB.Model(&model.Order{}).
		Where("status = ? AND updated_at BETWEEN ? AND ?", "selesai", today, tomorrow).
		Select("COALESCE(SUM(nominal), 0)"). // pakai nominal karena kamu pakai itu untuk tagihan kurir
		Scan(&totalPendapatan).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total pendapatan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_pendapatan": totalPendapatan,
	})
}

func GetTotalOrdersSelesaiToday(c *gin.Context) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	today := time.Now().In(loc).Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	var totalOrders int64
	err := config.DB.Model(&model.Order{}).
		Where("status = ? AND updated_at BETWEEN ? AND ?", "selesai", today, tomorrow).
		Count(&totalOrders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung pesanan selesai"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_orders_selesai": totalOrders,
	})
}

func GetOrdersSelesaiToday(c *gin.Context) {
	kurirID := c.Param("id")
	var orders []model.Order

	loc, _ := time.LoadLocation("Asia/Jakarta")
	today := time.Now().In(loc).Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	if err := config.DB.Preload("Customer").
		Where("kurir_id = ? AND status = ? AND updated_at BETWEEN ? AND ?", kurirID, "selesai", today, tomorrow).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data order selesai hari ini"})
		return
	}

	var response []map[string]interface{}
	for _, order := range orders {
		response = append(response, map[string]interface{}{
			"id":            order.ID,
			"layanan":       order.Layanan,
			"status":        order.Status,
			"nama_order":    fmt.Sprintf("Order #%d", order.ID),
			"nama_customer": order.Customer.Name,
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetOrdersForKurir(c *gin.Context) {
	kurirID := c.Param("id")
	var orders []model.Order

	// Ambil semua pesanan kurir (proses dan selesai)
	if err := config.DB.Preload("Customer").
		Where("kurir_id = ?", kurirID).
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data order"})
		return
	}

	var response []map[string]interface{}
	for _, order := range orders {
		// Safe access untuk pointer
		nominal := uint(0)
		if order.Nominal != nil {
			nominal = *order.Nominal
		}

		paymentStatus := "pending"
		if order.PaymentStatus != nil {
			paymentStatus = *order.PaymentStatus
		}

		response = append(response, map[string]interface{}{
			"id":             order.ID,
			"layanan":        order.Layanan,
			"status":         order.Status,
			"nominal":        nominal,
			"payment_status": paymentStatus,
			"updated_at":     order.UpdatedAt.Format("2006-01-02"),
			"nama_order":     fmt.Sprintf("Order #%d", order.ID),
			"nama_customer":  order.Customer.Name,
		})
	}

	c.JSON(http.StatusOK, response)
}
