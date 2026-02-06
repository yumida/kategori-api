package repositories

import (
	"cashier-api/models"
	"database/sql"
	"fmt"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	var insertQuery string = "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
	for i := range details {
		details[i].TransactionID = transactionID
		insertQuery += fmt.Sprintf("(%d, %d, %d, %d)", details[i].TransactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if i < len(details)-1 {
			insertQuery += ", "
		}
	}
	_, err = tx.Exec(insertQuery)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetReportToday() (*models.ReportResponse, error) {
	query := `
	SELECT 
		(SELECT COUNT(*) FROM transactions WHERE DATE(created_at) = CURRENT_DATE) AS total_transactions,
		(SELECT SUM(total_amount) FROM transactions WHERE DATE(created_at) = CURRENT_DATE) AS total_revenue`
	var report models.ReportResponse
	err := repo.db.QueryRow(query).Scan(&report.TotalTransactions, &report.TotalRevenue)
	if err != nil {
		return nil, err
	}

	productQuery := `
	SELECT 
		p.id, p.name, SUM(td.quantity) AS total_sold
	FROM
		transaction_details td
	JOIN
		products p ON td.product_id = p.id 
	JOIN
		transactions t ON td.transaction_id = t.id
	WHERE
		DATE(t.created_at) = CURRENT_DATE
	GROUP BY
		p.id, p.name
	ORDER BY
		total_sold DESC
	LIMIT 1`
	var soldProduct models.SoldProductReport
	err = repo.db.QueryRow(productQuery).Scan(&soldProduct.ProductID, &soldProduct.ProductName, &soldProduct.TotalSold)
	if err != nil {
		return nil, err
	}
	report.MostSoldProducts = soldProduct
	return &report, nil
}

func (repo *TransactionRepository) GetReportByDate(StartDate string, EndDate string) (*models.ReportResponse, error) {
	query := `
	SELECT
		(SELECT COUNT(*) FROM transactions WHERE DATE(created_at) BETWEEN $1 AND $2) AS total_transactions,
		(SELECT SUM(total_amount) FROM transactions WHERE DATE(created_at) BETWEEN $1 AND $2) AS total_revenue`
	var report models.ReportResponse
	err := repo.db.QueryRow(query, StartDate, EndDate).Scan(&report.TotalTransactions, &report.TotalRevenue)
	if err != nil {
		return nil, err
	}
	productQuery := `
	SELECT 
		p.id, p.name, SUM(td.quantity) AS total_sold
	FROM
		transaction_details td
	JOIN
		products p ON td.product_id = p.id
	JOIN
		transactions t ON td.transaction_id = t.id
	WHERE	
		DATE(t.created_at) BETWEEN $1 AND $2
	GROUP BY
		p.id, p.name	
	ORDER BY
		total_sold DESC
	LIMIT 1`
	var soldProduct models.SoldProductReport
	err = repo.db.QueryRow(productQuery, StartDate, EndDate).Scan(&soldProduct.ProductID, &soldProduct.ProductName, &soldProduct.TotalSold)
	if err != nil {
		return nil, err
	}
	report.MostSoldProducts = soldProduct
	return &report, nil
}
