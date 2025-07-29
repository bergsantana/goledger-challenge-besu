package api

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"time"

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
	//fromAddress common.Address
}

func LoadContract() (*ContractClient, error) {
	client, err := ethclient.Dial(os.Getenv("NODE_URL"))

	if err != nil {
		fmt.Printf("\nError loading node: %v\n", err)
		return nil, err
	}

	_, err = client.ChainID(context.Background())

	if err != nil {
		fmt.Printf("\nFailed to get chain ID: %v\n", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(os.Getenv("PRIVATE_KEY"), "0x"))
	if err != nil {
		fmt.Printf("\nError loading private key: %v\n", err)

		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Printf("\nError reading Public Key: %v\n", err)

		return nil, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("🔑 Using address:", fromAddress.Hex())

	chainID := big.NewInt(1337)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		fmt.Printf("\nError on auth: %v\n", err)

		return nil, err
	}

	address := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	code, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		fmt.Printf("\n❌ Error reading contract code: %v\n", err)
	}

	if len(code) == 0 {
		fmt.Printf("\n❌ No contract deployed at address %s\n", address.Hex())
	} else {
		fmt.Println("✅ Contract is deployed at", address.Hex())
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

	CallContract()

	return &ContractClient{
		client:      client,
		auth:        auth,
		address:     address,
		contractAbi: parsedAbi,
		//fromAddress: fromAddress,
	}, nil
}

func (c *ContractClient) SetValue(value int64) (*types.Transaction, error) {
	CallContract()
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
	CallContract()
	out, err := c.contractAbi.Pack("get")

	fmt.Println("c.address")
	fmt.Println(c.address)
	fmt.Println(&c.address)

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

func CallContract() {
	var result interface{}
	raw, err := os.ReadFile("abi/SimpleStorage.json")
	if err != nil {
		log.Fatalf("Error loading ABI file: %v", err)
	}

	var artifact struct {
		ABI json.RawMessage `json:"abi"`
	}
	err = json.Unmarshal(raw, &artifact)

	if err != nil {
		log.Fatalf("Error parsing ABI File: %v", err)
	}
	abi, err := abi.JSON(strings.NewReader(string(artifact.ABI))) // found under besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json
	if err != nil {
		log.Fatalf("error parsing abi: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, "http://localhost:8545") // e.g., http://localhost:8545
	if err != nil {
		log.Fatalf("error connecting to eth client: %v", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress("enode://38e0f5cd19120e150c10ba3754de21d7cda215f2bb3a8d7360abf55e4f7205e41d601cd7e55a6afd23ccf14975f2c3f24a129934c6d091bf0fa8ba7ea0463cbc@127.0.0.1:30303") // will be returned during startDev.sh execution
	caller := bind.CallOpts{
		Pending: false,
		Context: ctx,
	}

	boundContract := bind.NewBoundContract(
		contractAddress,
		abi,
		client,
		client,
		client,
	)

	var output []interface{}
	err = boundContract.Call(&caller, &output, "get")
	if err != nil {
		log.Fatalf("error calling contract: %v", err)
	}
	result = output

	fmt.Println("Successfully called contract!", result)
}
