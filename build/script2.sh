mused init musednode2 --chain-id musechain
musecli keys add validator2
mused add-genesis-account $(musecli keys show validator2 -a) 1000musetoken,100000000stake
musecli config chain-id musechain
musecli config output json
musecli config indent true
musecli config trust-node true
mused gentx --name validator2
mused collect-gentxs
mused validate-genesis
sed -i '.bak' "s/persistent_peers = .*/persistent_peers = '1@192.168.10.3:26656'/g" ~/.mused/config/config.toml
mused start