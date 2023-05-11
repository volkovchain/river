package repository

import (
	"context"

	"gitlab.midas.dev/back/river/internal/entity"
)

type PaymentStatus string

const (
	CreatedStatus    PaymentStatus = "created"
	ProcessingStatus               = "processing"
	DoneStatus                     = "done"
)

type EmployeeRepository interface {
	List(ctx context.Context) ([]*entity.Employee, error)
}

type SalaryRepository interface {
	Create(ctx context.Context) error
	UpdateStatusToProcessing(ctx context.Context, id int64) error
	UpdateStatusToDone(ctx context.Context, id int64) error
	UpdatePaymentStatusToProcessing(ctx context.Context, id int64) error
	UpdatePaymentStatusToDone(ctx context.Context, id int64) error
	ListByStatus(ctx context.Context, status PaymentStatus) ([]*entity.Salary, error)
	ListPaymentsBySalaryID(ctx context.Context, salaryID int64) ([]*entity.Payment, error)
}
