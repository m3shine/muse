rm -rf /$HOME/.mused
if [ ! -f ~/.mused/config/genesis.json ]; then
    echo 'Initialize the genesis.json file that will help you to bootstrap the network'
    `mused init musednode0 --chain-id musechain`

    echo 'Create a key to hold your validator account'
    `musecli keys add validator0`

    echo 'Add that key into the genesis.app_state.accounts array in the genesis file'
    `mused add-genesis-account $(musecli keys show validator0 -a) 1000musetoken,100000000stake`

    `musecli config chain-id musechain`
    `musecli config output json`
    `musecli config indent true`
    `musecli config trust-node true`
    `mused gentx --name validator0`

    echo ' Generate the transaction that creates your validator'
    `mused collect-gentxs`

    echo 'Add the generated bonding transaction to the genesis file'
    `mused validate-genesis`
fi
mused start