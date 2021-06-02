package main

import (
	"net/http"
	"webservice/database"
	"webservice/employee"
	"webservice/product"
	"webservice/receipt"

	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	receipt.SetupRoutes(apiBasePath)
	employee.SetupRoutes(apiBasePath)
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
