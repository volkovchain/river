package db

import (
	"context"
	"database/sql"

	"gitlab.midas.dev/back/river/internal/entity"
	"gitlab.midas.dev/back/river/internal/repository"
)

func NewSalaryRepository(db *sql.DB) repository.SalaryRepository {
	return &salaryRepositorySQLite{db: db}
}

type salaryRepositorySQLite struct {
	db *sql.DB
}

func (s *salaryRepositorySQLite) UpdatePaymentStatusToProcessing(ctx context.Context, id int64) error {
	return s.updatePaymentStatus(ctx, id, repository.ProcessingStatus)
}

func (s *salaryRepositorySQLite) UpdatePaymentStatusToDone(ctx context.Context, id int64) error {
	return s.updatePaymentStatus(ctx, id, repository.DoneStatus)
}

func (s *salaryRepositorySQLite) updatePaymentStatus(ctx context.Context, id int64, status repository.PaymentStatus) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE main.payments SET status = $1 WHERE id = $2`, status, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *salaryRepositorySQLite) Create(ctx context.Context) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO salaries (status) VALUES ($1) RETURNING id;
	`)

	if err != nil {
		return err
	}

	defer func() {
		_ = stmt.Close()
	}()

	var salaryID int
	err = stmt.QueryRowContext(ctx, repository.CreatedStatus).Scan(&salaryID)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO payments (salary_id, employee_id, amount, status, addr) 
		SELECT $1, id, amount_salary, $2, addr FROM employers;
	`, salaryID, repository.CreatedStatus)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *salaryRepositorySQLite) UpdateStatusToProcessing(ctx context.Context, id int64) error {
	return s.updateStatus(ctx, id, repository.ProcessingStatus)
}

func (s *salaryRepositorySQLite) updateStatus(ctx context.Context, id int64, status repository.PaymentStatus) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE salaries SET status = $1 WHERE id = $2`, status, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *salaryRepositorySQLite) UpdateStatusToDone(ctx context.Context, id int64) error {
	return s.updateStatus(ctx, id, repository.DoneStatus)
}

func (s *salaryRepositorySQLite) ListByStatus(ctx context.Context, status repository.PaymentStatus) ([]*entity.Salary, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, status, created_at FROM salaries WHERE status = $1`, status)

	if err != nil {
		return nil, err
	}

	salaries := make([]*entity.Salary, 0)

	for rows.Next() {
		salary := new(entity.Salary)

		err = rows.Scan(&salary.ID, &salary.Status, &salary.CreateAt)
		if err != nil {
			continue
		}

		salaries = append(salaries, salary)
	}

	return salaries, err
}

func (s *salaryRepositorySQLite) ListPaymentsBySalaryID(ctx context.Context, salaryID int64) ([]*entity.Payment, error) {
	rows, err := s.db.QueryContext(ctx, `
	SELECT id, employee_id, salary_id, amount, addr FROM payments WHERE salary_id = $1 AND status != $2
	`, salaryID, repository.DoneStatus)

	if err != nil {
		return nil, err
	}
	payments := make([]*entity.Payment, 0)
	for rows.Next() {
		payment := new(entity.Payment)
		err = rows.Scan(&payment.ID, &payment.EmployeeID, &payment.SalaryID, &payment.Amount, &payment.Addr)
		if err != nil {
			continue
		}

		payments = append(payments, payment)
	}

	return payments, err
}
