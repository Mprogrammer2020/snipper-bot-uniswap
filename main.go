package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to an Ethereum client
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + os.Getenv("INFURA_PROJECT_ID"))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Load the Uniswap contract ABI
	uniswapABI := os.Getenv("UNISWAP_CONTRACT_ABI")

	// Parse the contract ABI
	contractABI, err := abi.JSON(strings.NewReader(uniswapABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Create a new Uniswap contract instance
	contractAddress := os.Getenv("UNISWAP_CONTRACT_ADDRESS")
	uniswapContract, err := NewUniswap(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to instantiate Uniswap contract: %v", err)
	}

	// Set up the transaction parameters
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	fromAddress := os.Getenv("FROM_ADDRESS")
	toAddress := os.Getenv("TO_ADDRESS")
	amount := big.NewInt(0)
	amount.SetString(os.Getenv("AMOUNT"), 10)
	gasLimit := uint64(0)
	gasLimit, _ = strconv.ParseUint(os.Getenv("GAS_LIMIT"), 10, 64)
	gasPrice := big.NewInt(0)
	gasPrice.SetString(os.Getenv("GAS_PRICE"), 10)

	// Create a new transaction
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))
	if err != nil {
		log.Fatalf("Failed to retrieve nonce: %v", err)
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), amount, gasLimit, gasPrice, nil)

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	log.Printf("Transaction sent: %s", signedTx.Hash().Hex())
}

