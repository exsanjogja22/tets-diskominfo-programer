package report

import (
	"database/sql"
	"fmt"
	"strings"
)

type CustomerSale struct {
	CompanyName    string
	Country        string
	TotalPurchases float64
}

func TopCustomersBySales(db *sql.DB) ([]CustomerSale, error) {
	query := "SELECT c.CompanyName, c.Country, " +
		"COALESCE(SUM(od.Quantity * od.UnitPrice), 0) AS Total " +
		"FROM customers c " +
		"LEFT JOIN orders o ON c.CustomerID = o.CustomerID " +
		"LEFT JOIN order_details od ON o.OrderID = od.OrderID " +
		"GROUP BY c.CustomerID, c.CompanyName, c.Country " +
		"ORDER BY Total DESC LIMIT 10"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query top customers: %w", err)
	}
	defer rows.Close()

	var result []CustomerSale
	for rows.Next() {
		var row CustomerSale
		if err := rows.Scan(&row.CompanyName, &row.Country, &row.TotalPurchases); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}
	return result, nil
}

func FormatTopCustomersBySales(data []CustomerSale) string {
	const (
		colCompany = "Company Name"
		colCountry = "Country"
		colTotal   = "Total Penjualan"
	)

	wCompany := len(colCompany)
	wCountry := len(colCountry)
	wTotal := len(colTotal)

	for _, row := range data {
		if len(row.CompanyName) > wCompany {
			wCompany = len(row.CompanyName)
		}
		if len(row.Country) > wCountry {
			wCountry = len(row.Country)
		}
		formatted := formatCurrency(row.TotalPurchases)
		if len(formatted) > wTotal {
			wTotal = len(formatted)
		}
	}

	if wCompany > 22 {
		wCompany = 22
	}
	if wCountry > 15 {
		wCountry = 15
	}
	if wTotal > 18 {
		wTotal = 18
	}

	var b strings.Builder

	b.WriteString("Laporan Penjualan — Total penjualan per pelanggan (nilai transaksi tertinggi):\n")
	b.WriteString("Top 10 Pelanggan:\n")
	sep := "+" + strings.Repeat("-", wCompany+2) + "+" + strings.Repeat("-", wCountry+2) + "+" + strings.Repeat("-", wTotal+2) + "+\n"
	b.WriteString(sep)
	b.WriteString("| " + padRight(colCompany, wCompany) + " | " + padRight(colCountry, wCountry) + " | " + padRight(colTotal, wTotal) + " |\n")
	b.WriteString(sep)

	for _, row := range data {
		company := row.CompanyName
		if len(company) > wCompany {
			company = company[:wCompany-3] + "..."
		}
		country := row.Country
		if len(country) > wCountry {
			country = country[:wCountry-3] + "..."
		}
		b.WriteString("| " + padRight(company, wCompany) + " | " + padRight(country, wCountry) + " | " + padLeft(formatCurrency(row.TotalPurchases), wTotal) + " |\n")
	}
	b.WriteString(sep)
	return b.String()
}

func formatCurrency(v float64) string {
	return fmt.Sprintf("$%s", formatNumber(v))
}

func formatNumber(v float64) string {
	s := fmt.Sprintf("%.2f", v)

	var intPart, decPart string
	if i := strings.Index(s, "."); i >= 0 {
		intPart, decPart = s[:i], s[i:]
	} else {
		intPart = s
		decPart = ""
	}
	n := len(intPart)
	if n <= 3 {
		return s
	}
	var b strings.Builder
	for i, c := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteString(",")
		}
		b.WriteRune(c)
	}
	b.WriteString(decPart)
	return b.String()
}

func padRight(s string, w int) string {
	if len(s) >= w {
		return s[:w]
	}
	return s + strings.Repeat(" ", w-len(s))
}

func padLeft(s string, w int) string {
	if len(s) >= w {
		return s[:w]
	}
	return strings.Repeat(" ", w-len(s)) + s
}
