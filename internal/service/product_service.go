package service

import (
	"time"

	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"github.com/brunobarlari/inventorypulse/internal/repository"
	"github.com/brunobarlari/inventorypulse/pkg/websocket"
)

type ProductService interface {
	Create(req *models.CreateProductRequest) (*models.Product, error)
	GetByID(id uint) (*models.Product, error)
	Update(id uint, req *models.UpdateProductRequest) (*models.Product, error)
	Delete(id uint) error
	List(page, pageSize int, categoryID *uint, search string) ([]models.Product, int64, error)
	UpdateStock(id uint, stock int) (*models.Product, error)
	GetHistory(productID uint, start, end *time.Time, page, pageSize int) ([]models.ProductHistory, int64, error)
}

type productService struct {
	productRepo        repository.ProductRepository
	productHistoryRepo repository.ProductHistoryRepository
	wsHub              *websocket.Hub
}

func NewProductService(productRepo repository.ProductRepository, productHistoryRepo repository.ProductHistoryRepository, wsHub *websocket.Hub) ProductService {
	return &productService{
		productRepo:        productRepo,
		productHistoryRepo: productHistoryRepo,
		wsHub:              wsHub,
	}
}

func (s *productService) Create(req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Stock:       req.Stock,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
	}

	if err := s.productRepo.Create(product, req.CategoryIDs); err != nil {
		return nil, err
	}

	// Reload with category
	product, err := s.productRepo.FindByID(product.ID)
	if err != nil {
		return nil, err
	}

	// Save initial history record
	history := &models.ProductHistory{
		ProductID: product.ID,
		Price:     product.Price,
		Stock:     product.Stock,
		ChangedAt: time.Now(),
	}
	s.productHistoryRepo.Create(history)

	// Broadcast WebSocket event
	if s.wsHub != nil {
		s.wsHub.BroadcastMessage(websocket.EventProductCreated, product.ToResponse())
	}

	return product, nil
}

func (s *productService) GetByID(id uint) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *productService) Update(id uint, req *models.UpdateProductRequest) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Track if price or stock changed for history
	priceChanged := false
	stockChanged := false
	oldPrice := product.Price
	oldStock := product.Stock

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.Stock != nil {
		if *req.Stock != product.Stock {
			stockChanged = true
			product.Stock = *req.Stock
		}
	}
	if req.Price != nil && *req.Price > 0 {
		if *req.Price != product.Price {
			priceChanged = true
			product.Price = *req.Price
		}
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}

	if err := s.productRepo.Update(product, req.CategoryIDs); err != nil {
		return nil, err
	}

	// Reload with category
	product, err = s.productRepo.FindByID(product.ID)
	if err != nil {
		return nil, err
	}

	// Record history if price or stock changed
	if priceChanged || stockChanged {
		history := &models.ProductHistory{
			ProductID: product.ID,
			Price:     product.Price,
			Stock:     product.Stock,
			ChangedAt: time.Now(),
		}
		s.productHistoryRepo.Create(history)

		// Log changes for debugging
		if priceChanged {
			_ = oldPrice // suppress unused variable warning
		}
		if stockChanged {
			_ = oldStock // suppress unused variable warning
		}
	}

	// Broadcast WebSocket event
	if s.wsHub != nil {
		s.wsHub.BroadcastMessage(websocket.EventProductUpdated, product.ToResponse())
	}

	return product, nil
}

func (s *productService) Delete(id uint) error {
	// Get product before deleting for the event
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.productRepo.Delete(id); err != nil {
		return err
	}

	// Broadcast WebSocket event
	if s.wsHub != nil {
		s.wsHub.BroadcastMessage(websocket.EventProductDeleted, map[string]uint{"id": product.ID})
	}

	return nil
}

func (s *productService) List(page, pageSize int, categoryID *uint, search string) ([]models.Product, int64, error) {
	return s.productRepo.List(page, pageSize, categoryID, search)
}

func (s *productService) UpdateStock(id uint, stock int) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	oldStock := product.Stock

	if err := s.productRepo.UpdateStock(id, stock); err != nil {
		return nil, err
	}

	product, err = s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Record history if stock changed
	if oldStock != stock {
		history := &models.ProductHistory{
			ProductID: product.ID,
			Price:     product.Price,
			Stock:     product.Stock,
			ChangedAt: time.Now(),
		}
		s.productHistoryRepo.Create(history)
	}

	// Broadcast WebSocket event
	if s.wsHub != nil {
		s.wsHub.BroadcastMessage(websocket.EventStockUpdated, product.ToResponse())
	}

	return product, nil
}

func (s *productService) GetHistory(productID uint, start, end *time.Time, page, pageSize int) ([]models.ProductHistory, int64, error) {
	// Verify product exists
	_, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, 0, err
	}

	return s.productHistoryRepo.FindByProductID(productID, start, end, page, pageSize)
}
