package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Product struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Size     string `db:"size"`
	Code     string `db:"code"`
	Quantity int    `db:"quantity"`
}

type Warehouse struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	IsAvailable bool   `db:"is_available"`
}

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

func main() {
	db, err := sql.Open("postgres", "dbname=your-db-name user=your-db-user password=your-db-password sslmode=require")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := gin.Default()

	router.POST("/products/reserve", func(c *gin.Context) {
		var productCodes []string
		if err := c.ShouldBindJSON(&productCodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ReserveProducts(db, productCodes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Products reserved successfully"})
	})

	router.POST("/products/release", func(c *gin.Context) {
		var productCodes []string
		if err := c.ShouldBindJSON(&productCodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ReleaseProducts(db, productCodes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Products released successfully"})
	})

	router.GET("/products/remaining/:warehouseID", func(c *gin.Context) {
		warehouseIDStr := c.Param("warehouseID")
		warehouseID, err := strconv.Atoi(warehouseIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
			return
		}

		products, err := GetRemainingProducts(db, warehouseID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get remain product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"products": products})
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
