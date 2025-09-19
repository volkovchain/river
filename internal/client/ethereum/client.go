package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Client defines the interface for interacting with Ethereum blockchain
type Client interface {
	// NetworkID retrieves the current network ID
	NetworkID(ctx context.Context) (*big.Int, error)

	// PendingNonceAt retrieves the pending nonce for an account
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)

	// SuggestGasPrice retrieves the currently suggested gas price
	SuggestGasPrice(ctx context.Context) (*big.Int, error)

	// SendTransaction injects a signed transaction into the pending pool for execution
	SendTransaction(ctx context.Context, tx *types.Transaction) error

	// TransactionReceipt returns the receipt of a transaction by transaction hash
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)

	// CallContract executes a message call transaction, which is directly executed in the VM
	// of the node, but never mined into the blockchain
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)

	// EstimateGas tries to estimate the gas needed to execute a specific transaction
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
}
