package main

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"encoding/json"
	"io/ioutil"
)

type ContractClient struct {
	client         *ethclient.Client
	auth           *bind.TransactOpts
	address        common.Address
	contractAbi    abi.ABI
}

func LoadContract() (*ContractClient, error) {
	client, err := ethclient.Dial(os.Getenv("NODE_URL"))
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(os.Getenv("PRIVATE_KEY"), "0x"))
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	chainID := big.NewInt(1337)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	// Load ABI
	raw, err := ioutil.ReadFile("abi/SimpleStorage.json")
	if err != nil {
		return nil, err
	}
	parsedAbi, err := abi.JSON(strings.NewReader(string(raw)))
	if err != nil {
		return nil, err
	}

	return &ContractClient{
		client:         client,
		auth:           auth,
		address:        address,
		contractAbi:    parsedAbi,
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
