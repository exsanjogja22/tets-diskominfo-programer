package product

type ProductListItem struct {
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	SupplierID   int     `json:"supplier_id"`
	CategoryID   int     `json:"category_id"`
	UnitPrice    float64 `json:"unit_price"`
	UnitsInStock int     `json:"units_in_stock"`
	Discontinued bool    `json:"discontinued"`
	CategoryName string  `json:"category_name"`
	SupplierName string  `json:"supplier_name"`
}

type ProductDetail struct {
	ProductID       int     `json:"product_id"`
	ProductName     string  `json:"product_name"`
	SupplierID      int     `json:"supplier_id"`
	CategoryID      int     `json:"category_id"`
	QuantityPerUnit string  `json:"quantity_per_unit"`
	UnitPrice       float64 `json:"unit_price"`
	UnitsInStock    int     `json:"units_in_stock"`
	UnitsOnOrder    int     `json:"units_on_order"`
	ReorderLevel    int     `json:"reorder_level"`
	Discontinued    bool    `json:"discontinued"`
	CategoryName    string  `json:"category_name"`
	SupplierName    string  `json:"supplier_name"`
	TotalSold       int     `json:"total_sold"`
}

type CategoryItem struct {
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type SupplierItem struct {
	SupplierID   int    `json:"supplier_id"`
	CompanyName  string `json:"company_name"`
}


type Meta struct {
	Pagination PaginationMeta `json:"pagination"`
	Keyword    string        `json:"keyword"`
	Sort       string        `json:"sort"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}


type ListResponse struct {
	Message string           `json:"message"`
	Status  int              `json:"status"`
	Data    []ProductListItem `json:"data"`
	Meta    Meta             `json:"meta"`
}

type DetailResponse struct {
	Message string        `json:"message"`
	Status  int           `json:"status"`
	Data    ProductDetail `json:"data"`
}

type CreateRequest struct {
	ProductName  string  `json:"product_name"`
	SupplierID   int     `json:"supplier_id"`
	CategoryID   int     `json:"category_id"`
	UnitPrice    float64 `json:"unit_price"`
	UnitsInStock *int    `json:"units_in_stock"`
	Discontinued *bool   `json:"discontinued"`
}

type CreateResponse struct {
	Message string        `json:"message"`
	Status  int           `json:"status"`
	Data    ProductDetail `json:"data"`
}
