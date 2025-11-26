package models

// ProductCategory represents the many-to-many relationship between products and categories
type ProductCategory struct {
	ProductID  uint     `gorm:"primaryKey" json:"product_id"`
	CategoryID uint     `gorm:"primaryKey" json:"category_id"`
	Product    Product  `gorm:"foreignKey:ProductID" json:"-"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"-"`
}

// TableName specifies the table name for ProductCategory model
func (ProductCategory) TableName() string {
	return "product_categories"
}

