package models

import (
	"time"
)

// ProductHistory tracks changes to product price and stock
type ProductHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"not null;index" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"-"`
	Price     float64   `gorm:"not null;type:decimal(10,2)" json:"price"`
	Stock     int       `gorm:"not null" json:"stock"`
	ChangedAt time.Time `gorm:"not null;index" json:"changed_at"`
}

// TableName specifies the table name for ProductHistory model
func (ProductHistory) TableName() string {
	return "product_history"
}

// ProductHistoryResponse is the DTO for product history responses
type ProductHistoryResponse struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	ChangedAt time.Time `json:"changed_at"`
}

// ToResponse converts ProductHistory to ProductHistoryResponse
func (h *ProductHistory) ToResponse() ProductHistoryResponse {
	return ProductHistoryResponse{
		ID:        h.ID,
		ProductID: h.ProductID,
		Price:     h.Price,
		Stock:     h.Stock,
		ChangedAt: h.ChangedAt,
	}
}

// ProductHistoryQuery is the DTO for history query parameters
type ProductHistoryQuery struct {
	PaginationRequest
	Start string `form:"start"` // Format: YYYY-MM-DD or RFC3339
	End   string `form:"end"`   // Format: YYYY-MM-DD or RFC3339
}

