package repository

import (
	"errors"

	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductSKUExists     = errors.New("product with this SKU already exists")
	ErrInvalidCategory      = errors.New("invalid category")
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	FindBySKU(sku string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	List(page, pageSize int, categoryID *uint) ([]models.Product, int64, error)
	UpdateStock(id uint, quantity int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	// Check if SKU already exists
	var count int64
	r.db.Model(&models.Product{}).Where("sku = ?", product.SKU).Count(&count)
	if count > 0 {
		return ErrProductSKUExists
	}

	// Verify category exists
	var catCount int64
	r.db.Model(&models.Category{}).Where("id = ?", product.CategoryID).Count(&catCount)
	if catCount == 0 {
		return ErrInvalidCategory
	}

	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Where("sku = ?", sku).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *models.Product) error {
	// Check if another product has the same SKU
	var count int64
	r.db.Model(&models.Product{}).Where("sku = ? AND id != ?", product.SKU, product.ID).Count(&count)
	if count > 0 {
		return ErrProductSKUExists
	}

	// Verify category exists
	var catCount int64
	r.db.Model(&models.Category{}).Where("id = ?", product.CategoryID).Count(&catCount)
	if catCount == 0 {
		return ErrInvalidCategory
	}

	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Product{}, id)
	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}
	return result.Error
}

func (r *productRepository) List(page, pageSize int, categoryID *uint) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	if categoryID != nil && *categoryID > 0 {
		query = query.Where("category_id = ?", *categoryID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Category").Offset(offset).Limit(pageSize).Order("id ASC").Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) UpdateStock(id uint, quantity int) error {
	result := r.db.Model(&models.Product{}).Where("id = ?", id).Update("quantity", quantity)
	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}
	return result.Error
}

