package db_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/flpcastro/fullcycle-repository/hexagonal-repository/adapters/db"
	"github.com/flpcastro/fullcycle-repository/hexagonal-repository/application"
	"github.com/stretchr/testify/require"
)

var Db *sql.DB

func setUp() {
	Db, _ = sql.Open("sqlite3", ":memory:")
	createTable(Db)
	createProduct(Db)
}

func createTable(db *sql.DB) {
	table := `CREATE TABLE products (
		"id" string,
		"name" string,
		"price" float,
		"status" string
	);`

	stmt, err := db.Prepare(table)
	if err != nil {
		log.Fatal(err.Error())
	}

	stmt.Exec()
}

func createProduct(db *sql.DB) {
	insert := `INSERT INTO products values("abc", "Product Test", 0, "disabled")`

	stmt, err := db.Prepare(insert)
	if err != nil {
		log.Fatal(err.Error())
	}

	stmt.Exec()
}

func TestProductDb_Get(t *testing.T) {
	setUp()
	defer Db.Close()
	productDb := db.NewProductDb(Db)
	product, err := productDb.Get("abc")
	require.Nil(t, err)
	require.Equal(t, "Product Test", product.GetName())
	require.Equal(t, 0.0, product.GetPrice())
	require.Equal(t, "disabled", product.GetStatus())
}

func TestProductDb_Save(t *testing.T) {
	setUp()
	defer Db.Close()
	productDb := db.NewProductDb(Db)

	product := application.NewProduct()
	product.Name = "Product Test"
	product.Price = 25

	productCreated, err := productDb.Save(product)
	require.Nil(t, err)
	require.Equal(t, product.Name, productCreated.GetName())
	require.Equal(t, product.Price, productCreated.GetPrice())
	require.Equal(t, product.Status, productCreated.GetStatus())

	product.Status = "enabled"
	productUpdated, err := productDb.Save(product)
	require.Nil(t, err)
	require.Equal(t, product.Name, productUpdated.GetName())
	require.Equal(t, product.Price, productUpdated.GetPrice())
	require.Equal(t, product.Status, productUpdated.GetStatus())
}
