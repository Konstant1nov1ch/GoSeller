package model

type Product struct {
	ID         int64  `gorm:"column:id;primaryKey" json:"Id"`
	Name       string `gorm:"column:product_name" json:"Name"`
	SalePriceU int64  `gorm:"column:product_price" json:"SalePriceU"`
}
