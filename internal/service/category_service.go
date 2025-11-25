package service

import (
	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"github.com/brunobarlari/inventorypulse/internal/repository"
)

type CategoryService interface {
	Create(req *models.CreateCategoryRequest) (*models.Category, error)
	GetByID(id uint) (*models.Category, error)
	Update(id uint, req *models.UpdateCategoryRequest) (*models.Category, error)
	Delete(id uint) error
	List(page, pageSize int) ([]models.Category, int64, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) Create(req *models.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *categoryService) Update(id uint, req *models.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(id uint) error {
	return s.categoryRepo.Delete(id)
}

func (s *categoryService) List(page, pageSize int) ([]models.Category, int64, error) {
	return s.categoryRepo.List(page, pageSize)
}

