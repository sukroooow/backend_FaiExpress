package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mubarok-ridho/misi-paket.backend/config"
	"github.com/mubarok-ridho/misi-paket.backend/model"
	// "github.com/mubarok-ridho/misi-paket.backend/config"
	// "github.com/mubarok-ridho/misi-paket.backend/model"
)

// Struct untuk payload publish ke Centrifugo
type CentrifugoPublishPayload struct {
	Method string                  `json:"method"`
	Params CentrifugoPublishParams `json:"params"`
}

type SendChatInput struct {
	OrderIDStr string `json:"order_id" binding:"required"`
	SenderID   uint   `json:"sender_id" binding:"required"`
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Sender     string `json:"sender" binding:"required"`
	Content    string `json:"message" binding:"required"`
}
type CentrifugoPublishParams struct {
	Channel string                 `json:"channel"`
	Data    map[string]interface{} `json:"data"`
}

// Struct untuk menerima body request Flutter
type ChatMessage struct {
	OrderID string `json:"order_id"`
	Sender  string `json:"sender"` // "kurir" atau "customer"
	Message string `json:"message"`
}

// Handler kirim chat ke Centrifugo (POST /chat/send)
func SendChatMessage(c *gin.Context) {
	var input SendChatInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderID, err := strconv.ParseUint(input.OrderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id harus berupa angka"})
		return
	}
	// 1. Simpan ke databaseh
	newMessage := model.Message{
		OrderID:    uint(orderID),
		SenderID:   input.SenderID,
		ReceiverID: input.ReceiverID,
		Content:    input.Content,
		SentAt:     time.Now(),
		IsRead:     false,
	}

	if err := config.DB.Create(&newMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pesan"})
		return
	}

	// Send ke Centrifugow
	msg := ChatMessage{
		OrderID: fmt.Sprint(orderID),
		Sender:  input.Sender, // ✅ penting!
		Message: input.Content,
	}

	channel := fmt.Sprintf("chat:%s", msg.OrderID)

	payload := CentrifugoPublishParams{
		Channel: channel,
		Data: map[string]interface{}{
			"sender":  msg.Sender,
			"message": msg.Message,
			"time":    time.Now().Format(time.RFC3339),
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode payload"})
		return
	}
	fmt.Println(string(jsonPayload))
	centrifugoURL := os.Getenv("CENTRIFUGO_API_URL")    // contoh: http://localhost:8000/api
	centrifugoAPIKey := os.Getenv("CENTRIFUGO_API_KEY") // contoh: secret_api_key

	req, err := http.NewRequest("POST", centrifugoURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", centrifugoAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send to Centrifugo"})
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "centrifugo publish failed",
			"details": string(bodyBytes),
		})
		return
	}
	fmt.Println(string(bodyBytes))
	c.JSON(http.StatusOK, gin.H{"status": "message sent"})
}

// Handler untuk generate token JWT Centrifugo (GET /centrifugo/token?user_id=xxx)
func GenerateCentrifugoToken(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	secret := os.Getenv("CENTRIFUGO_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "missing CENTRIFUGO_SECRET env"})
		return
	}

	// Membuat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token signing failed"})
		return
	}

	fmt.Println("DEBUG Backend Generated JWT:", tokenString) // Tambahkan ini

	// ✅ Debug log untuk development
	if os.Getenv("APP_ENV") == "development" {
		log.Println("✅ JWT token generated with:")
		log.Println("  - userID:", userID)
		log.Println("  - secret : [", secret, "]") // HATI-HATI! hanya untuk dev!
		log.Println("  - token  :", tokenString)
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"userId": userID,
	})
}

func GetMessagesByOrderID(c *gin.Context) {
	orderIDParam := c.Param("order_id")
	orderID, err := strconv.ParseUint(orderIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id tidak valid"})
		return
	}

	var messages []model.Message
	if err := config.DB.
		Preload("Sender").
		Where("order_id = ?", orderID).
		Order("sent_at ASC").
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil pesan dari database"})
		return
	}

	// Boleh langsung kirim semua data lengkap
	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}
