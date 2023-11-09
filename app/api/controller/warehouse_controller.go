package controller

//	@title			Warehouse API Documentation
//	@description	This is a sample API for a warehouse application
//	@version		1
//	@host			localhost:8080

import (
	"database/sql"
	"errors"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Size        string `json:"size"`
	Code        string `json:"code"`
	Quantity    int    `json:"quantity"`
	WarehouseID int    `json:"warehouse_id"`
}

type Warehouse struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	IsAvailable bool   `json:"is_available" db:"is_available"`
}

//	@Summary		Create a new warehouse.
//	@Description	Create a new warehouse in the database.
//	@Tags			warehouses
//	@Accept			json
//	@Produce		json
//	@Param			warehouse	body		Warehouse		true	"Warehouse information"
//	@Success		200			{string}	string			"Warehouse created"
//	@Failure		400			{object}	ErrorResponse	"Invalid request format"
//	@Failure		500			{object}	ErrorResponse	"Internal server error"
//	@Router			/create-warehouse [post]
//
func CreateWarehouse(db *sql.DB, w *Warehouse) error {
	stmt, err := db.Prepare("INSERT INTO warehouse(name, is_available) VALUES($1, $2) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(w.Name, w.IsAvailable).Scan(&w.ID)
	if err != nil {
		return err
	}

	return nil
}

//	@Summary		Create a new product.
//	@Description	Create a new product on a specified warehouse.
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			product	body		Product			true	"Product information"
//	@Success		200		{string}	string			"Product created"
//	@Failure		400		{object}	ErrorResponse	"Invalid request format"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/create-product [post]
//
func CreateProduct(db *sql.DB, p *Product) error {
	stmt, err := db.Prepare("INSERT INTO products(name, size, code, quantity, warehouse_id) VALUES($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(p.Name, p.Size, p.Code, p.Quantity, p.WarehouseID).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

//	@Summary		Delete a product
//	@Description	Delete a product by its ID.
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"Product ID"
//	@Success		200	{string}	string			"Product deleted successfully"
//	@Failure		400	{object}	ErrorResponse	"Invalid request format"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/delete-product/:id [delete]
//
func DeleteProduct(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

//	@Summary		Reserves products
//	@Description	Reserves products and updates their quantities
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			productCodes	query		[]string	true	"Product codes"
//	@Success		204				{string}	string		""
//	@Failure		400				{object}	ErrorResponse
//	@Failure		500				{object}	ErrorResponse
//	@Router			/reserve-products [post]
//
func ReserveProducts(db *sql.DB, productCodes []string) error {
	if len(productCodes) == 0 {
		return errors.New("empty product codes")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, code := range productCodes {
		row := tx.QueryRow("SELECT id, name, size, code, quantity FROM products WHERE code = $1 FOR UPDATE", code)

		var p Product
		err := row.Scan(&p.ID, &p.Name, &p.Size, &p.Code, &p.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}

		if p.Quantity < 1 {
			tx.Rollback()
			return errors.New("product is out of stock")
		}

		_, err = tx.Exec("UPDATE products SET quantity = quantity - 1 WHERE id = $1", p.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//	@Summary		Releases products
//	@Description	Releases reserved products and updates their quantities
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			productCodes	query		[]string	true	"Product codes"
//	@Success		204				{string}	string		""
//	@Failure		400				{object}	ErrorResponse
//	@Failure		500				{object}	ErrorResponse
//	@Router			/release-products [post]
//
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

// @Description Get remaining products for a given warehouse.
// @Tags products
// @Accept json
// @Produce json
// @Param warehouseID path int true "Warehouse ID"
// @Success 200 {array} Product "Remaining products"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /remaining-products/{warehouseID} [get]
// GetRemainingProducts возвращает оставшееся количество продуктов на складе
func GetRemainingProducts(db *sql.DB, warehouseID int) ([]Product, error) {
	rows, err := db.Query("SELECT code, quantity FROM products WHERE warehouse_id = $1", warehouseID)
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
		p.WarehouseID = warehouseID
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
