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

	DelegateABI               = `[
  {
    "type": "function",
    "name": "associateOperatorWithStaker",
    "inputs": [
      {
        "name": "clientChainID",
        "type": "uint32",
        "internalType": "uint32"
      },
      {
        "name": "staker",
        "type": "bytes",
        "internalType": "bytes"
      },
      {
        "name": "operator",
        "type": "bytes",
        "internalType": "bytes"
      }
    ],
    "outputs": [
      {
        "name": "success",
        "type": "bool",
        "internalType": "bool"
      }
    ],
    "stateMutability": "nonpayable"
  },
  {
    "type": "function",
    "name": "delegate",
    "inputs": [
      {
        "name": "clientChainID",
        "type": "uint32",
        "internalType": "uint32"
      },
      {
        "name": "assetsAddress",
        "type": "bytes",
        "internalType": "bytes"
      },
      {
        "name": "stakerAddress",
        "type": "bytes",
        "internalType": "bytes"
      },
      {
        "name": "operatorAddr",
        "type": "bytes",
        "internalType": "bytes"
      },
      {
        "name": "opAmount",
        "type": "uint256",
        "internalType": "uint256"
      }
    ],
    "outputs": [
      {
        "name": "success",
        "type": "bool",
        "internalType": "bool"
      }
    ],
    "stateMutability": "nonpayable"
  },
  {
    "type": "function",
    "name": "dissociateOperatorFromStaker",
    "inputs": [
      {
        "name": "clientChainID",
        "type": "uint32",
        "internalType": "uint32"
      },
      {
        "name": "staker",
        "type": "bytes",
        "internalType": "bytes"
      }
    ],
    "outputs": [
      {
        "name": "success",
        "type": "bool",
        "internalType": "bool"
      }
    ],
    "stateMutability": "nonpayable"
  },
  {
    "type": "function",
    "name": "undelegate",
    "inputs": [
      {
        "name": "clientChainID",
        "type": "uint32",
        "internalType": "uint32"
      },
      {
        "name": "assetsAddress",
        "type": "bytes",
        "internalType": "bytes"
      },
      {
        "name": "stakerAddress",
        "type": "bytes",
        "internalType": "bytes"
      },
      {
        "name": "operatorAddr",
        "type": "bytes",
        "internalType": "bytes"
      },
      {
        "name": "opAmount",
        "type": "uint256",
        "internalType": "uint256"
      },
      {
        "name": "instantUnbond",
        "type": "bool",
        "internalType": "bool"
      }
    ],
    "outputs": [
      {
        "name": "success",
        "type": "bool",
        "internalType": "bool"
      }
    ],
    "stateMutability": "nonpayable"
  }
]
`
	depositPrecompileAddress  = "0x0000000000000000000000000000000000000804"
	delegatePrecompileAddress = "0x0000000000000000000000000000000000000805"
)

var (
	privateKey     string
	defaultAssetID string
	layerZeroID    uint32
)

var rootCmd = &cobra.Command{
	Use:     "assetcli",
	Short:   "Asset CLI tool",
	Version: "0.0.7",
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

var depositNSTCmd = &cobra.Command{
	Use:   "depositNST",
	Short: "DepositNST to Exocore",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpcUrl")
		staker, _ := cmd.Flags().GetString("staker")
		amountStr, _ := cmd.Flags().GetString("amount")
		amount, ok := new(big.Int).SetString(amountStr, 10)
		pubkey, _ := cmd.Flags().GetString("pubkey")
		if !ok {
			log.Fatalf("Invalid amount: %s", amountStr)
		}
		err := depositNST_(rpcUrl, pubkey, staker, amount)
		if err != nil {
			log.Fatalf("Failed to depositNST: %v", err)
		}
	},
}

var withdrawNSTCmd = &cobra.Command{
	Use:   "withdrawNST",
	Short: "WithdrawNST to Exocore",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpcUrl")
		staker, _ := cmd.Flags().GetString("staker")
		amountStr, _ := cmd.Flags().GetString("amount")
		amount, ok := new(big.Int).SetString(amountStr, 10)
		pubkey, _ := cmd.Flags().GetString("pubkey")
		if !ok {
			log.Fatalf("Invalid amount: %s", amountStr)
		}
		err := withdrawNST_(rpcUrl, pubkey, staker, amount)
		if err != nil {
			log.Fatalf("Failed to withdrawNST: %v", err)
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
		instantUnbond, _ := cmd.Flags().GetBool("instantUnbond")
		if !ok {
			log.Fatalf("Invalid amount: %s", amountStr)
		}
		err := undelegate_(rpcUrl, staker, operator, amount, instantUnbond)
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

var cancelSelfDelegateCmd = &cobra.Command{
	Use:   "cancel-self-delegate",
	Short: "Cancel self delegate to Exocore",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpcUrl")
		staker, _ := cmd.Flags().GetString("staker")
		err := cancelSelfDelegate_(rpcUrl, staker[2:])
		if err != nil {
			log.Fatalf("Failed to cancel self delegate: %v", err)
		}
	},
}

var registerTokenCmd = &cobra.Command{
	Use:   "register-token",
	Short: "Register token to Exocore",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpcUrl")
		assetAddress, _ := cmd.Flags().GetString("assetAddress")
		decimals, _ := cmd.Flags().GetUint8("decimals")
		name, _ := cmd.Flags().GetString("name")
		metaData, _ := cmd.Flags().GetString("metaData")
		oracleInfo, _ := cmd.Flags().GetString("oracleInfo")
		err := registerToken_(rpcUrl, assetAddress, decimals, name, metaData, oracleInfo)
		if err != nil {
			log.Fatalf("Failed to register token: %v", err)
		}
	},
}

var updateTokenCmd = &cobra.Command{
	Use:  "update-token",
	Short: "Update token to Exocore",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpcUrl")
		assetAddress, _ := cmd.Flags().GetString("assetAddress")
		metaData, _ := cmd.Flags().GetString("metaData")
		err := updateToken_(rpcUrl, assetAddress, metaData)
		if err != nil {
			log.Fatalf("Failed to update token: %v", err)
		}
	},
}

var registerOrUpdateClientChainCmd = &cobra.Command{
	Use:   "register-or-update-client-chain",
	Short: "Register or update client chain to Exocore",
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl, _ := cmd.Flags().GetString("rpcUrl")
		clientChainID, _ := cmd.Flags().GetUint32("clientChainID")
		addressLength, _ := cmd.Flags().GetUint8("addressLength")
		name, _ := cmd.Flags().GetString("name")
		metaInfo, _ := cmd.Flags().GetString("metaInfo")
		signatureType, _ := cmd.Flags().GetString("signatureType")
		err := registerOrUpdateClientChain_(rpcUrl, clientChainID, addressLength, name, metaInfo, signatureType)
		if err != nil {
			log.Fatalf("Failed to register or update client chain: %v", err)
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
	rootCmd.AddCommand(depositNSTCmd)
	rootCmd.AddCommand(withdrawNSTCmd)
	rootCmd.AddCommand(cancelSelfDelegateCmd)
	rootCmd.AddCommand(registerTokenCmd)
	rootCmd.AddCommand(updateTokenCmd)
	rootCmd.AddCommand(registerOrUpdateClientChainCmd)

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

	depositNSTCmd.Flags().String("rpcUrl", "http://localhost:8545", "Exocore RPC URL")
	depositNSTCmd.Flags().String("staker", "", "Staker address")
	depositNSTCmd.Flags().String("amount", "0", "Amount to deposit")
	depositNSTCmd.Flags().String("pubkey", "", "pubkey")

	withdrawNSTCmd.Flags().String("rpcUrl", "http://localhost:8545", "Exocore RPC URL")
	withdrawNSTCmd.Flags().String("staker", "", "Staker address")
	withdrawNSTCmd.Flags().String("amount", "0", "Amount to deposit")
	withdrawNSTCmd.Flags().String("pubkey", "", "pubkey")

	cancelSelfDelegateCmd.Flags().String("rpcUrl", "http://localhost:8545", "Exocore RPC URL")
	cancelSelfDelegateCmd.Flags().String("staker", "", "Staker address")

	registerTokenCmd.Flags().String("rpcUrl", "http://localhost:8545", "Exocore RPC URL")
	registerTokenCmd.Flags().String("assetAddress", "", "Asset address")
	registerTokenCmd.Flags().Uint8("decimals", 0, "Decimals")
	registerTokenCmd.Flags().String("name", "", "Token name")
	registerTokenCmd.Flags().String("metaData", "", "Meta data")
	registerTokenCmd.Flags().String("oracleInfo", "", "Oracle info")

	updateTokenCmd.Flags().String("rpcUrl", "http://localhost:8545", "Exocore RPC URL")
	updateTokenCmd.Flags().String("assetAddress", "", "Asset address")
	updateTokenCmd.Flags().String("metaData", "", "Meta data")

	registerOrUpdateClientChainCmd.Flags().String("rpcUrl", "http://localhost:8545", "Exocore RPC URL")
	registerOrUpdateClientChainCmd.Flags().Uint32("clientChainID", 0, "Client chain ID")
	registerOrUpdateClientChainCmd.Flags().Uint8("addressLength", 0, "Address length")
	registerOrUpdateClientChainCmd.Flags().String("name", "", "Name")
	registerOrUpdateClientChainCmd.Flags().String("metaInfo", "", "Meta info")
	registerOrUpdateClientChainCmd.Flags().String("signatureType", "", "Signature type")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

func deposit_(rpcUrl, stakerAddress string, amount *big.Int) error {
	depositAddr := common.HexToAddress(depositPrecompileAddress)
	assetAddr, err := assetToBytes(defaultAssetID)
	if err != nil {
		return err
	}
	// assetAddr := common.HexToAddress(defaultAssetID)
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

	data, err := depositAbi.Pack("depositLST", layerZeroID, assetAddr, paddingAddressTo32(stakerAddr), opAmount)
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
	assetAddr, err := assetToBytes(defaultAssetID)
	if err != nil {
		return err
	}
	// assetAddr := common.HexToAddress(defaultAssetID)
	stakerAddr := common.HexToAddress(stakerAddress)
	operatorAddr := []byte(operatorBench32Str)
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

	delegateAbi, err := abi.JSON(strings.NewReader(DelegateABI))
	if err != nil {
		return err
	}

	data, err := delegateAbi.Pack("delegate", layerZeroID, assetAddr, paddingAddressTo32(stakerAddr), operatorAddr, opAmount)
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

func undelegate_(rpcUrl, stakerAddress, operatorBench32Str string, amount *big.Int, instantUnbond bool) error {
	delegateAddr := common.HexToAddress(delegatePrecompileAddress)
	assetAddr, err := assetToBytes(defaultAssetID)
	if err != nil {
		return err
	}
	// assetAddr := common.HexToAddress(defaultAssetID)
	stakerAddr := common.HexToAddress(stakerAddress)
	operatorAddr := []byte(operatorBench32Str)
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

	delegateAbi, err := abi.JSON(strings.NewReader(DelegateABI))
	if err != nil {
		return err
	}

	data, err := delegateAbi.Pack("undelegate", layerZeroID, assetAddr, paddingAddressTo32(stakerAddr), operatorAddr, opAmount, instantUnbond)
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

func cancelSelfDelegate_(rpcUrl, stakerAddr string) error {
	delegateAddr := common.HexToAddress(delegatePrecompileAddress)

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

	data, err := delegateAbi.Pack("dissociateOperatorFromStaker", layerZeroID, staker)
	if err != nil {
		return err
	}

	txID, err := sendTransaction(ethClient, chainID, callAddr, sk, delegateAddr, data)
	if err != nil {
		return err
	}

	fmt.Println("Cancel Self Delegate Transaction ID:", txID)
	return waitForTransaction(ethClient, txID)
}

func withdrawLST_(rpcUrl, stakerAddress string, amount *big.Int) error {
	depositAddr := common.HexToAddress(depositPrecompileAddress)
	assetAddr, err := assetToBytes(defaultAssetID)
	if err != nil {
		return err
	}
	// assetAddr := common.HexToAddress(defaultAssetID)
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

	data, err := depositAbi.Pack("withdrawLST", layerZeroID, assetAddr, paddingAddressTo32(stakerAddr), opAmount)
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

func depositNST_(rpcUrl, pubkey string, stakerAddress string, amount *big.Int) error {
	depositAddr := common.HexToAddress(depositPrecompileAddress)
	stakerAddr := common.HexToAddress(stakerAddress)
	opAmount := amount
	if len(pubkey) != 64 {
		return fmt.Errorf("invalid pubkey length: %d", len(pubkey))
	}
	pubkeyBytes := common.Hex2Bytes(pubkey)
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

	data, err := depositAbi.Pack("depositNST", layerZeroID, pubkeyBytes, paddingAddressTo32(stakerAddr), opAmount)
	if err != nil {
		return err
	}

	txID, err := sendTransaction(ethClient, chainID, callAddr, sk, depositAddr, data)
	if err != nil {
		return err
	}

	fmt.Println("Deposit NST Transaction ID:", txID)
	return waitForTransaction(ethClient, txID)
}

func withdrawNST_(rpcUrl, pubkey string, stakerAddress string, amount *big.Int) error {
	depositAddr := common.HexToAddress(depositPrecompileAddress)
	stakerAddr := common.HexToAddress(stakerAddress)
	opAmount := amount
	if len(pubkey) != 64 {
		return fmt.Errorf("invalid pubkey length: %d", len(pubkey))
	}
	pubkeyBytes := common.Hex2Bytes(pubkey)
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

	data, err := depositAbi.Pack("withdrawNST", layerZeroID, pubkeyBytes, paddingAddressTo32(stakerAddr), opAmount)
	if err != nil {
		return err
	}

	txID, err := sendTransaction(ethClient, chainID, callAddr, sk, depositAddr, data)
	if err != nil {
		return err
	}

	fmt.Println("Withdraw NST Transaction ID:", txID)
	return waitForTransaction(ethClient, txID)
}

func registerToken_(rpcUrl, assetAddress string, decimals uint8, name string, metaData string, oracleInfo string) error {
	depositAddr := common.HexToAddress(depositPrecompileAddress)
	var token []byte
	if len(assetAddress) == 42 {
		assetAddr := common.HexToAddress(assetAddress)
		token = paddingAddressTo32(assetAddr)
	} else if len(assetAddress) == 66 {
		token = common.Hex2Bytes(strings.TrimPrefix(assetAddress, "0x"))
	} else {
		return fmt.Errorf("invalid asset address length: %d", len(assetAddress))
	}

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

	data, err := depositAbi.Pack("registerToken", layerZeroID, token, decimals, name, metaData, oracleInfo)
	if err != nil {
		return err
	}

	txID, err := sendTransaction(ethClient, chainID, callAddr, sk, depositAddr, data)
	if err != nil {
		return err
	}

	fmt.Println("RegisterToken Transaction ID:", txID)
	return waitForTransaction(ethClient, txID)
}

func updateToken_(rpcUrl, assetAddress string, metaData string) error {
	depositAddr := common.HexToAddress(depositPrecompileAddress)
	assetAddr := common.HexToAddress(assetAddress)
	token := paddingAddressTo32(assetAddr)

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

	data, err := depositAbi.Pack("updateToken", layerZeroID, token, metaData)
	if err != nil {
		return err
	}

	txID, err := sendTransaction(ethClient, chainID, callAddr, sk, depositAddr, data)
	if err != nil {
		return err
	}

	fmt.Println("updateToken Transaction ID:", txID)
	return waitForTransaction(ethClient, txID)
}

func registerOrUpdateClientChain_(rpcUrl string, clientChainID uint32, addressLength uint8, name string, metaInfo string, signatureType string) error {
	depositAddr := common.HexToAddress(depositPrecompileAddress)

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

	data, err := depositAbi.Pack("registerOrUpdateClientChain", clientChainID, addressLength, name, metaInfo, signatureType)
	if err != nil {
		return err
	}

	txID, err := sendTransaction(ethClient, chainID, callAddr, sk, depositAddr, data)
	if err != nil {
		return err
	}

	fmt.Println("registerOrUpdateClientChain Transaction ID:", txID)
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

func assetToBytes(assetAddress string) ([]byte, error) {
	if len(assetAddress) == 42 {
		assetAddr := common.HexToAddress(assetAddress)
		return paddingAddressTo32(assetAddr), nil
	} else if len(assetAddress) == 66 {
		return common.Hex2Bytes(strings.TrimPrefix(assetAddress, "0x")), nil
	}
	return nil, fmt.Errorf("invalid asset address length: %d", len(assetAddress))
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
