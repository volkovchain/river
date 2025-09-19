package ethereum

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestMockClient(t *testing.T) {
	t.Run("default behavior", func(t *testing.T) {
		mock := &MockClient{}

		// Test NetworkID
		networkID, err := mock.NetworkID(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(1), networkID)

		// Test PendingNonceAt
		nonce, err := mock.PendingNonceAt(context.Background(), common.Address{})
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), nonce)

		// Test SuggestGasPrice
		gasPrice, err := mock.SuggestGasPrice(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(1000000000), gasPrice)

		// Test SendTransaction
		err = mock.SendTransaction(context.Background(), &types.Transaction{})
		assert.NoError(t, err)

		// Test TransactionReceipt
		receipt, err := mock.TransactionReceipt(context.Background(), common.Hash{})
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), receipt.Status)

		// Test CallContract
		result, err := mock.CallContract(context.Background(), ethereum.CallMsg{}, nil)
		assert.NoError(t, err)
		assert.Empty(t, result)

		// Test EstimateGas
		gas, err := mock.EstimateGas(context.Background(), ethereum.CallMsg{})
		assert.NoError(t, err)
		assert.Equal(t, uint64(100000), gas)
	})

	t.Run("custom behavior", func(t *testing.T) {
		expectedErr := assert.AnError
		mock := &MockClient{
			NetworkIDFn: func(ctx context.Context) (*big.Int, error) {
				return nil, expectedErr
			},
		}

		_, err := mock.NetworkID(context.Background())
		assert.Equal(t, expectedErr, err)
	})

	t.Run("all custom behaviors", func(t *testing.T) {
		expectedErr := assert.AnError
		mock := &MockClient{
			NetworkIDFn: func(ctx context.Context) (*big.Int, error) {
				return big.NewInt(2), nil
			},
			PendingNonceAtFn: func(ctx context.Context, account common.Address) (uint64, error) {
				return 5, nil
			},
			SuggestGasPriceFn: func(ctx context.Context) (*big.Int, error) {
				return big.NewInt(2000000000), nil
			},
			SendTransactionFn: func(ctx context.Context, tx *types.Transaction) error {
				return expectedErr
			},
			TransactionReceiptFn: func(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
				return &types.Receipt{Status: 0}, nil
			},
			CallContractFn: func(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
				return []byte("test"), nil
			},
			EstimateGasFn: func(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
				return 200000, nil
			},
		}

		// Test NetworkID
		networkID, err := mock.NetworkID(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(2), networkID)

		// Test PendingNonceAt
		nonce, err := mock.PendingNonceAt(context.Background(), common.Address{})
		assert.NoError(t, err)
		assert.Equal(t, uint64(5), nonce)

		// Test SuggestGasPrice
		gasPrice, err := mock.SuggestGasPrice(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(2000000000), gasPrice)

		// Test SendTransaction
		err = mock.SendTransaction(context.Background(), &types.Transaction{})
		assert.Equal(t, expectedErr, err)

		// Test TransactionReceipt
		receipt, err := mock.TransactionReceipt(context.Background(), common.Hash{})
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), receipt.Status)

		// Test CallContract
		result, err := mock.CallContract(context.Background(), ethereum.CallMsg{}, nil)
		assert.NoError(t, err)
		assert.Equal(t, []byte("test"), result)

		// Test EstimateGas
		gas, err := mock.EstimateGas(context.Background(), ethereum.CallMsg{})
		assert.NoError(t, err)
		assert.Equal(t, uint64(200000), gas)
	})
}
