1、编译linux安装包
```
make build-linux
```

2、制作docker镜像
```
docker build -t tygeth/muse:latest .
```

3、登录docker hub
```
docker login
```

4、推送docker镜像
```
docker push tygeth/muse:latest
```


5、启动4个docker节点

```
make localnet-start
```

更新可执行程序，只需要重新编译并重启节点:
```$xslt
make build-linux localnet-start
```


