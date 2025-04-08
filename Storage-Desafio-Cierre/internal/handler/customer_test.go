package handler_test

import (
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/service"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func init() {
	dsn := mysql.Config{
		User:   "root",
		Passwd: "root",
		Addr:   "127.0.0.1:3306",
		Net:    "tcp",
		DBName: "fantasy_products_test",
	}
	// Register
	txdb.Register("txdb", "mysql", dsn.FormatDSN())
}

// reset db for new test (R)
func ResetDB(db *sql.DB) error {
	queries := []string{
		// clear all tables
		"DELETE FROM invoices;",
		"DELETE FROM customers;",
		// reset all auto increment
		"ALTER TABLE invoices AUTO_INCREMENT = 1;",
		"ALTER TABLE customers AUTO_INCREMENT = 1;",
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

// Insert all data necessary
func SetupTestData(db *sql.DB) error {
	customersSQL := `
		INSERT INTO customers (id, first_name, last_name, ` + "`condition`" + `) VALUES
		(1, 'Michael', 'Jordan', 1),
		(2, 'Sarah', 'Connor', 1),
		(3, 'Albert', 'Einstein', 1),
		(4, 'Isaac', 'Newton', 1),
		(5, 'Marie', 'Curie', 1),
		(6, 'Nikola', 'Tesla', 1),
		(7, 'Ada', 'Lovelace', 0),
		(8, 'Alan', 'Turing', 0);
	`
	if _, err := db.Exec(customersSQL); err != nil {
		return err
	}

	invoicesSQL := `
		INSERT INTO invoices (id, customer_id, total) VALUES
		(1, 1, 1200),
		(2, 2, 650),
		(3, 3, 300),
		(4, 4, 150),
		(5, 5, 75),
		(6, 6, 40),
		(7, 7, 15),
		(8, 8, 8);
	`
	if _, err := db.Exec(invoicesSQL); err != nil {
		return err
	}

	return nil
}

func TestGetTopActiveCustomers(t *testing.T) {
	t.Run("should return top active customers based on amount spent", func(t *testing.T) {
		// Open connection db
		db, err := sql.Open("txdb", "")
		assert.NoError(t, err)
		defer db.Close()

		// Config db test
		assert.NoError(t, SetupTestData(db))

		// inject dependency
		repo := repository.NewCustomersMySQL(db)
		service := service.NewCustomersDefault(repo)
		handler := handler.NewCustomersDefault(service)

		handlerFunc := handler.GetTopActiveCustomersByAmountSpent()
		request := httptest.NewRequest(http.MethodGet, "/customers/top-active-customers-by-amount-spent", nil)
		response := httptest.NewRecorder()

		handlerFunc(response, request)

		expectedCode := http.StatusOK
		expectedResponse := `
			{
				"message": "customers found",
				"data": [
					{
						"first_name": "Michael",
						"last_name": "Jordan",
						"total": 1200
					},
					{
						"first_name": "Sarah",
						"last_name": "Connor",
						"total": 650
					},
					{
						"first_name": "Albert",
						"last_name": "Einstein",
						"total": 300
					},
					{
						"first_name": "Isaac",
						"last_name": "Newton",
						"total": 150
					},
					{
						"first_name": "Marie",
						"last_name": "Curie",
						"total": 75
					}
				]
			}
		`

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponse, response.Body.String())

		// Reset DB
		assert.NoError(t, ResetDB(db))
	})

	t.Run("should return empty list when no customers match", func(t *testing.T) {
		// Open db
		db, err := sql.Open("txdb", "")
		assert.NoError(t, err)
		defer db.Close()

		// Reset DB
		assert.NoError(t, ResetDB(db))

		// inject
		repo := repository.NewCustomersMySQL(db)
		service := service.NewCustomersDefault(repo)
		handler := handler.NewCustomersDefault(service)

		handlerFunc := handler.GetTopActiveCustomersByAmountSpent()
		request := httptest.NewRequest(http.MethodGet, "/customers/top-active-customers-by-amount-spent", nil)
		response := httptest.NewRecorder()

		// Exec
		handlerFunc(response, request)

		expectedCode := http.StatusOK
		expectedResponse := `{"message": "customers found", "data": []}`

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponse, response.Body.String())

		// Reset DB
		assert.NoError(t, ResetDB(db))
	})
}

func TestGetInvoicesByCondition(t *testing.T) {
	t.Run("should return invoices grouped by condition", func(t *testing.T) {
		db, err := sql.Open("txdb", "")
		assert.NoError(t, err)
		defer db.Close()

		assert.NoError(t, SetupTestData(db))

		repo := repository.NewCustomersMySQL(db)
		service := service.NewCustomersDefault(repo)
		handler := handler.NewCustomersDefault(service)

		handlerFunc := handler.GetInvoicesByCondition()
		request := httptest.NewRequest(http.MethodGet, "/customers/invoices-by-condition", nil)
		response := httptest.NewRecorder()

		handlerFunc(response, request)

		expectedCode := http.StatusOK
		expectedResponse := `
			{
				"message": "customers found",
				"data": [
					{
						"condition": 1,
						"total": 2300
					},
					{
						"condition": 0,
						"total": 23
					}
				]
			}
		`

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponse, response.Body.String())

		assert.NoError(t, ResetDB(db))
	})

	t.Run("should return empty result if no invoices exist", func(t *testing.T) {
		db, err := sql.Open("txdb", "")
		assert.NoError(t, err)
		defer db.Close()

		assert.NoError(t, ResetDB(db))

		repo := repository.NewCustomersMySQL(db)
		service := service.NewCustomersDefault(repo)
		handler := handler.NewCustomersDefault(service)

		handlerFunc := handler.GetInvoicesByCondition()
		request := httptest.NewRequest(http.MethodGet, "/customers/invoices-by-condition", nil)
		response := httptest.NewRecorder()

		handlerFunc(response, request)

		expectedCode := http.StatusOK
		expectedResponse := `{"message": "customers found", "data": []}`

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponse, response.Body.String())

		assert.NoError(t, ResetDB(db))
	})
}
