package ethclient

import (
	"context"
	"eolia-signer/types"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	eth               *ethclient.Client
	entrypoint        *abi.ABI
	entrypointAddress *common.Address
	factory           *abi.ABI
	factoryAddress    *common.Address
	account           *abi.ABI
}

func NewClient(url string, entrypoint string, factory string) *Client {
	cli, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	abiData, err := os.ReadFile("ethclient/abi/entrypoint.abi.json")
	if err != nil {
		log.Fatalf("Failed to read entrypoint ABI file: %v", err)
	}

	entryPointAddress := common.HexToAddress(entrypoint)
	entryAbi, err := abi.JSON(strings.NewReader(string(abiData)))
	if err != nil {
		log.Fatalf("Failed to parse entrypoint ABI: %v", err)
	}

	abiData, err = os.ReadFile("ethclient/abi/factory.abi.json")
	if err != nil {
		log.Fatalf("Failed to read factory ABI file: %v", err)
	}

	factoryAddress := common.HexToAddress(factory)
	factoryAbi, err := abi.JSON(strings.NewReader(string(abiData)))
	if err != nil {
		log.Fatalf("Failed to parse factory ABI: %v", err)
	}

	abiData, err = os.ReadFile("ethclient/abi/account.abi.json")
	if err != nil {
		log.Fatalf("Failed to read account ABI file: %v", err)
	}

	accountAbi, err := abi.JSON(strings.NewReader(string(abiData)))
	if err != nil {
		log.Fatalf("Failed to parse account ABI: %v", err)
	}

	return &Client{eth: cli, entrypoint: &entryAbi, entrypointAddress: &entryPointAddress, factory: &factoryAbi, factoryAddress: &factoryAddress, account: &accountAbi}
}

func (c *Client) GetCalculatedAddress(address common.Address, salt *big.Int) (common.Address, error) {
	data, err := c.factory.Pack("getAddress", address, salt)
	if err != nil {
		return common.Address{}, err
	}

	callMsg := ethereum.CallMsg{
		To:   c.factoryAddress,
		Data: data,
	}

	output, err := c.eth.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return common.Address{}, err
	}

	var result common.Address
	err = c.factory.UnpackIntoInterface(&result, "getAddress", output)
	if err != nil {
		return common.Address{}, err
	}

	fmt.Printf("Calculated address: %s\n", result.Hex())
	return result, nil
}

func (c *Client) AccountNeedsInitialization(sender common.Address, owner common.Address) []byte {
	byteCode, err := c.eth.CodeAt(context.Background(), sender, nil)
	if err != nil {
		fmt.Printf("Error fetching bytecode for %s: %v\n", sender.Hex(), err)
		return nil
	}

	fmt.Printf("Bytecode length for %s: %d\n", sender.Hex(), len(byteCode))

	if len(byteCode) == 0 {
		encodedFunction, err := c.factory.Pack("createAccount", owner, big.NewInt(0))
		if err != nil {
			return nil
		}

		addressInBytes := c.factoryAddress.Bytes()
		initCode := append(addressInBytes, encodedFunction...)

		fmt.Printf("Account needs initialization, init code: %x\n", initCode)
		fmt.Printf("Factory address: %s\n", c.factoryAddress.Hex())
		fmt.Printf("Sender address: %s\n", sender.Hex())
		fmt.Printf("Owner address: %s\n", owner.Hex())
		fmt.Printf("Encoded function: %x\n", encodedFunction)

		return initCode
	}

	return nil
}

func (c *Client) GetNonce(sender common.Address) *big.Int {
	byteCode, err := c.eth.CodeAt(context.Background(), sender, nil)
	if err != nil {
		fmt.Printf("Error fetching bytecode for %s: %v\n", sender.Hex(), err)
		return nil
	}

	if len(byteCode) == 0 {
		return big.NewInt(0)
	} else {
		encodedFunction, err := c.account.Pack("getNonce")
		if err != nil {
			return nil
		}
		callMsg := ethereum.CallMsg{
			To:   &sender,
			Data: encodedFunction,
		}
		output, err := c.eth.CallContract(context.Background(), callMsg, nil)
		if err != nil {
			fmt.Printf("Error calling getNonce for %s: %v\n", sender.Hex(), err)
			return nil
		}
		var nonce *big.Int
		err = c.account.UnpackIntoInterface(&nonce, "getNonce", output)
		if err != nil {
			fmt.Printf("Error unpacking getNonce result for %s: %v\n", sender.Hex(), err)
			return nil
		}
		return nonce
	}
}

func (c *Client) GetUserOpHash(userOp *types.PackedUserOperation) common.Hash {
	encodedFunction, err := c.entrypoint.Pack("getUserOpHash", userOp)
	if err != nil {
		log.Fatalf("Failed to pack getUserOpHash: %v", err)
	}

	callMsg := ethereum.CallMsg{
		To:   c.entrypointAddress,
		Data: encodedFunction,
	}

	output, err := c.eth.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Failed to call getUserOpHash: %v", err)
	}

	var result common.Hash
	err = c.entrypoint.UnpackIntoInterface(&result, "getUserOpHash", output)
	if err != nil {
		log.Fatalf("Failed to unpack getUserOpHash result: %v", err)
	}

	return result
}
