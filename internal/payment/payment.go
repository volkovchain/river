package payment

import (
	"context"
	"crypto/ecdsa"
	_ "embed"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

//go:embed abi/erc20.json
var erc20abi string

const (
	USDCContractAddress = "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	MethodErc20Balance  = "balanceOf"
)

type Service struct {
	abi         abi.ABI
	privateKeys []*ecdsa.PrivateKey
	client      *ethclient.Client
}

func New(node string, privateKeys []string) *Service {
	ab, err := abi.JSON(strings.NewReader(erc20abi))
	if err != nil {
		panic(err)
	}

	if len(privateKeys) == 0 {
		log.Fatal("")
	}

	s := &Service{}

	s.abi = ab

	for _, key := range privateKeys {
		pk, err := crypto.HexToECDSA(key)

		if err != nil {
			log.Fatal(err)
		}

		s.privateKeys = append(s.privateKeys, pk)
	}

	client, err := ethclient.Dial(node)

	if err != nil {
		log.Fatal(err)
	}

	s.client = client
	return s
}

func (s *Service) Send(ctx context.Context, to common.Address, valueAmount int64) error {
	var (
		pk          *ecdsa.PrivateKey
		fromAddress common.Address
	)

	tokenAddress := common.HexToAddress(USDCContractAddress)
	for _, pk = range s.privateKeys {
		publicKeyECDSA, ok := pk.Public().(*ecdsa.PublicKey)

		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}

		fromAddress = crypto.PubkeyToAddress(*publicKeyECDSA)

		balance, err := s.FetchTokenBalance(tokenAddress, fromAddress)

		if err != nil {
			log.Printf("fetch balance error - %v", err)
			continue
		}
		fmt.Printf("employee amount - %d, wallet balance - %d\n", valueAmount, balance.Int64())
		if valueAmount < balance.Int64() {
			log.Println("balance normal")
			break
		} else {
			log.Println("balance - insufficient funds")
		}
	}

	nonce, err := s.client.PendingNonceAt(ctx, fromAddress)
	log.Printf("nonce %d\n", nonce)

	value := big.NewInt(0)
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	log.Printf("gas price %d\n", gasPrice.Int64())

	if err != nil {
		log.Fatal(err)
		return err
	}

	transferFnSignature := []byte("transfer(address,uint256)")

	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	log.Printf("method id %s\n", hexutil.Encode(methodID))

	paddedAddress := common.LeftPadBytes(to.Bytes(), 32)
	amount := new(big.Int)
	// USDc decimal 6
	amount.SetInt64(valueAmount)

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte

	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	//gasLimit, err := s.client.EstimateGas(ctx, ethereum.CallMsg{
	//	To:   &to,
	//	Data: data,
	//})
	gasLimit := uint64(100_000)
	log.Printf("gas limit - %d\n", gasLimit)

	if err != nil {
		log.Println(err)
		return err

	}

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := s.client.NetworkID(ctx)
	//chainID := big.NewInt(1)
	log.Printf("chain id - %d\n", chainID.Int64())

	if err != nil {
		log.Println(err)
		return err
	}
	//types.HomesteadSigner{}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), pk)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("send transaction to %s\n", to.String())

	err = s.client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println(err)
		return err
	}
	ctx, _ = context.WithTimeout(ctx, 2*time.Second)

	var retryCount = 0
	for {
		time.Sleep(2 * time.Second)
		log.Printf("get receipt hash #%s retry number #%d ", tx.Hash().String(), retryCount)
		receipt, err := s.client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			return err
		}

		if receipt != nil && receipt.BlockNumber != nil {
			if receipt.Status == 1 || retryCount > 9 {
				log.Printf("receipt status - %d", receipt.Status)
				break
			} else {
				retryCount++
				continue
			}
		}
	}

	return nil
}

func (s *Service) FetchTokenBalance(tokenAddress common.Address, address common.Address) (*big.Int, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	out, err := s.abiCall(ctx, &(tokenAddress), MethodErc20Balance, address)
	if err != nil {
		return big.NewInt(0), err
	}
	tokenBalance := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return tokenBalance, nil
}

func (s *Service) encodeAbiCall(method string, params ...interface{}) ([]byte, error) {

	input, err := s.abi.Pack(method, params...)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func (s *Service) abiCall(
	ctx context.Context,
	contractAddress *common.Address,
	method string,
	params ...interface{},
) ([]interface{}, error) {

	input, err := s.encodeAbiCall(method, params...)
	if err != nil {
		return nil, err
	}

	callData := ethereum.CallMsg{
		To:   contractAddress,
		Data: input,
	}

	data, callErr := s.client.CallContract(ctx, callData, nil)

	if callErr != nil {
		return nil, callErr
	}

	response, unpackErr := s.abi.Unpack(method, data)

	if unpackErr != nil {
		return nil, unpackErr
	}
	return response, nil
}
