package salary

import (
	"context"

	"gitlab.midas.dev/back/river/internal/entity"
	"gitlab.midas.dev/back/river/internal/repository"
	"gitlab.midas.dev/back/river/internal/types"
)

// Service handles salary-related business logic
type Service struct {
	salaryRepository repository.SalaryRepository
	paymentService   PaymentService
}

// PaymentService defines the interface for payment operations
type PaymentService interface {
	Send(ctx context.Context, to types.Address, valueAmount int64) error
}

// New creates a new salary service
func New(salaryRepository repository.SalaryRepository, paymentService PaymentService) *Service {
	return &Service{
		salaryRepository: salaryRepository,
		paymentService:   paymentService,
	}
}

// Repay processes salaries that are in processing status
func (s *Service) Repay(ctx context.Context) error {
	salaries, err := s.salaryRepository.ListByStatus(ctx, repository.ProcessingStatus)
	if err != nil {
		return err
	}

	if len(salaries) > 0 {
		err = s.pay(ctx, salaries)
		if err != nil {
			return err
		}
	}
	return nil
}

// Pay creates a new salary and processes it
func (s *Service) Pay(ctx context.Context) error {
	err := s.startPay(ctx)
	if err != nil {
		return err
	}

	salaries, err := s.salaryRepository.ListByStatus(ctx, repository.CreatedStatus)
	if err != nil {
		return err
	}

	if len(salaries) > 0 {
		err = s.pay(ctx, salaries)
		if err != nil {
			return err
		}
	}

	return nil
}

// startPay creates a new salary record
func (s *Service) startPay(ctx context.Context) error {
	err := s.salaryRepository.Create(ctx)
	if err != nil {
		return err
	}
	return nil
}

// pay processes the actual payment for salaries
func (s *Service) pay(ctx context.Context, salaries []*entity.Salary) error {
	// Implementation will be migrated from the existing salary package
	return nil
}
