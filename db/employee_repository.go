package db

import (
	"context"
	"database/sql"

	"gitlab.midas.dev/back/river/internal/entity"
	"gitlab.midas.dev/back/river/internal/repository"
)

func NewEmployeeRepository(db *sql.DB) repository.EmployeeRepository {
	return &employeeRepositorySQLLite{db: db}
}

type employeeRepositorySQLLite struct {
	db *sql.DB
}

func (e *employeeRepositorySQLLite) List(ctx context.Context) ([]*entity.Employee, error) {
	rows, err := e.db.QueryContext(ctx, `
		SELECT id, name, addr, amount_salary FROM employers;`)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	emps := make([]*entity.Employee, 0)
	for rows.Next() {
		emp := new(entity.Employee)
		err = rows.Scan(&emp.ID, &emp.Name, &emp.Addr, &emp.SalaryAmount)

		if err != nil {
			continue
		}

		emps = append(emps, emp)
	}
	return emps, err
}
