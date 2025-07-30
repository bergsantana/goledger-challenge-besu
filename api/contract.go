package api

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"

	//"time"

	log "github.com/sirupsen/logrus"

	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

type ContractClient struct {
	client      *ethclient.Client
	auth        *bind.TransactOpts
	address     common.Address
	contractAbi abi.ABI
}

func LoadContract() (*ContractClient, error) {
	client, err := ethclient.Dial(os.Getenv("NODE_URL"))

	if err != nil {
		log.Fatalf("\nError loading node: %v\n", err)
		return nil, err
	}

	_, err = client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("\nFailed to get chain ID: %v\n", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(os.Getenv("PRIVATE_KEY"), "0x"))
	if err != nil {
		log.Fatalf("\nError loading private key: %v\n", err)

		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("\nError reading Public Key: %v\n", err)

		return nil, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("🔑 Using address:", fromAddress.Hex())

	chainID := big.NewInt(1337)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("\nError on auth: %v\n", err)

		return nil, err
	}

	address := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	code, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("\n❌ Error reading contract code: %v\n", err)
	}

	if len(code) == 0 {
		log.Fatalf("\n❌ No contract deployed at address %s\n", address.Hex())
	} else {
		log.Println("✅ Contract is deployed at", address.Hex())
	}

	// Load ABI
	raw, err := os.ReadFile("abi/SimpleStorage.json")
	if err != nil {
		return nil, err
	}

	var artifact struct {
		ABI json.RawMessage `json:"abi"`
	}
	err = json.Unmarshal(raw, &artifact)

	if err != nil {
		return nil, err
	}

	parsedAbi, err := abi.JSON(strings.NewReader(string(artifact.ABI)))
	if err != nil {
		return nil, err
	}

	return &ContractClient{
		client:      client,
		auth:        auth,
		address:     address,
		contractAbi: parsedAbi,
	}, nil
}

func (c *ContractClient) SetValue(value int64) (*types.Transaction, error) {

	input, err := c.contractAbi.Pack("set", big.NewInt(value))
	if err != nil {
		return nil, err
	}

	nonce, err := c.client.PendingNonceAt(context.Background(), c.auth.From)
	if err != nil {
		return nil, err
	}

	gasPrice, _ := c.client.SuggestGasPrice(context.Background())

	tx := types.NewTransaction(nonce, c.address, big.NewInt(0), 300000, gasPrice, input)

	signedTx, err := c.auth.Signer(c.auth.From, tx)
	if err != nil {
		return nil, err
	}

	err = c.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (c *ContractClient) GetValue() (*big.Int, error) {

	out, err := c.contractAbi.Pack("get")

	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		To:   &c.address,
		Data: out,
	}

	res, err := c.client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, err
	}

	var value *big.Int
	err = c.contractAbi.UnpackIntoInterface(&value, "get", res)
	return value, err
}
