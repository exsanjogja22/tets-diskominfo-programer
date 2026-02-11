package main

import (
	"log"
	"net/http"

	"github.com/exsanjogja22/test-pemrograman-go.git/api/product"
	"github.com/exsanjogja22/test-pemrograman-go.git/database"
	"github.com/exsanjogja22/test-pemrograman-go.git/report"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan env dari sistem")
	}

	db := database.InitialDb()

	err := db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Tampilkan dan println laporan ke konsol saat startup
	data, err := report.TopCustomersBySales(db)
	if err != nil {
		log.Println("Laporan penjualan:", err)
	} else {
		text := report.FormatTopCustomersBySales(data)
		log.Println("\n" + text)
	}

	e := echo.New()

	e.GET("/report/sales/top-customers", func(c echo.Context) error {
		data, err := report.TopCustomersBySales(db)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error: "+err.Error())
		}
		text := report.FormatTopCustomersBySales(data)
		return c.String(http.StatusOK, text)
	})

	productHandler := product.NewHandler(db)
	e.GET("/api/products", productHandler.ListProducts)
	e.GET("/api/products/:id", productHandler.GetProductByID)
	e.POST("/api/products", productHandler.CreateProduct)
	e.GET("/api/categories", productHandler.GetCategories)
	e.GET("/api/suppliers", productHandler.GetSuppliers)

	e.Start(":8083")
}
