package controller

import (
	"database/sql"
	"github.com/DmitriiKumancev/lamoda-test/utils"
	"testing"

	_ "github.com/lib/pq"
)

func TestCreateWarehouse(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=root password=secret dbname=lamoda_db sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	w := &Warehouse{
		Name:        utils.RandomString(6),
		IsAvailable: true,
	}

	err = CreateWarehouse(db, w)
	if err != nil {
		t.Fatal(err)
	}

	if w.ID == 0 {
		t.Errorf("Expected warehouse ID to be non-zero, got %d", w.ID)
	}
}

func TestCreateProduct(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=root password=secret dbname=lamoda_db sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	w := &Warehouse{
		Name:        utils.RandomString(6),
		IsAvailable: true,
	}

	err = CreateWarehouse(db, w)
	if err != nil {
		t.Fatal(err)
	}

	p := &Product{
		Name:        utils.RandomString(6),
		Size:        utils.RandomString(6),
		Code:        utils.RandomString(6),
		Quantity:    utils.RandomInt(6),
		WarehouseID: w.ID,
	}

	err = CreateProduct(db, p)
	if err != nil {
		t.Fatal(err)
	}

	if p.ID == 0 {
		t.Errorf("Expected product ID to be non-zero, got %d", p.ID)
	}
}

func TestReserveProductsEmptyProductCodes(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=root password=secret dbname=lamoda_db sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = ReserveProducts(db, []string{})
	if err == nil {
		t.Error("Expected an error with empty product codes, but got nil")
	}
}

func TestReserveProductsInvalidProductCode(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=root password=secret dbname=lamoda_db sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	w := &Warehouse{
		Name:        utils.RandomString(6),
		IsAvailable: true,
	}
	err = CreateWarehouse(db, w)
	if err != nil {
		t.Fatal(err)
	}

	p := &Product{
		Name:        utils.RandomString(6),
		Size:        utils.RandomString(6),
		Code:        utils.RandomString(6),
		Quantity:    1,
		WarehouseID: w.ID,
	}
	err = CreateProduct(db, p)
	if err != nil {
		t.Fatal(err)
	}

	err = ReserveProducts(db, []string{"invalid-code"})
	if err == nil {
		t.Error("Expected an error with invalid product code, but got nil")
	}
}

func TestReserveProductsProductOutOfStock(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=root password=secret dbname=lamoda_db sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	w := &Warehouse{
		Name:        utils.RandomString(6),
		IsAvailable: true,
	}
	err = CreateWarehouse(db, w)
	if err != nil {
		t.Fatal(err)
	}

	p := &Product{
		Name:        utils.RandomString(6),
		Size:        utils.RandomString(6),
		Code:        utils.RandomString(6),
		Quantity:    0, // устанавливаем количество 0, чтобы продукт был недоступен для бронирования
		WarehouseID: w.ID,
	}
	err = CreateProduct(db, p)
	if err != nil {
		t.Fatal(err)
	}

	err = ReserveProducts(db, []string{p.Code})
	if err == nil {
		t.Errorf("Expected error, but got nil")
	} else if err.Error() != "product is out of stock" {
		t.Errorf("Expected error message 'product is out of stock', but got '%s'", err.Error())
	}
}
