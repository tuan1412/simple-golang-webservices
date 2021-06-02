package employee

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"webservice/middleware"
)

const employeesPath = "employees"

func handleEmployees(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		employeeList, err := getEmployeeList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		j, err := json.Marshal(employeeList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPost:
		var employee Employee
		err := json.NewDecoder(r.Body).Decode(&employee)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = insertEmployee(employee)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleEmployee(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", employeesPath))
	if len(urlPathSegments[:1]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	employeeID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		product, err := getEmployee(employeeID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if product == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}

	case http.MethodPut:
		var employee Employee
		err := json.NewDecoder(r.Body).Decode(&employee)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if employee.EmployeeID != employeeID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = updateEmployee(employee)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		removeEmployee(employeeID)

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func SetupRoutes(apiBasePath string) {
	productsHandler := http.HandlerFunc(handleEmployees)
	productHandler := http.HandlerFunc(handleEmployee)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, employeesPath), middleware.Cors(productsHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, employeesPath), middleware.Cors(productHandler))
}
