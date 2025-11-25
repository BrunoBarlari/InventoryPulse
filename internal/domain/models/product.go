package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:200" json:"name"`
	Description string         `gorm:"size:1000" json:"description"`
	SKU         string         `gorm:"uniqueIndex;not null;size:50" json:"sku"`
	Quantity    int            `gorm:"not null;default:0" json:"quantity"`
	Price       float64        `gorm:"not null;type:decimal(10,2)" json:"price"`
	CategoryID  uint           `gorm:"not null;index" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}

// ProductResponse is the DTO for product responses
type ProductResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	SKU         string           `json:"sku"`
	Quantity    int              `json:"quantity"`
	Price       float64          `json:"price"`
	CategoryID  uint             `json:"category_id"`
	Category    CategoryResponse `json:"category,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// ToResponse converts Product to ProductResponse
func (p *Product) ToResponse() ProductResponse {
	resp := ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		SKU:         p.SKU,
		Quantity:    p.Quantity,
		Price:       p.Price,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Category.ID != 0 {
		resp.Category = p.Category.ToResponse()
	}
	return resp
}

// CreateProductRequest is the DTO for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=200"`
	Description string  `json:"description" binding:"max=1000"`
	SKU         string  `json:"sku" binding:"required,min=1,max=50"`
	Quantity    int     `json:"quantity" binding:"gte=0"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

// UpdateProductRequest is the DTO for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"omitempty,min=1,max=200"`
	Description string  `json:"description" binding:"max=1000"`
	SKU         string  `json:"sku" binding:"omitempty,min=1,max=50"`
	Quantity    int     `json:"quantity" binding:"gte=0"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	CategoryID  uint    `json:"category_id" binding:"omitempty"`
}

// UpdateStockRequest is the DTO for updating product stock
type UpdateStockRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

