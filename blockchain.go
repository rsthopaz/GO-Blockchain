package main

import (
	"log"
	"math/big"
	"os"
	"strings"
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)


var ethClient *ethclient.Client
var contractABI abi.ABI
var contractAddress common.Address
var auth *bind.TransactOpts

func initBlockchain() {
	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	privateKeyHex := os.Getenv("BLOCKCHAIN_PRIVATE_KEY")
	contractAddr := os.Getenv("CONTRACT_ADDRESS")

	if rpcURL == "" || privateKeyHex == "" || contractAddr == "" {
		log.Fatal("Blockchain env not set")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("Failed to connect to Ethereum:", err)
	}
	ethClient = client

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		log.Fatal("Invalid private key:", err)
	}

	chainID := big.NewInt(11155111) // Sepolia
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	auth.Value = big.NewInt(0)       // no ETH sent
	auth.GasLimit = uint64(300000)   // cukup aman

	abiJSON := `[
  {
    "inputs": [
      {
        "internalType": "enum AssetEventRegistry.EventType",
        "name": "eventType",
        "type": "uint8"
      },
      {
        "internalType": "string",
        "name": "assetId",
        "type": "string"
      },
      {
        "internalType": "string",
        "name": "assetCode",
        "type": "string"
      },
      {
        "internalType": "string",
        "name": "assetName",
        "type": "string"
      },
      {
        "internalType": "string",
        "name": "category",
        "type": "string"
      },
      {
        "internalType": "uint256",
        "name": "value",
        "type": "uint256"
      }
    ],
    "name": "recordEvent",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  }
]`

	contractABI, err = abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		log.Fatal("Invalid ABI:", err)
	}

	contractAddress = common.HexToAddress(contractAddr)

	log.Println("✅ Blockchain initialized")
}

func sendAssetEvent(eventType uint8, asset Asset) error {

	valueInt := big.NewInt(int64(asset.Value))

	input, err := contractABI.Pack(
		"recordEvent",
		eventType,
		asset.ID.String(),   // ✅ uuid → string
		asset.Code,          // string
		asset.Name,          // string
		asset.Category,      // string
		valueInt,            // ✅ uint256
	)
	if err != nil {
		return err
	}

	nonce, err := ethClient.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return err
	}

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	msg := ethereum.CallMsg{
		From: auth.From,
		To:   &contractAddress,
		Data: input,
	}

	gasLimit, err := ethClient.EstimateGas(context.Background(), msg)
	if err != nil {
		return err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = gasLimit

	tx := types.NewTransaction(
		nonce,
		contractAddress,
		big.NewInt(0),
		gasLimit,
		gasPrice,
		input,
	)

	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return err
	}

	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	log.Println("📦 Blockchain event sent:", eventType, asset.ID.String())
	log.Println("   TxHash :", signedTx.Hash().Hex())
	log.Println("   Event  :", eventType)
	log.Println("   AssetID:", asset.ID.String())

	
	return nil
}

