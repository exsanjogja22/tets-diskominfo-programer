package product

import (
	"database/sql"
	"fmt"
	"strings"
)

type ListParams struct {
	CategoryID *int
	SupplierID *int
	MinPrice   *float64
	MaxPrice   *float64
	Keyword    string
	Page       int
	Limit      int
	Sort       string
}

func List(db *sql.DB, p ListParams) ([]ProductListItem, int, error) {
	offset := (p.Page - 1) * p.Limit
	if offset < 0 {
		offset = 0
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	var args []interface{}
	var conds []string
	if p.CategoryID != nil {
		conds = append(conds, "p.CategoryID = ?")
		args = append(args, *p.CategoryID)
	}
	if p.SupplierID != nil {
		conds = append(conds, "p.SupplierID = ?")
		args = append(args, *p.SupplierID)
	}
	if p.MinPrice != nil {
		conds = append(conds, "p.UnitPrice >= ?")
		args = append(args, *p.MinPrice)
	}
	if p.MaxPrice != nil {
		conds = append(conds, "p.UnitPrice <= ?")
		args = append(args, *p.MaxPrice)
	}
	if p.Keyword != "" {
		conds = append(conds, "p.ProductName LIKE ?")
		args = append(args, "%"+p.Keyword+"%")
	}
	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	orderBy := "p.ProductID ASC"
	if p.Sort != "" {
		parts := strings.SplitN(p.Sort, ":", 2)
		col := strings.TrimSpace(parts[0])
		dir := "ASC"
		if len(parts) == 2 && strings.ToLower(strings.TrimSpace(parts[1])) == "desc" {
			dir = "DESC"
		}
		switch col {
		case "product_name":
			orderBy = "p.ProductName " + dir
		case "unit_price":
			orderBy = "p.UnitPrice " + dir
		case "units_in_stock":
			orderBy = "p.UnitsInStock " + dir
		}
	}

	countQuery := "SELECT COUNT(*) FROM products p LEFT JOIN categories c ON p.CategoryID = c.CategoryID LEFT JOIN suppliers s ON p.SupplierID = s.SupplierID " + where
	var total int
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	sel := "SELECT p.ProductID, p.ProductName, COALESCE(p.SupplierID,0), COALESCE(p.CategoryID,0), p.UnitPrice, COALESCE(p.UnitsInStock,0), COALESCE(p.Discontinued,0), COALESCE(c.CategoryName,''), COALESCE(s.CompanyName,'') "
	q := sel + " FROM products p LEFT JOIN categories c ON p.CategoryID = c.CategoryID LEFT JOIN suppliers s ON p.SupplierID = s.SupplierID " + where + " ORDER BY " + orderBy + " LIMIT ? OFFSET ?"
	args = append(args, p.Limit, offset)
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	var list []ProductListItem
	for rows.Next() {
		var r ProductListItem
		var disc int
		if err := rows.Scan(&r.ProductID, &r.ProductName, &r.SupplierID, &r.CategoryID, &r.UnitPrice, &r.UnitsInStock, &disc, &r.CategoryName, &r.SupplierName); err != nil {
			return nil, 0, err
		}
		r.Discontinued = disc != 0
		list = append(list, r)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func GetByID(db *sql.DB, id int) (*ProductDetail, error) {
	q := `SELECT p.ProductID, p.ProductName, COALESCE(p.SupplierID,0), COALESCE(p.CategoryID,0),
	COALESCE(p.QuantityPerUnit,''), p.UnitPrice, COALESCE(p.UnitsInStock,0), COALESCE(p.UnitsOnOrder,0), COALESCE(p.ReorderLevel,0), COALESCE(p.Discontinued,0),
	COALESCE(c.CategoryName,''), COALESCE(s.CompanyName,''),
	COALESCE((SELECT SUM(od.Quantity) FROM order_details od WHERE od.ProductID = p.ProductID), 0)
	FROM products p
	LEFT JOIN categories c ON p.CategoryID = c.CategoryID
	LEFT JOIN suppliers s ON p.SupplierID = s.SupplierID
	WHERE p.ProductID = ?`
	var d ProductDetail
	var disc int
	err := db.QueryRow(q, id).Scan(&d.ProductID, &d.ProductName, &d.SupplierID, &d.CategoryID, &d.QuantityPerUnit, &d.UnitPrice, &d.UnitsInStock, &d.UnitsOnOrder, &d.ReorderLevel, &disc, &d.CategoryName, &d.SupplierName, &d.TotalSold)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {

		q2 := `SELECT p.ProductID, p.ProductName, COALESCE(p.SupplierID,0), COALESCE(p.CategoryID,0), p.UnitPrice, COALESCE(p.UnitsInStock,0), COALESCE(p.Discontinued,0),
	COALESCE(c.CategoryName,''), COALESCE(s.CompanyName,''),
	COALESCE((SELECT SUM(od.Quantity) FROM order_details od WHERE od.ProductID = p.ProductID), 0)
	FROM products p
	LEFT JOIN categories c ON p.CategoryID = c.CategoryID
	LEFT JOIN suppliers s ON p.SupplierID = s.SupplierID
	WHERE p.ProductID = ?`
		err2 := db.QueryRow(q2, id).Scan(&d.ProductID, &d.ProductName, &d.SupplierID, &d.CategoryID, &d.UnitPrice, &d.UnitsInStock, &disc, &d.CategoryName, &d.SupplierName, &d.TotalSold)
		if err2 == sql.ErrNoRows {
			return nil, nil
		}
		if err2 != nil {
			return nil, fmt.Errorf("get product: %w", err)
		}
		d.QuantityPerUnit = ""
		d.UnitsOnOrder = 0
		d.ReorderLevel = 0
		d.Discontinued = disc != 0
		return &d, nil
	}
	d.Discontinued = disc != 0
	return &d, nil
}

func Create(db *sql.DB, req CreateRequest) (int64, error) {
	units := 0
	if req.UnitsInStock != nil {
		units = *req.UnitsInStock
	}
	disc := 0
	if req.Discontinued != nil && *req.Discontinued {
		disc = 1
	}
	res, err := db.Exec(
		"INSERT INTO products (ProductName, SupplierID, CategoryID, UnitPrice, UnitsInStock, Discontinued) VALUES (?,?,?,?,?,?)",
		req.ProductName, req.SupplierID, req.CategoryID, req.UnitPrice, units, disc,
	)
	if err != nil {
		return 0, fmt.Errorf("insert product: %w", err)
	}
	return res.LastInsertId()
}

func ExistsSupplier(db *sql.DB, id int) (bool, error) {
	var n int
	err := db.QueryRow("SELECT 1 FROM suppliers WHERE SupplierID = ? LIMIT 1", id).Scan(&n)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExistsCategory(db *sql.DB, id int) (bool, error) {
	var n int
	err := db.QueryRow("SELECT 1 FROM categories WHERE CategoryID = ? LIMIT 1", id).Scan(&n)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func ListCategories(db *sql.DB) ([]CategoryItem, error) {
	rows, err := db.Query("SELECT CategoryID, COALESCE(CategoryName,'') FROM categories ORDER BY CategoryName")
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()
	var list []CategoryItem
	for rows.Next() {
		var c CategoryItem
		if err := rows.Scan(&c.CategoryID, &c.CategoryName); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func ListSuppliers(db *sql.DB) ([]SupplierItem, error) {
	rows, err := db.Query("SELECT SupplierID, COALESCE(CompanyName,'') FROM suppliers ORDER BY CompanyName")
	if err != nil {
		return nil, fmt.Errorf("list suppliers: %w", err)
	}
	defer rows.Close()
	var list []SupplierItem
	for rows.Next() {
		var s SupplierItem
		if err := rows.Scan(&s.SupplierID, &s.CompanyName); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}
