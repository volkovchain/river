package payment

import (
	"context"
	"math/big"

	"gitlab.midas.dev/back/river/internal/client/ethereum"
	"gitlab.midas.dev/back/river/internal/types"
)

// Service handles payment-related business logic
type Service struct {
	client ethereum.Client
	// Add other fields as needed
}

// New creates a new payment service
func New(client ethereum.Client) *Service {
	return &Service{
		client: client,
	}
}

// Send transfers tokens to the specified address
func (s *Service) Send(ctx context.Context, to types.Address, valueAmount int64) error {
	// Implementation will be migrated from the existing payment package
	return nil
}

// Balance represents a token balance
type Balance struct {
	Amount *big.Int
}
