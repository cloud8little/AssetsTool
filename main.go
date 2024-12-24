package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
)

const (
	// DepositABI                = `[{"inputs":[{"internalType":"uint32","name":"clientChainID","type":"uint32"},{"internalType":"bytes","name":"assetsAddress","type":"bytes"},{"internalType":"bytes","name":"stakerAddress","type":"bytes"},{"internalType":"uint256","name":"opAmount","type":"uint256"}],"name":"depositLST","outputs":[{"internalType":"bool","name":"success","type":"bool"},{"internalType":"uint256","name":"latestAssetState","type":"uint256"}],"stateMutability":"nonpayable","type":"function"}]`
	DepositABI = `[
    {
      "inputs":
      [
        {
          "internalType": "uint32",
          "name": "clientChainID",
          "type": "uint32"
        },
        {
          "internalType": "bytes",
          "name": "assetsAddress",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "stakerAddress",
          "type": "bytes"
        },
        {
          "internalType": "uint256",
          "name": "opAmount",
          "type": "uint256"
        }
      ],
      "name": "depositLST",
      "outputs":
      [
        {
          "internalType": "bool",
          "name": "success",
          "type": "bool"
        },
        {
          "internalType": "uint256",
          "name": "latestAssetState",
          "type": "uint256"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs":
      [
        {
          "internalType": "uint32",
          "name": "clientChainID",
          "type": "uint32"
        },
        {
          "internalType": "bytes",
          "name": "validatorID",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "stakerAddress",
          "type": "bytes"
        },
        {
          "internalType": "uint256",
          "name": "opAmount",
          "type": "uint256"
        }
      ],
      "name": "depositNST",
      "outputs":
      [
        {
          "internalType": "bool",
          "name": "success",
          "type": "bool"
        },
        {
          "internalType": "uint256",
          "name": "latestAssetState",
          "type": "uint256"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getClientChains",
      "outputs":
      [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        },
        {
          "internalType": "uint32[]",
          "name": "",
          "type": "uint32[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs":
      [
        {
          "internalType": "uint32",
          "name": "clientChainID",
          "type": "uint32"
        }
      ],
      "name": "isRegisteredClientChain",
      "outputs":
      [
        {
          "internalType": "bool",
          "name": "success",
          "type": "bool"
        },
        {
          "internalType": "bool",
          "name": "isRegistered",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs":
      [
        {
          "internalType": "uint32",
          "name": "clientChainID",
          "type": "uint32"
        },
        {
          "internalType": "uint8",
          "name": "addressLength",
          "type": "uint8"
        },
        {
          "internalType": "string",
          "name": "name",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "metaInfo",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "signatureType",
          "type": "string"
        }
      ],
      "name": "registerOrUpdateClientChain",
      "outputs":
      [
        {
          "internalType": "bool",
          "name": "success",
          "type": "bool"
        },
        {
          "internalType": "bool",
          "name": "updated",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs":
      [
        {
          "internalType": "uint32",
          "name": "clientChainId",
          "type": "uint32"
        },
        {
          "internalType": "bytes",
          "name": "token",
          "type": "bytes"
        },
        {
          "internalType": "uint8",
          "name": "decimals",
          "type": "uint8"
        },
        {
          "internalType": "string",
          "name": "name",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "metaData",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "oracleInfo",
          "type": "string"
        }
      ],
      "name": "registerToken",
      "outputs":
      [
        {
          "internalType": "bool",
          "name": "success",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs":
      [
        {
          "internalType": "uint32",
          "name": "clientChainId",
          "type": "uint32"
        },
        {
          "internalType": "bytes",
          "name": "token",
          "type": "bytes"
        },
        {
          "internalType": "string",
          "name": "metaData",
          "type": "string"
        }
      ],
      "name": "updateToken",
      "outputs":
      [
        {
          "internalType": "bool",
          "name": "success",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs":
      [
        {
          "internalType": "uint32",
          "name": "clientChainID",
          "type": "uint32"
        },
        {
          "internalType": "bytes",
          "name": "assetsAddress",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "withdrawAddress",
          "type": "bytes"
        },
        {
          "internalType": "uint256",
          "name": "opAmount",
          "type": "uint256"
        }
      ],
      "name": "withdrawLST",
      "outputs":
      [
        {
          "internalType": "bool",
          "name": "success",
          "type": "bool"
        },
        {
          "internalType": "uint256",
          "name": "latestAssetState",
          "type": "uint256"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs":
      [
        {
          "internalType": "uint32",
          "name": "clientChainID",
          "type": "uint32"
        },
        {
          "internalType": "bytes",
          "name": "validatorID",
          "type": "bytes"
        },
        {
          "internalType": "bytes",
          "name": "withdrawAddress",
          "type": "bytes"
        },
        {
          "internalType": "uint256",
          "name": "opAmount",
          "type": "uint256"
        }
      ],
      "name": "withdrawNST",
      "outputs":
      [
        {
          "internalType": "bool",
          "name": "success",
          "type": "bool"
        },
        {
          "internalType": "uint256",
          "name": "latestAssetState",
          "type": "uint256"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`
	DelegateABI               = `[{"inputs":[{"internalType":"uint32","name":"clientChainID","type":"uint32"},{"internalType":"bytes","name":"staker","type":"bytes"},{"internalType":"bytes","name":"operator","type":"bytes"}],"name":"associateOperatorWithStaker","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint32","name":"clientChainID","type":"uint32"},{"internalType":"uint64","name":"lzNonce","type":"uint64"},{"internalType":"bytes","name":"assetsAddress","type":"bytes"},{"internalType":"bytes","name":"stakerAddress","type":"bytes"},{"internalType":"bytes","name":"operatorAddr","type":"bytes"},{"internalType":"uint256","name":"opAmount","type":"uint256"}],"name":"delegate","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint32","name":"clientChainID","type":"uint32"},{"internalType":"bytes","name":"staker","type":"bytes"}],"name":"dissociateOperatorFromStaker","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint32","name":"clientChainID","type":"uint32"},{"internalType":"uint64","name":"lzNonce","type":"uint64"},{"internalType":"bytes","name":"assetsAddress","type":"bytes"},{"internalType":"bytes","name":"stakerAddress","type":"bytes"},{"internalType":"bytes","name":"operatorAddr","type":"bytes"},{"internalType":"uint256","name":"opAmount","type":"uint256"}],"name":"undelegate","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`
	depositPrecompileAddress  = "0x0000000000000000000000000000000000000804"
	delegatePrecompileAddress = "0x0000000000000000000000000000000000000805"
)

var (
	privateKey     string
	defaultAssetID string
	layerZeroID    uint32
)

var rootCmd = &cobra.Command{
	Use:   "assetcli",
	Short: "Asset CLI tool",
}

var depositCmd = &cobra.Command{
	Use:   "deposit",
	Short: "Deposit to Exocore",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpcUrl")
		staker, _ := cmd.Flags().GetString("staker")
		amountStr, _ := cmd.Flags().GetString("amount")
		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			log.Fatalf("Invalid amount: %s", amountStr)
		}
		err := deposit_(rpcUrl, staker, amount)
		if err != nil {
			log.Fatalf("Failed to deposit: %v", err)
		}
	},
}

var delegateCmd = &cobra.Command{
	Use:   "delegate",
	Short: "Delegate to Exocore",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpcUrl")
		staker, _ := cmd.Flags().GetString("staker")
		operator, _ := cmd.Flags().GetString("operator")
		amountStr, _ := cmd.Flags().GetString("amount")
		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			log.Fatalf("Invalid amount: %s", amountStr)
		}
		err := delegateTo_(rpcUrl, staker, operator, amount)
		if err != nil {
			log.Fatalf("Failed to delegate: %v", err)
		}
	},
}

// undelegte command
var undelegateCmd = &cobra.Command{
	Use:   "undelegate",
	Short: "Undelegate from Exocore",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpcUrl")
		staker, _ := cmd.Flags().GetString("staker")
		operator, _ := cmd.Flags().GetString("operator")
		amountStr, _ := cmd.Flags().GetString("amount")
		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			log.Fatalf("Invalid amount: %s", amountStr)
		}
		err := undelegate_(rpcUrl, staker, operator, amount)
		if err != nil {
			log.Fatalf("Failed to undelegate: %v", err)
		}
	},
}

var selfDelegateCmd = &cobra.Command{
	Use:   "self-delegate",
	Short: "Self delegate to Exocore",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpcUrl")
		staker, _ := cmd.Flags().GetString("staker")
		operator, _ := cmd.Flags().GetString("operator")
		err := selfDelegate_(rpcUrl, staker[2:], operator)
		if err != nil {
			log.Fatalf("Failed to self delegate: %v", err)
		}
	},
}

var withdrawLSTCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "WithdrawLST from Exocore",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpcUrl")
		staker, _ := cmd.Flags().GetString("staker")
		amountStr, _ := cmd.Flags().GetString("amount")
		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			log.Fatalf("Invalid amount: %s", amountStr)
		}
		err := withdrawLST_(rpcUrl, staker, amount)
		if err != nil {
			log.Fatalf("Failed to withdrawfe: %v", err)
		}
	},
}

func main() {
	rootCmd.PersistentFlags().StringVar(&privateKey, "privateKey", "", "Private key for transactions")
	rootCmd.PersistentFlags().StringVar(&defaultAssetID, "defaultAssetID", "", "Default asset ID")
	rootCmd.PersistentFlags().Uint32Var(&layerZeroID, "layerZeroID", 101, "LayerZero ID")

	rootCmd.AddCommand(depositCmd)
	rootCmd.AddCommand(delegateCmd)
	rootCmd.AddCommand(selfDelegateCmd)
	rootCmd.AddCommand(undelegateCmd)
	rootCmd.AddCommand(withdrawLSTCmd)

	depositCmd.Flags().String("rpcUrl", "http://localhost:8545", "Exocore RPC URL")
	depositCmd.Flags().String("staker", "", "Staker address")
	depositCmd.Flags().String("amount", "0", "Amount to deposit")

	withdrawLSTCmd.Flags().String("rpcUrl", "http://localhost:8545", "Exocore RPC URL")
	withdrawLSTCmd.Flags().String("staker", "", "Staker address")
	withdrawLSTCmd.Flags().String("amount", "0", "Amount to deposit")

	delegateCmd.Flags().String("rpcUrl", "http://localhost:8545", "Exocore RPC URL")
	delegateCmd.Flags().String("staker", "", "Staker address")
	delegateCmd.Flags().String("operator", "", "Operator address")
	delegateCmd.Flags().String("amount", "0", "Amount to delegate")

	undelegateCmd.Flags().String("rpcUrl", "http://localhost:8545", "Exocore RPC URL")
	undelegateCmd.Flags().String("staker", "", "Staker address")
	undelegateCmd.Flags().String("operator", "", "Operator address")
	undelegateCmd.Flags().String("amount", "0", "Amount to undelegate")

	selfDelegateCmd.Flags().String("rpcUrl", "http://localhost:8545", "Exocore RPC URL")
	selfDelegateCmd.Flags().String("staker", "", "Staker address")
	selfDelegateCmd.Flags().String("operator", "", "Operator address")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func deposit_(rpcUrl, stakerAddress string, amount *big.Int) error {
	depositAddr := common.HexToAddress(depositPrecompileAddress)
	assetAddr := common.HexToAddress(defaultAssetID)
	stakerAddr := common.HexToAddress(stakerAddress)
	opAmount := amount

	_, ethClient, err := connectToEthereum(rpcUrl)
	if err != nil {
		return err
	}

	sk, callAddr, err := getPrivateKeyAndAddress(privateKey)
	if err != nil {
		return err
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return err
	}

	depositAbi, err := abi.JSON(strings.NewReader(DepositABI))
	if err != nil {
		return err
	}

	data, err := depositAbi.Pack("depositLST", layerZeroID, paddingAddressTo32(assetAddr), paddingAddressTo32(stakerAddr), opAmount)
	if err != nil {
		return err
	}

	txID, err := sendTransaction(ethClient, chainID, callAddr, sk, depositAddr, data)
	if err != nil {
		return err
	}

	fmt.Println("Deposit Transaction ID:", txID)
	return waitForTransaction(ethClient, txID)
}

func delegateTo_(rpcUrl, stakerAddress, operatorBench32Str string, amount *big.Int) error {
	delegateAddr := common.HexToAddress(delegatePrecompileAddress)
	assetAddr := common.HexToAddress(defaultAssetID)
	stakerAddr := common.HexToAddress(stakerAddress)
	operatorAddr := []byte(operatorBench32Str)
	opAmount := amount
	lzNonce := uint64(0)
	_, ethClient, err := connectToEthereum(rpcUrl)
	if err != nil {
		return err
	}

	sk, callAddr, err := getPrivateKeyAndAddress(privateKey)
	if err != nil {
		return err
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return err
	}

	delegateAbi, err := abi.JSON(strings.NewReader(DelegateABI))
	if err != nil {
		return err
	}

	data, err := delegateAbi.Pack("delegate", layerZeroID, lzNonce, paddingAddressTo32(assetAddr), paddingAddressTo32(stakerAddr), operatorAddr, opAmount)
	if err != nil {
		return err
	}

	txID, err := sendTransaction(ethClient, chainID, callAddr, sk, delegateAddr, data)
	if err != nil {
		return err
	}

	fmt.Println("Delegate To Transaction ID:", txID)
	return waitForTransaction(ethClient, txID)
}

func undelegate_(rpcUrl, stakerAddress, operatorBench32Str string, amount *big.Int) error {
	delegateAddr := common.HexToAddress(delegatePrecompileAddress)
	assetAddr := common.HexToAddress(defaultAssetID)
	stakerAddr := common.HexToAddress(stakerAddress)
	operatorAddr := []byte(operatorBench32Str)
	opAmount := amount
	lzNonce := uint64(0)

	_, ethClient, err := connectToEthereum(rpcUrl)
	if err != nil {
		return err
	}

	sk, callAddr, err := getPrivateKeyAndAddress(privateKey)
	if err != nil {
		return err
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return err
	}

	delegateAbi, err := abi.JSON(strings.NewReader(DelegateABI))
	if err != nil {
		return err
	}

	data, err := delegateAbi.Pack("undelegate", layerZeroID, lzNonce, paddingAddressTo32(assetAddr), paddingAddressTo32(stakerAddr), operatorAddr, opAmount)
	if err != nil {
		return err
	}

	txID, err := sendTransaction(ethClient, chainID, callAddr, sk, delegateAddr, data)
	if err != nil {
		return err
	}

	fmt.Println("Undelegate Transaction ID:", txID)
	return waitForTransaction(ethClient, txID)
}

func selfDelegate_(rpcUrl, stakerAddr, operatorBench32Str string) error {
	delegateAddr := common.HexToAddress(delegatePrecompileAddress)
	operatorAddr := []byte(operatorBench32Str)

	_, ethClient, err := connectToEthereum(rpcUrl)
	if err != nil {
		return err
	}

	sk, callAddr, err := getPrivateKeyAndAddress(privateKey)
	if err != nil {
		return err
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return err
	}

	delegateAbi, err := abi.JSON(strings.NewReader(DelegateABI))
	if err != nil {
		return err
	}

	staker, err := hex.DecodeString(stakerAddr)
	if err != nil {
		return err
	}

	data, err := delegateAbi.Pack("associateOperatorWithStaker", layerZeroID, staker, operatorAddr)
	if err != nil {
		return err
	}

	txID, err := sendTransaction(ethClient, chainID, callAddr, sk, delegateAddr, data)
	if err != nil {
		return err
	}

	fmt.Println("Self Delegate Transaction ID:", txID)
	return waitForTransaction(ethClient, txID)
}

func withdrawLST_(rpcUrl, stakerAddress string, amount *big.Int) error {
	depositAddr := common.HexToAddress(depositPrecompileAddress)
	assetAddr := common.HexToAddress(defaultAssetID)
	stakerAddr := common.HexToAddress(stakerAddress)
	opAmount := amount

	_, ethClient, err := connectToEthereum(rpcUrl)
	if err != nil {
		return err
	}

	sk, callAddr, err := getPrivateKeyAndAddress(privateKey)
	if err != nil {
		return err
	}

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		return err
	}

	depositAbi, err := abi.JSON(strings.NewReader(DepositABI))
	if err != nil {
		return err
	}

	data, err := depositAbi.Pack("withdrawLST", layerZeroID, paddingAddressTo32(assetAddr), paddingAddressTo32(stakerAddr), opAmount)
	if err != nil {
		return err
	}

	txID, err := sendTransaction(ethClient, chainID, callAddr, sk, depositAddr, data)
	if err != nil {
		return err
	}

	fmt.Println("Withdraw LST Transaction ID:", txID)
	return waitForTransaction(ethClient, txID)
}

func connectToEthereum(nodeURL string) (*rpc.Client, *ethclient.Client, error) {
	client, err := rpc.DialContext(context.Background(), nodeURL)
	if err != nil {
		return nil, nil, err
	}
	ethClient := ethclient.NewClient(client)
	return client, ethClient, nil
}

func getPrivateKeyAndAddress(privateKey string) (*ecdsa.PrivateKey, common.Address, error) {
	sk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, common.Address{}, err
	}
	callAddr := crypto.PubkeyToAddress(sk.PublicKey)
	return sk, callAddr, nil
}

func paddingAddressTo32(address common.Address) []byte {
	paddingLen := 32 - len(address)
	ret := make([]byte, len(address))
	copy(ret, address[:])
	for i := 0; i < paddingLen; i++ {
		ret = append(ret, 0)
	}
	fmt.Println("Padded address:", hexutil.Encode(ret))
	return ret
}

func sendTransaction(client *ethclient.Client, chainID *big.Int, from common.Address, sk *ecdsa.PrivateKey, to common.Address, data []byte) (string, error) {
	ctx := context.Background()
	nonce, err := client.NonceAt(ctx, from, nil)
	if err != nil {
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	gasLimit := uint64(500000)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    big.NewInt(0),
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})
	signer := types.LatestSignerForChainID(chainID)
	signTx, err := types.SignTx(tx, signer, sk)
	if err != nil {
		return "", err
	}

	fmt.Println("the txID is:", signTx.Hash().String())
	msg := ethereum.CallMsg{
		From: from,
		To:   tx.To(),
		Data: tx.Data(),
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	fmt.Println("The bool value returned by the contract:", result)
	if err != nil {
		fmt.Println("Failed to call contract:", err)
	}
	if isZeroByteArray(result) {
		fmt.Println("Failed to call contract,The bool value returned by the contract is false")
	}
	err = client.SendTransaction(ctx, signTx)
	if err != nil {
		return "", err
	}
	return signTx.Hash().String(), nil
}
func isZeroByteArray(byteArray []byte) bool {
	for _, b := range byteArray {
		if b != 0 {
			return false
		}
	}
	return true
}
func waitForTransaction(client *ethclient.Client, txID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	txHash := common.HexToHash(txID)
	tx, _, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %v", err)
	}

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with status: %v", receipt.Status)
	}

	// fmt.Println("Transaction mined successfully with receipt:", receipt)
	return nil
}
