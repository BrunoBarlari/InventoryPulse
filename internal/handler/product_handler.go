package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"github.com/brunobarlari/inventorypulse/internal/repository"
	"github.com/brunobarlari/inventorypulse/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

type ProductListQuery struct {
	models.PaginationRequest
	CategoryID uint `form:"category_id"`
}

// List godoc
// @Summary      List products
// @Description  Get paginated list of products with optional category filter
// @Tags         products
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        page_size query int false "Page size" default(10)
// @Param        category_id query int false "Filter by category ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /products [get]
func (h *ProductHandler) List(c *gin.Context) {
	var query ProductListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	page := query.GetPage()
	pageSize := query.GetPageSize()

	var categoryID *uint
	if query.CategoryID > 0 {
		categoryID = &query.CategoryID
	}

	products, total, err := h.productService.List(page, pageSize, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve products",
		})
		return
	}

	responses := make([]models.ProductResponse, len(products))
	for i, prod := range products {
		responses[i] = prod.ToResponse()
	}

	c.JSON(http.StatusOK, models.NewPaginatedResponse(responses, page, pageSize, total))
}

// Get godoc
// @Summary      Get product by ID
// @Description  Get a single product by its ID
// @Tags         products
// @Produce      json
// @Param        id path int true "Product ID"
// @Success      200  {object}  models.ProductResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /products/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid product ID",
		})
		return
	}

	product, err := h.productService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "not_found",
				Message: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve product",
		})
		return
	}

	c.JSON(http.StatusOK, product.ToResponse())
}

// Create godoc
// @Summary      Create product
// @Description  Create a new product (admin only)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request body models.CreateProductRequest true "Product data"
// @Success      201  {object}  models.ProductResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      409  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	product, err := h.productService.Create(&req)
	if err != nil {
		if errors.Is(err, repository.ErrProductSKUExists) {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "conflict",
				Message: "Product with this SKU already exists",
			})
			return
		}
		if errors.Is(err, repository.ErrInvalidCategory) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid category ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, product.ToResponse())
}

// Update godoc
// @Summary      Update product
// @Description  Update an existing product (admin only)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path int true "Product ID"
// @Param        request body models.UpdateProductRequest true "Product data"
// @Success      200  {object}  models.ProductResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      409  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid product ID",
		})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	product, err := h.productService.Update(uint(id), &req)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "not_found",
				Message: "Product not found",
			})
			return
		}
		if errors.Is(err, repository.ErrProductSKUExists) {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "conflict",
				Message: "Product with this SKU already exists",
			})
			return
		}
		if errors.Is(err, repository.ErrInvalidCategory) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid category ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update product",
		})
		return
	}

	c.JSON(http.StatusOK, product.ToResponse())
}

// Delete godoc
// @Summary      Delete product
// @Description  Delete a product (admin only)
// @Tags         products
// @Produce      json
// @Param        id path int true "Product ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid product ID",
		})
		return
	}

	err = h.productService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "not_found",
				Message: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete product",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// UpdateStock godoc
// @Summary      Update product stock
// @Description  Update the stock quantity of a product (admin only)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path int true "Product ID"
// @Param        request body models.UpdateStockRequest true "Stock data"
// @Success      200  {object}  models.ProductResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /products/{id}/stock [patch]
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid product ID",
		})
		return
	}

	var req models.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	product, err := h.productService.UpdateStock(uint(id), req.Quantity)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "not_found",
				Message: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update stock",
		})
		return
	}

	c.JSON(http.StatusOK, product.ToResponse())
}

