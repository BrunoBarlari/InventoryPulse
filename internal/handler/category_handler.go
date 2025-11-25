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

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// List godoc
// @Summary      List categories
// @Description  Get paginated list of categories
// @Tags         categories
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        page_size query int false "Page size" default(10)
// @Success      200  {object}  models.PaginatedResponse[models.CategoryResponse]
// @Failure      500  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	var pagination models.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	page := pagination.GetPage()
	pageSize := pagination.GetPageSize()

	categories, total, err := h.categoryService.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve categories",
		})
		return
	}

	// Convert to response DTOs
	responses := make([]models.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = cat.ToResponse()
	}

	c.JSON(http.StatusOK, models.NewPaginatedResponse(responses, page, pageSize, total))
}

// Get godoc
// @Summary      Get category by ID
// @Description  Get a single category by its ID
// @Tags         categories
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200  {object}  models.CategoryResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /categories/{id} [get]
func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid category ID",
		})
		return
	}

	category, err := h.categoryService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "not_found",
				Message: "Category not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve category",
		})
		return
	}

	c.JSON(http.StatusOK, category.ToResponse())
}

// Create godoc
// @Summary      Create category
// @Description  Create a new category (admin only)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        request body models.CreateCategoryRequest true "Category data"
// @Success      201  {object}  models.CategoryResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      409  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	category, err := h.categoryService.Create(&req)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryAlreadyExists) {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "conflict",
				Message: "Category with this name already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, category.ToResponse())
}

// Update godoc
// @Summary      Update category
// @Description  Update an existing category (admin only)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id path int true "Category ID"
// @Param        request body models.UpdateCategoryRequest true "Category data"
// @Success      200  {object}  models.CategoryResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      409  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid category ID",
		})
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	category, err := h.categoryService.Update(uint(id), &req)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "not_found",
				Message: "Category not found",
			})
			return
		}
		if errors.Is(err, repository.ErrCategoryAlreadyExists) {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "conflict",
				Message: "Category with this name already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update category",
		})
		return
	}

	c.JSON(http.StatusOK, category.ToResponse())
}

// Delete godoc
// @Summary      Delete category
// @Description  Delete a category (admin only)
// @Tags         categories
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      409  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid category ID",
		})
		return
	}

	err = h.categoryService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "not_found",
				Message: "Category not found",
			})
			return
		}
		if errors.Is(err, repository.ErrCategoryHasProducts) {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "conflict",
				Message: "Cannot delete category with associated products",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete category",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Category deleted successfully",
	})
}

