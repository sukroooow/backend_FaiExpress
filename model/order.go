package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CustomerID   uint           `json:"customer_id"`
	KurirID      uint           `json:"kurir_id"` // bisa null, makanya pointer
	AlamatJemput string         `json:"alamat_jemput"`
	AlamatAntar  string         `json:"alamat_antar"`
	NamaBarang   *string        `json:"nama_barang"` // pakai pointer kalau bisa null
	NamaMakanan  *string        `json:"nama_makanan"`
	Status       string         `json:"status"`
	Layanan      string         `json:"layanan"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
