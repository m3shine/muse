#单节点部署
```
#编译安装
make install

# Initialize configuration files and genesis file
mused init --chain-id musenetwork

# Copy the `Address` output here and save it for later use 
# [optional] add "--ledger" at the end to use a Ledger Nano S 
musecli keys add jon


# Add both accounts, with coins to the genesis file
mused add-genesis-account $(musecli keys show jon -a) 1000nametoken,1000jackcoin

# Configure your CLI to eliminate need for chain-id flag
musecli config chain-id namechain
musecli config output json
musecli config indent true
musecli config trust-node true

#启动节点
mused start

#启动rest服务
musecli rest-server --chain-id musenetwork --trust-node





###################################################

命令行操作

#歌词上链
musecli tx muse set-lyric 110rr jon bigbiggirl 35465456 btoken --from jon
#查询
musecli query muse lyric 110rr

###################################################

http请求


#查询账户信息
curl -s http://localhost:1317/auth/accounts/$(musecli keys show jon -a)

#查询上链歌词
curl -s localhost:1317/muse/lyric/ufrug33q3iqd

#歌词上链
curl -XPUT -s http://localhost:1317/muse/lyric --data-binary '{"base_req":{"from":"cosmos193jutxkx74xx8yaufcx9pcp3cwd90nwsegklpa","password":"xxxx","chain_id":"musenetwork","sequence":"3","account_number":"0"},"lyric_code":"d0911","author":"jon","title":"helloworld","hash":"333333333333333333333333","owner":"cosmos193jutxkx74xx8yaufcx9pcp3cwd90nwsegklpa","token_name":"hwtoken"}'
```





#docker部署
未成功

```
#编译linux安装包
make build-linux

#制作docker镜像
docker build -t tygeth/muse:latest .

#登录docker hub
docker login

#推送docker镜像
docker push tygeth/muse:latest

#启动4个docker节点
make localnet-start

#更新可执行程序，只需要重新编译并重启节点:
make build-linux localnet-start
```


