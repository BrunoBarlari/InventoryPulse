package handler

import (
	"net/http"

	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"github.com/brunobarlari/inventorypulse/internal/service"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	productService  service.ProductService
	categoryService service.CategoryService
}

func NewSearchHandler(productService service.ProductService, categoryService service.CategoryService) *SearchHandler {
	return &SearchHandler{
		productService:  productService,
		categoryService: categoryService,
	}
}

// SearchQuery represents the search query parameters
type SearchQuery struct {
	models.PaginationRequest
	Query string `form:"q" binding:"required"`
	Type  string `form:"type"` // "product", "category", or empty for both
}

// SearchResult represents the unified search result
type SearchResult struct {
	Products   interface{} `json:"products,omitempty"`
	Categories interface{} `json:"categories,omitempty"`
}

// Search godoc
// @Summary      Search products and categories
// @Description  Unified search endpoint for products and/or categories
// @Tags         search
// @Produce      json
// @Param        q query string true "Search query"
// @Param        type query string false "Type to search: 'product', 'category', or empty for both"
// @Param        page query int false "Page number" default(1)
// @Param        page_size query int false "Page size" default(10)
// @Success      200  {object}  SearchResult
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Security     BearerAuth
// @Router       /search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	var query SearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Query parameter 'q' is required",
		})
		return
	}

	page := query.GetPage()
	pageSize := query.GetPageSize()

	result := SearchResult{}

	// Search products if type is empty or "product"
	if query.Type == "" || query.Type == "product" {
		products, total, err := h.productService.List(page, pageSize, nil, query.Query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to search products",
			})
			return
		}

		productResponses := make([]models.ProductResponse, len(products))
		for i, p := range products {
			productResponses[i] = p.ToResponse()
		}

		result.Products = models.NewPaginatedResponse(productResponses, page, pageSize, total)
	}

	// Search categories if type is empty or "category"
	if query.Type == "" || query.Type == "category" {
		categories, total, err := h.categoryService.Search(query.Query, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to search categories",
			})
			return
		}

		categoryResponses := make([]models.CategoryResponse, len(categories))
		for i, cat := range categories {
			categoryResponses[i] = cat.ToResponse()
		}

		result.Categories = models.NewPaginatedResponse(categoryResponses, page, pageSize, total)
	}

	c.JSON(http.StatusOK, result)
}

