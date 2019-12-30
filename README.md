# switch license

### licensemgr是生成激活码程序(目录位置licensemgr)
```
licensemgr是生成激活码程序,需要输入产品,过期时间,机器码
```

### cli(目录位置licensemgr/tools/cli)
```
cli是安装程序,其中第一步获取机器码,下一步输入激活码
```


### lib链接库(目录位置licensemgr/tools/linklib)
```
libplugin.go是go的license链接库
libshared.go是java的链接库
```

### license测试验证程序(目录位置licensemgr/tools/chkLicense)
```
chkplugin go测试license.dat程序
输入参数如:
./chkplugin -l ../cli/license.dat -lib ../linklib/libplugin.so -p switch-directory-chain

chkshard java测试license.dat程序
输入参数如:
./chkshard  "../cli/license.dat"  "switch-directory-chain" "../linklib/libshared.so"
```