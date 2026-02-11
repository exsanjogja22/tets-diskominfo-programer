package product

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) ListProducts(c echo.Context) error {
	params := ListParams{
		Page:  1,
		Limit: 10,
	}
	if v := c.QueryParam("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			params.Page = n
		}
	}
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			params.Limit = n
		}
	}
	if v := c.QueryParam("category_id"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.CategoryID = &n
		}
	}
	if v := c.QueryParam("supplier_id"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.SupplierID = &n
		}
	}
	if v := c.QueryParam("min_price"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			params.MinPrice = &f
		}
	}
	if v := c.QueryParam("max_price"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			params.MaxPrice = &f
		}
	}
	if v := c.QueryParam("keyword"); v != "" {
		params.Keyword = v
	}
	if v := c.QueryParam("search"); v != "" {
		params.Keyword = v
	}
	if v := c.QueryParam("product_name"); v != "" {
		params.Keyword = v
	}

	if v := c.QueryParam("sort"); v != "" {
		params.Sort = v
	}

	list, total, err := List(h.DB, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to get products",
			"status":  http.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	totalPages := total / params.Limit
	if total%params.Limit > 0 {
		totalPages++
	}
	if total == 0 {
		totalPages = 0
	}

	return c.JSON(http.StatusOK, ListResponse{
		Message: "Success",
		Status:  200,
		Data:    list,
		Meta: Meta{
			Pagination: PaginationMeta{
				Page:       params.Page,
				Limit:      params.Limit,
				Total:      total,
				TotalPages: totalPages,
			},
			Keyword: params.Keyword,
			Sort:    params.Sort,
		},
	})
}

// GetProductByID GET /api/products/:id
func (h *Handler) GetProductByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid product ID",
			"status":  http.StatusBadRequest,
		})
	}

	detail, err := GetByID(h.DB, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to get product",
			"status":  http.StatusInternalServerError,
			"error":   err.Error(),
		})
	}
	if detail == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Product not found",
			"status":  http.StatusNotFound,
		})
	}

	return c.JSON(http.StatusOK, DetailResponse{
		Message: "Success",
		Status:  200,
		Data:    *detail,
	})
}

func (h *Handler) GetCategories(c echo.Context) error {
	list, err := ListCategories(h.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to get categories",
			"status":  http.StatusInternalServerError,
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success",
		"status":  200,
		"data":    list,
	})
}

func (h *Handler) GetSuppliers(c echo.Context) error {
	list, err := ListSuppliers(h.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to get suppliers",
			"status":  http.StatusInternalServerError,
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success",
		"status":  200,
		"data":    list,
	})
}

func (h *Handler) CreateProduct(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
			"status":  http.StatusBadRequest,
		})
	}

	if len(req.ProductName) < 3 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "product_name is required and must be at least 3 characters",
			"status":  http.StatusBadRequest,
		})
	}

	if req.SupplierID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "supplier_id is required",
			"status":  http.StatusBadRequest,
		})
	}

	if req.CategoryID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "category_id is required",
			"status":  http.StatusBadRequest,
		})
	}

	ok, err := ExistsSupplier(h.DB, req.SupplierID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to validate supplier",
			"status":  http.StatusInternalServerError,
		})
	}
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "supplier_id does not exist",
			"status":  http.StatusBadRequest,
		})
	}

	ok, err = ExistsCategory(h.DB, req.CategoryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to validate category",
			"status":  http.StatusInternalServerError,
		})
	}
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "category_id does not exist",
			"status":  http.StatusBadRequest,
		})
	}

	if req.UnitPrice <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "unit_price is required and must be greater than 0",
			"status":  http.StatusBadRequest,
		})
	}

	newID, err := Create(h.DB, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to create product",
			"status":  http.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	detail, err := GetByID(h.DB, int(newID))
	if err != nil || detail == nil {
		return c.JSON(http.StatusOK, CreateResponse{
			Message: "Success",
			Status:  201,
			Data: ProductDetail{
				ProductID:       int(newID),
				ProductName:     req.ProductName,
				SupplierID:      req.SupplierID,
				CategoryID:      req.CategoryID,
				QuantityPerUnit: "",
				UnitPrice:       req.UnitPrice,
				UnitsInStock:    0,
				UnitsOnOrder:    0,
				ReorderLevel:    0,
				Discontinued:    false,
				CategoryName:    "",
				SupplierName:    "",
				TotalSold:       0,
			},
		})
	}

	return c.JSON(http.StatusCreated, CreateResponse{
		Message: "Success",
		Status:  201,
		Data:    *detail,
	})
}
