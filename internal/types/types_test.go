package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// mockAddress is a simple implementation of the Address interface for testing
type mockAddress struct {
	addr string
}

func (m mockAddress) String() string {
	return m.addr
}

func TestAddressInterface(t *testing.T) {
	// Test that common.Address implements the Address interface
	var _ Address = common.Address{}

	// Test that our mock implements the Address interface
	var _ Address = mockAddress{}

	// Test the mock address
	mockAddr := mockAddress{addr: "test"}
	assert.Equal(t, "test", mockAddr.String())
}
