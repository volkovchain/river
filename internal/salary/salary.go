package salary

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.midas.dev/back/river/internal/entity"
	"gitlab.midas.dev/back/river/internal/payment"
	"gitlab.midas.dev/back/river/internal/repository"
)

type Service struct {
	salaryRepository repository.SalaryRepository
	paymentService   *payment.Service
}

func New(salary repository.SalaryRepository, paymentService *payment.Service) *Service {
	return &Service{paymentService: paymentService, salaryRepository: salary}
}

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

func (s *Service) startPay(ctx context.Context) error {
	err := s.salaryRepository.Create(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) pay(ctx context.Context, salaries []*entity.Salary) error {
	log.Println("start pay")
	log.Printf("%d salaries\n", len(salaries))

	for _, salary := range salaries {
		log.Printf("salary id - %d\n", salary.ID)
		payments, err := s.salaryRepository.ListPaymentsBySalaryID(ctx, salary.ID)
		if err != nil {
			log.Println(err)
			return err
		}

		err = s.salaryRepository.UpdateStatusToProcessing(ctx, salary.ID)

		if err != nil {
			log.Println(err)
			return err
		}
		countPayments := len(payments)

		// Time wait - range 30 minute
		var tWait time.Duration
		if tWait > 1 {
			tWait = time.Duration(math.Round(float64((30 * 60) / countPayments)))
		}
		log.Printf("wait - %d seconds", tWait)
		var countErrPayments int
		for _, paymt := range payments {
			if paymt.Status == string(repository.CreatedStatus) {
				err := s.salaryRepository.UpdatePaymentStatusToProcessing(ctx, paymt.ID)
				if err != nil {
					countErrPayments++
					log.Println(err)
					return err
				}
			}

			if paymt.Status == string(repository.CreatedStatus) {
				if err != nil {
					countErrPayments++
					return err
				}
			}

			if !common.IsHexAddress(paymt.Addr) {
				countErrPayments++
				log.Printf("in not hex address - %s", paymt.Addr)
				continue
			}

			addr := common.HexToAddress(paymt.Addr)
			err = s.paymentService.Send(ctx, addr, paymt.Amount)

			if err == nil {

				err := s.salaryRepository.UpdatePaymentStatusToDone(ctx, paymt.ID)
				if err != nil {
					countErrPayments++
					log.Println(err)
					return err
				}
			} else {
				countErrPayments++
			}
			time.Sleep(tWait * time.Second)
		}
		if countErrPayments == 0 {

			err = s.salaryRepository.UpdateStatusToDone(ctx, salary.ID)

			if err != nil {
				log.Println(err)
				return err
			}
		}

	}
	return nil
}
