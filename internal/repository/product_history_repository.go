package repository

import (
	"time"

	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"gorm.io/gorm"
)

type ProductHistoryRepository interface {
	Create(history *models.ProductHistory) error
	FindByProductID(productID uint, start, end *time.Time, page, pageSize int) ([]models.ProductHistory, int64, error)
	GetLatestByProductID(productID uint) (*models.ProductHistory, error)
}

type productHistoryRepository struct {
	db *gorm.DB
}

func NewProductHistoryRepository(db *gorm.DB) ProductHistoryRepository {
	return &productHistoryRepository{db: db}
}

func (r *productHistoryRepository) Create(history *models.ProductHistory) error {
	return r.db.Create(history).Error
}

func (r *productHistoryRepository) FindByProductID(productID uint, start, end *time.Time, page, pageSize int) ([]models.ProductHistory, int64, error) {
	var history []models.ProductHistory
	var total int64

	query := r.db.Model(&models.ProductHistory{}).Where("product_id = ?", productID)

	// Apply date filters
	if start != nil {
		query = query.Where("changed_at >= ?", *start)
	}
	if end != nil {
		// Add 1 day to end date to include the entire day
		endPlusDay := end.Add(24 * time.Hour)
		query = query.Where("changed_at < ?", endPlusDay)
	}

	// Get total count
	query.Count(&total)

	// Apply pagination and ordering
	offset := (page - 1) * pageSize
	err := query.Order("changed_at DESC").Offset(offset).Limit(pageSize).Find(&history).Error
	if err != nil {
		return nil, 0, err
	}

	return history, total, nil
}

func (r *productHistoryRepository) GetLatestByProductID(productID uint) (*models.ProductHistory, error) {
	var history models.ProductHistory
	err := r.db.Where("product_id = ?", productID).Order("changed_at DESC").First(&history).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &history, nil
}

