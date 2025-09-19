package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ClientImpl implements the Client interface using the go-ethereum client
type ClientImpl struct {
	client *ethclient.Client
}

// NewClient creates a new Ethereum client
func NewClient(rawClient *ethclient.Client) Client {
	return &ClientImpl{client: rawClient}
}

// NetworkID retrieves the current network ID
func (c *ClientImpl) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.client.NetworkID(ctx)
}

// PendingNonceAt retrieves the pending nonce for an account
func (c *ClientImpl) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return c.client.PendingNonceAt(ctx, account)
}

// SuggestGasPrice retrieves the currently suggested gas price
func (c *ClientImpl) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.client.SuggestGasPrice(ctx)
}

// SendTransaction injects a signed transaction into the pending pool for execution
func (c *ClientImpl) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.client.SendTransaction(ctx, tx)
}

// TransactionReceipt returns the receipt of a transaction by transaction hash
func (c *ClientImpl) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return c.client.TransactionReceipt(ctx, txHash)
}

// CallContract executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain
func (c *ClientImpl) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return c.client.CallContract(ctx, call, blockNumber)
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction
func (c *ClientImpl) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	return c.client.EstimateGas(ctx, call)
}
