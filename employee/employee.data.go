package employee

import (
	"context"
	"database/sql"
	"time"
	"webservice/database"
)

func getEmployeeList() ([]Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(ctx, `SELECT employeeId, name FROM employees`)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	employees := make([]Employee, 0)
	for results.Next() {
		var employee Employee
		results.Scan(&employee.EmployeeID, &employee.Name)
		employees = append(employees, employee)
	}
	return employees, nil
}

func getEmployee(employeeID int) (*Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(ctx, `SELECT employeeId, name FROM employees WHERE employeeId = ?`, employeeID)
	employee := &Employee{}
	err := row.Scan(&employee.EmployeeID, &employee.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return employee, nil
}

func updateEmployee(employee Employee) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(ctx, `UPDATE employees SET name = ? WHERE employeeId = ?`, employee.Name, employee.EmployeeID)
	if err != nil {
		return err
	}
	return nil
}

func insertEmployee(employee Employee) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO employees values ?`)
	if err != nil {
		return 0, nil
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(insertID), nil
}

func removeEmployee(employeeId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM employees WHERE employeeId = ?`, employeeId)
	if err != nil {
		return err
	}
	return nil
}
