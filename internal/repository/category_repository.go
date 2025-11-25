package repository

import (
	"errors"

	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category with this name already exists")
	ErrCategoryHasProducts   = errors.New("category has associated products")
)

type CategoryRepository interface {
	Create(category *models.Category) error
	FindByID(id uint) (*models.Category, error)
	FindByName(name string) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
	List(page, pageSize int) ([]models.Category, int64, error)
	HasProducts(id uint) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	// Check if category with same name already exists
	var count int64
	r.db.Model(&models.Category{}).Where("name = ?", category.Name).Count(&count)
	if count > 0 {
		return ErrCategoryAlreadyExists
	}

	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *models.Category) error {
	// Check if another category with same name exists
	var count int64
	r.db.Model(&models.Category{}).Where("name = ? AND id != ?", category.Name, category.ID).Count(&count)
	if count > 0 {
		return ErrCategoryAlreadyExists
	}

	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	// Check if category has products
	hasProducts, err := r.HasProducts(id)
	if err != nil {
		return err
	}
	if hasProducts {
		return ErrCategoryHasProducts
	}

	result := r.db.Delete(&models.Category{}, id)
	if result.RowsAffected == 0 {
		return ErrCategoryNotFound
	}
	return result.Error
}

func (r *categoryRepository) List(page, pageSize int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	r.db.Model(&models.Category{}).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("id ASC").Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) HasProducts(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("category_id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

