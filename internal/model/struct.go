package model

type Product struct {
	ID         int64  `gorm:"product_id" json:"Id"`
	Name       string `gorm:"product_name" json:"Name"`
	SalePriceU int64  `gorm:"product_price" json:"SalePriceU"`
}
