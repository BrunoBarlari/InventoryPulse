package repository

import (
	"errors"

	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrProductSKUExists = errors.New("product with this SKU already exists")
	ErrInvalidCategory  = errors.New("invalid category")
)

type ProductRepository interface {
	Create(product *models.Product, categoryIDs []uint) error
	FindByID(id uint) (*models.Product, error)
	FindBySKU(sku string) (*models.Product, error)
	Update(product *models.Product, categoryIDs []uint) error
	Delete(id uint) error
	List(page, pageSize int, categoryID *uint, search string) ([]models.Product, int64, error)
	UpdateStock(id uint, stock int) error
	Search(query string, page, pageSize int) ([]models.Product, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product, categoryIDs []uint) error {
	// Check if SKU already exists
	var count int64
	r.db.Model(&models.Product{}).Where("sku = ?", product.SKU).Count(&count)
	if count > 0 {
		return ErrProductSKUExists
	}

	// Verify primary category exists if provided
	if product.CategoryID > 0 {
		var catCount int64
		r.db.Model(&models.Category{}).Where("id = ?", product.CategoryID).Count(&catCount)
		if catCount == 0 {
			return ErrInvalidCategory
		}
	}

	// Verify all category IDs exist
	if len(categoryIDs) > 0 {
		var validCount int64
		r.db.Model(&models.Category{}).Where("id IN ?", categoryIDs).Count(&validCount)
		if int(validCount) != len(categoryIDs) {
			return ErrInvalidCategory
		}
	}

	// Create product
	if err := r.db.Create(product).Error; err != nil {
		return err
	}

	// Associate categories
	if len(categoryIDs) > 0 {
		var categories []models.Category
		r.db.Where("id IN ?", categoryIDs).Find(&categories)
		if err := r.db.Model(product).Association("Categories").Replace(categories); err != nil {
			return err
		}
	}

	return nil
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Preload("Categories").First(&product, id).Error
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
	err := r.db.Preload("Category").Preload("Categories").Where("sku = ?", sku).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *models.Product, categoryIDs []uint) error {
	// Check if another product has the same SKU
	var count int64
	r.db.Model(&models.Product{}).Where("sku = ? AND id != ?", product.SKU, product.ID).Count(&count)
	if count > 0 {
		return ErrProductSKUExists
	}

	// Verify primary category exists if provided
	if product.CategoryID > 0 {
		var catCount int64
		r.db.Model(&models.Category{}).Where("id = ?", product.CategoryID).Count(&catCount)
		if catCount == 0 {
			return ErrInvalidCategory
		}
	}

	// Verify all category IDs exist
	if len(categoryIDs) > 0 {
		var validCount int64
		r.db.Model(&models.Category{}).Where("id IN ?", categoryIDs).Count(&validCount)
		if int(validCount) != len(categoryIDs) {
			return ErrInvalidCategory
		}
	}

	// Update product
	if err := r.db.Save(product).Error; err != nil {
		return err
	}

	// Update categories association if provided
	if len(categoryIDs) > 0 {
		var categories []models.Category
		r.db.Where("id IN ?", categoryIDs).Find(&categories)
		if err := r.db.Model(product).Association("Categories").Replace(categories); err != nil {
			return err
		}
	}

	return nil
}

func (r *productRepository) Delete(id uint) error {
	// First remove category associations
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	// Clear associations
	r.db.Model(&product).Association("Categories").Clear()

	// Delete product
	result := r.db.Delete(&models.Product{}, id)
	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}
	return result.Error
}

func (r *productRepository) List(page, pageSize int, categoryID *uint, search string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	// Filter by primary category
	if categoryID != nil && *categoryID > 0 {
		query = query.Where("category_id = ?", *categoryID)
	}

	// Search by name or SKU (case-insensitive)
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?) OR LOWER(sku) LIKE LOWER(?)", searchPattern, searchPattern)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Category").Preload("Categories").Offset(offset).Limit(pageSize).Order("id ASC").Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) UpdateStock(id uint, stock int) error {
	result := r.db.Model(&models.Product{}).Where("id = ?", id).Update("stock", stock)
	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}
	return result.Error
}

func (r *productRepository) Search(query string, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	searchPattern := "%" + query + "%"
	dbQuery := r.db.Model(&models.Product{}).Where(
		"LOWER(name) LIKE LOWER(?) OR LOWER(sku) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)",
		searchPattern, searchPattern, searchPattern,
	)

	dbQuery.Count(&total)

	offset := (page - 1) * pageSize
	err := dbQuery.Preload("Category").Preload("Categories").Offset(offset).Limit(pageSize).Order("id ASC").Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
