package service

import (
	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"github.com/brunobarlari/inventorypulse/internal/repository"
	"github.com/brunobarlari/inventorypulse/pkg/websocket"
)

type CategoryService interface {
	Create(req *models.CreateCategoryRequest) (*models.Category, error)
	GetByID(id uint) (*models.Category, error)
	Update(id uint, req *models.UpdateCategoryRequest) (*models.Category, error)
	Delete(id uint) error
	List(page, pageSize int) ([]models.Category, int64, error)
	Search(query string, page, pageSize int) ([]models.Category, int64, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	wsHub        *websocket.Hub
}

func NewCategoryService(categoryRepo repository.CategoryRepository, wsHub *websocket.Hub) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		wsHub:        wsHub,
	}
}

func (s *categoryService) Create(req *models.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	// Broadcast WebSocket event
	if s.wsHub != nil {
		s.wsHub.BroadcastMessage(websocket.EventCategoryCreated, category.ToResponse())
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

	// Broadcast WebSocket event
	if s.wsHub != nil {
		s.wsHub.BroadcastMessage(websocket.EventCategoryUpdated, category.ToResponse())
	}

	return category, nil
}

func (s *categoryService) Delete(id uint) error {
	// Get category before deleting for the event
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		return err
	}

	// Broadcast WebSocket event
	if s.wsHub != nil {
		s.wsHub.BroadcastMessage(websocket.EventCategoryDeleted, map[string]uint{"id": category.ID})
	}

	return nil
}

func (s *categoryService) List(page, pageSize int) ([]models.Category, int64, error) {
	return s.categoryRepo.List(page, pageSize)
}

func (s *categoryService) Search(query string, page, pageSize int) ([]models.Category, int64, error) {
	return s.categoryRepo.Search(query, page, pageSize)
}
