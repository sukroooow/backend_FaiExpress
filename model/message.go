// model/messages.go
package model

import "time"

type Message struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `json:"order_id"`
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
	IsRead     bool      `json:"is_read"`

	Sender User `gorm:"foreignKey:SenderID"` // 👈

}

func (Message) TableName() string {
	return "public.messages"
}
