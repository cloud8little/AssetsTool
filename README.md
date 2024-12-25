# AssetsTool

## Overview

The AssetsTool provide interface for interacting with exocore assets and delegation module, enable directly access to the precompile contract.

## Quick Start

### Steps

#### LST ReStaking

Example: 
1. ./deposit.sh
2. exocored q assets QueStakerAssetInfos 0xa53f68563D22EB0dAFAA871b6C08a6852f91d627_0x9ce1 --node http://localhost:20000
3. ./delegate.sh
4. exocored q assets QueOperatorAssetInfos exo1hj3qk6wg7se6l8g3s3ept7aas37dc75fk3lm2s --node http://localhost:20000
5. ./selfdelegate.sh
6. exocored q delegation QueryAssociatedOperatorByStaker 0xa53f68563D22EB0dAFAA871b6C08a6852f91d627_0x9ce1 --node http://localhost:20000

#### NST Restaking
1. ./depositNST.sh
2. exocored q oracle show-staker-list 0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee_0x9d19 --node http://localhost:20000
3. exocored q assets QueStakerAssetInfos 0xa53f68563D22EB0dAFAA871b6C08a6852f91d627_0x9d19 --node http://localhost:20000
4. ./delegateNST.sh
5. exocored q assets QueOperatorAssetInfos exo1hj3qk6wg7se6l8g3s3ept7aas37dc75fk3lm2s --node http://localhost:20000
6. ./undelegateNST.sh
7. exocored q assets QueOperatorAssetInfos exo1hj3qk6wg7se6l8g3s3ept7aas37dc75fk3lm2s --node http://localhost:20000


## License

This project is licensed under the MIT License.