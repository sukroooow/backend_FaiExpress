package model

type User struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	Name             string  `json:"name"`
	Email            string  `gorm:"unique" json:"email"`
	Password         string  `json:"password"`
	Role             string  `json:"role"` // admin, kurir, customer
	Phone            string  `json:"phone"`
	Kendaraan        *string `json:"kendaraan,omitempty"` // nullable (hanya kurir)
	OrdersAsCustomer []Order `gorm:"foreignKey:CustomerID" json:"orders_as_customer,omitempty"`
	OrdersAsKurir    []Order `gorm:"foreignKey:KurirID" json:"orders_as_kurir,omitempty"`
	Status           string  `json:"status"` // online, offline
}
