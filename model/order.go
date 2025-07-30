package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	CustomerID uint `json:"customer_id"`
	Customer   User `gorm:"foreignKey:CustomerID" json:"customer"`

	KurirID uint `json:"kurir_id"`
	Kurir   User `gorm:"foreignKey:KurirID" json:"kurir"`

	MetodeBayar string `json:"metode_bayar"`

	Status        string  `json:"status"`
	Layanan       string  `json:"layanan"`
	Nominal       *uint   `json:"nominal"`        // 💰 Total tagihan (opsional)
	PaymentStatus *string `json:"payment_status"` // ⏳ "pending" atau ✅ "done"

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// func (Order) TableName() string {
// 	return "public.orders"
// }
