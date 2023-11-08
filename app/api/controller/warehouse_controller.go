package controller

import (
	"database/sql"
	"errors"
)

type Product struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Size     string `json:"size" db:"size"`
	Code     string `json:"code" db:"code"`
	Quantity int    `json:"quantity" db:"quantity"`
}

type Warehouse struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	IsAvailable bool   `json:"is_available" db:"is_available"`
}

// ReserveProducts reserves products
// @Summary Reserves products
// @Description Reserves products and updates their quantities
// @Tags products
// @Accept json
// @Produce json
// @Param productCodes query []string true "Product codes"
// @Success 204 {string} string ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/reserve [post]
// ReserveProducts резервирует продукты
func ReserveProducts(db *sql.DB, productCodes []string) error {
	if len(productCodes) == 0 {
		return errors.New("empty product codes")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateStmt, err := tx.Prepare("UPDATE products SET quantity = quantity - 1 WHERE code = $1")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer updateStmt.Close()

	for _, code := range productCodes {
		var p Product
		err := db.QueryRow("SELECT id, name, size, code, quantity FROM products WHERE code = $1", code).Scan(&p.ID, &p.Name, &p.Size, &p.Code, &p.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}

		if p.Quantity < 1 {
			tx.Rollback()
			return errors.New("product is out of stock")
		}

		_, err = updateStmt.Exec(p.Code)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// ReleaseProducts releases products
// @Summary Releases products
// @Description Releases reserved products and updates their quantities
// @Tags products
// @Accept json
// @Produce json
// @Param productCodes query []string true "Product codes"
// @Success 204 {string} string ""
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/release [post]
// ReleaseProducts отменяет резервирование товаров.
func ReleaseProducts(db *sql.DB, productCodes []string) error {
	if len(productCodes) == 0 {
		return errors.New("empty product codes")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateStmt, err := tx.Prepare("UPDATE products SET quantity = quantity + 1 WHERE code = $1")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer updateStmt.Close()

	for _, code := range productCodes {
		var p Product
		err := db.QueryRow("SELECT id, name, size, code, quantity FROM products WHERE code = $1", code).Scan(&p.ID, &p.Name, &p.Size, &p.Code, &p.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = updateStmt.Exec(p.Code)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetRemainingProducts returns remaining products
// @Summary Returns remaining products
// @Description Returns the remaining products in the warehouse
// @Tags products
// @Accept json
// @Produce json
// @Param warehouseID query int true "Warehouse ID"
// @Success 200 {array} Product
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/remaining [get]
// GetRemainingProducts возвращает оставшееся количество продуктов на складе
func GetRemainingProducts(db *sql.DB, warehouseID int) ([]Product, error) {
	rows, err := db.Query("SELECT code, quantity FROM products WHERE warehouse_id=$1", warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Code, &p.Quantity); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
