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
	Stock       int            `gorm:"not null;default:0" json:"stock"`
	Price       float64        `gorm:"not null;type:decimal(10,2)" json:"price"`
	CategoryID  uint           `gorm:"index" json:"category_id,omitempty"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Categories  []Category     `gorm:"many2many:product_categories;" json:"categories,omitempty"`
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
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	SKU         string             `json:"sku"`
	Stock       int                `json:"stock"`
	Price       float64            `json:"price"`
	CategoryID  uint               `json:"category_id,omitempty"`
	Category    CategoryResponse   `json:"category,omitempty"`
	Categories  []CategoryResponse `json:"categories,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ToResponse converts Product to ProductResponse
func (p *Product) ToResponse() ProductResponse {
	resp := ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		SKU:         p.SKU,
		Stock:       p.Stock,
		Price:       p.Price,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Category.ID != 0 {
		resp.Category = p.Category.ToResponse()
	}
	if len(p.Categories) > 0 {
		resp.Categories = make([]CategoryResponse, len(p.Categories))
		for i, cat := range p.Categories {
			resp.Categories[i] = cat.ToResponse()
		}
	}
	return resp
}

// CreateProductRequest is the DTO for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=200"`
	Description string  `json:"description" binding:"max=1000"`
	SKU         string  `json:"sku" binding:"required,min=1,max=50"`
	Stock       int     `json:"stock" binding:"gte=0"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	CategoryID  uint    `json:"category_id"`
	CategoryIDs []uint  `json:"category_ids"`
}

// UpdateProductRequest is the DTO for updating a product
type UpdateProductRequest struct {
	Name        string   `json:"name" binding:"omitempty,min=1,max=200"`
	Description string   `json:"description" binding:"max=1000"`
	SKU         string   `json:"sku" binding:"omitempty,min=1,max=50"`
	Stock       *int     `json:"stock"`
	Price       *float64 `json:"price"`
	CategoryID  uint     `json:"category_id" binding:"omitempty"`
	CategoryIDs []uint   `json:"category_ids"`
}

// UpdateStockRequest is the DTO for updating product stock
type UpdateStockRequest struct {
	Stock int `json:"stock" binding:"required"`
}

