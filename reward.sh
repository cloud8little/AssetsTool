./assetcli set-staker-reward-params --rpcUrl http://localhost:8545 --staker 0x884958ffb23885a13b5d51617f05f2b8b20f0405  --redelegateReward true --redelegateOperator im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5 --privateKey C26A874A75B028638D477DDF31EB8627899CB505798DF70D2DD2A631F9CAE7A4 --layerZeroID 40217
echo "check staker reward params"

imuad q feedistribution staker-reward-params 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19

imuad q feedistribution staker-unclaimed-rewards 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19

echo "then undelegate 1 erc20..."

./assetcli undelegate --rpcUrl http://localhost:8545 --staker 0x884958ffb23885a13b5d51617f05f2b8b20f0405 --amount 1000000000000000000 --privateKey C26A874A75B028638D477DDF31EB8627899CB505798DF70D2DD2A631F9CAE7A4  --defaultAssetID 0x80e5bb3a04554e54b40dd6e14ca0f97212d9428d --layerZeroID 40217 --operator im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5 
echo "check epoch number"

watch "imuad q epochs current-epoch minute"

imuad q feedistribution  staker-claimed-rewards 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19 0xedb7a6077ab45df72e57bc2ea091f9183429720e

imuad q assets operator-assets im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5

imuad q feedistribution reward-assets-by-avs 0xedb7a6077ab45df72e57bc2ea091f9183429720e

imuad q feedistribution operator-unclaimed-rewards im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5 0xedb7a6077ab45df72e57bc2ea091f9183429720e

imuad q feedistribution delegation-starting-info 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19 0x80e5bb3a04554e54b40dd6e14ca0f97212d9428d_0x9d19 im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5 minute

imuad q tendermint-validator-set

# instantUnbonding 

./assetcli undelegate --rpcUrl http://localhost:8545 --staker 0x884958ffb23885a13b5d51617f05f2b8b20f0405 --amount 1499000000000000000000 --privateKey C26A874A75B028638D477DDF31EB8627899CB505798DF70D2DD2A631F9CAE7A4  --defaultAssetID 0x80e5bb3a04554e54b40dd6e14ca0f97212d9428d --layerZeroID 40217 --operator im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5 --instantUnbond true


imuad q delegation delegations 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19 0x0000000000000000000000000000000000000000_0x0 

imuad q delegation delegations 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19 0x80e5bb3a04554e54b40dd6e14ca0f97212d9428d_0x9d19

imuad q feedistribution delegation-starting-info 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19 0x80e5bb3a04554e54b40dd6e14ca0f97212d9428d_0x9d19 im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5 minute

imuad q feedistribution delegation-starting-info 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x0 0x0000000000000000000000000000000000000000_0x0 im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5 minute

./assetcli undelegate --rpcUrl http://localhost:8545 --staker 0x884958ffb23885a13b5d51617f05f2b8b20f0405 --amount 200000000000000000000 --privateKey C26A874A75B028638D477DDF31EB8627899CB505798DF70D2DD2A631F9CAE7A4  --defaultAssetID 0x80e5bb3a04554e54b40dd6e14ca0f97212d9428d --layerZeroID 40217 --operator im192mgay9r8mdgwcprcvgdm47td36n974h4ay6c9 --instantUnbond true

./assetcli undelegate --rpcUrl http://localhost:8545 --staker 0x884958ffb23885a13b5d51617f05f2b8b20f0405 --amount 300000000000000000000 --privateKey C26A874A75B028638D477DDF31EB8627899CB505798DF70D2DD2A631F9CAE7A4  --defaultAssetID 0x80e5bb3a04554e54b40dd6e14ca0f97212d9428d --layerZeroID 40217 --operator im1m4mewv4fr2qphdjuft8cvlz4vphur93s3347e2 --instantUnbond true

imuad q assets staker-assets 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19  

# let IM token has voting power

assetcli register-token --rpcUrl http://localhost:8545 --privateKey C26A874A75B028638D477DDF31EB8627899CB505798DF70D2DD2A631F9CAE7A4 --assetAddress 0x0000000000000000000000000000000000000000 --decimals 18 --name ETH --metaData "ETH for Ethereum" --oracleInfo "ETH,Ethereum,8" --layerZeroID 0


