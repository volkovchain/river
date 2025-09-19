package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// MockClient is a mock implementation of the Client interface for testing
type MockClient struct {
	NetworkIDFn          func(ctx context.Context) (*big.Int, error)
	PendingNonceAtFn     func(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPriceFn    func(ctx context.Context) (*big.Int, error)
	SendTransactionFn    func(ctx context.Context, tx *types.Transaction) error
	TransactionReceiptFn func(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	CallContractFn       func(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	EstimateGasFn        func(ctx context.Context, call ethereum.CallMsg) (uint64, error)
}

func (m *MockClient) NetworkID(ctx context.Context) (*big.Int, error) {
	if m.NetworkIDFn != nil {
		return m.NetworkIDFn(ctx)
	}
	return big.NewInt(1), nil
}

func (m *MockClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if m.PendingNonceAtFn != nil {
		return m.PendingNonceAtFn(ctx, account)
	}
	return 0, nil
}

func (m *MockClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if m.SuggestGasPriceFn != nil {
		return m.SuggestGasPriceFn(ctx)
	}
	return big.NewInt(1000000000), nil
}

func (m *MockClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if m.SendTransactionFn != nil {
		return m.SendTransactionFn(ctx, tx)
	}
	return nil
}

func (m *MockClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	if m.TransactionReceiptFn != nil {
		return m.TransactionReceiptFn(ctx, txHash)
	}
	return &types.Receipt{Status: 1}, nil
}

func (m *MockClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if m.CallContractFn != nil {
		return m.CallContractFn(ctx, call, blockNumber)
	}
	return []byte{}, nil
}

func (m *MockClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	if m.EstimateGasFn != nil {
		return m.EstimateGasFn(ctx, call)
	}
	return 100000, nil
}
