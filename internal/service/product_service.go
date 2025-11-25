package service

import (
	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"github.com/brunobarlari/inventorypulse/internal/repository"
)

type ProductService interface {
	Create(req *models.CreateProductRequest) (*models.Product, error)
	GetByID(id uint) (*models.Product, error)
	Update(id uint, req *models.UpdateProductRequest) (*models.Product, error)
	Delete(id uint) error
	List(page, pageSize int, categoryID *uint) ([]models.Product, int64, error)
	UpdateStock(id uint, quantity int) (*models.Product, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) Create(req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Quantity:    req.Quantity,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	// Reload with category
	return s.productRepo.FindByID(product.ID)
}

func (s *productService) GetByID(id uint) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *productService) Update(id uint, req *models.UpdateProductRequest) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.Quantity >= 0 {
		product.Quantity = req.Quantity
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	// Reload with category
	return s.productRepo.FindByID(product.ID)
}

func (s *productService) Delete(id uint) error {
	return s.productRepo.Delete(id)
}

func (s *productService) List(page, pageSize int, categoryID *uint) ([]models.Product, int64, error) {
	return s.productRepo.List(page, pageSize, categoryID)
}

func (s *productService) UpdateStock(id uint, quantity int) (*models.Product, error) {
	if err := s.productRepo.UpdateStock(id, quantity); err != nil {
		return nil, err
	}
	return s.productRepo.FindByID(id)
}

