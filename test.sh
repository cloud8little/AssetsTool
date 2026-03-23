identifier="minute"

# assetcli set-staker-reward-params --rpcUrl http://localhost:8545 --staker 0x884958ffb23885a13b5d51617f05f2b8b20f0405  --redelegateReward true --redelegateOperator im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5 --privateKey C26A874A75B028638D477DDF31EB8627899CB505798DF70D2DD2A631F9CAE7A4 --layerZeroID 40217
# echo "check staker reward params"

# imuad q feedistribution staker-reward-params 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19

echo "check current epoch number"

# get the epoch number
get_current_epoch() {
    curl -s http://localhost:1317/imuachain/epochs/v1/current_epoch/${identifier} | jq ".current_epoch" | sed 's/"//g'
}

current_epoch=$(get_current_epoch)
echo "Current epoch number: $current_epoch"
expected_epoch=2
if [ "$current_epoch" -lt "$expected_epoch" ]; then
    echo "Waiting for epoch to reach $expected_epoch..."
    while [ "$current_epoch" -lt "$expected_epoch" ]; do
        sleep 10
        current_epoch=$(get_current_epoch)
        echo "Current epoch number: $current_epoch"
    done
fi

echo "send a undelegate transaction"
assetcli undelegate --rpcUrl http://localhost:8545 --staker 0x884958ffb23885a13b5d51617f05f2b8b20f0405 --amount 1000000000000000000 --privateKey C26A874A75B028638D477DDF31EB8627899CB505798DF70D2DD2A631F9CAE7A4  --defaultAssetID 0x80e5bb3a04554e54b40dd6e14ca0f97212d9428d --layerZeroID 40217 --operator im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5
echo "check staker unclaimed rewards at epoch 2"
imuad q feedistribution staker-unclaimed-rewards 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19
echo "check staker claimed rewards at epoch 2"
imuad q feedistribution staker-claimed-rewards 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19 0xedb7a6077ab45df72e57bc2ea091f9183429720e

echo "wait for epoch 3"
expected_epoch=3
if [ "$current_epoch" -lt "$expected_epoch" ]; then
    echo "Waiting for epoch to reach $expected_epoch..."
    while [ "$current_epoch" -lt "$expected_epoch" ]; do
        sleep 10
        current_epoch=$(get_current_epoch)
        echo "Current epoch number: $current_epoch"
    done
fi

echo "check staker unclaimed rewards"
imuad q feedistribution staker-unclaimed-rewards 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19
echo "check staker claimed rewards"
imuad q feedistribution staker-claimed-rewards 0x884958ffb23885a13b5d51617f05f2b8b20f0405_0x9d19 0xedb7a6077ab45df72e57bc2ea091f9183429720e

echo "operator assets"
imuad q assets operator-assets im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5 
echo "operator unclaimed rewards"
imuad q feedistribution operator-unclaimed-rewards im13py43laj8zz6zw6a29sh7p0jhzeq7pq923hdp5 0xedb7a6077ab45df72e57bc2ea091f9183429720e
