# switch license
```
1.
根目录make生成激活程序
licensemgr linux平台
licensemgr.exe win平台

2.获取机器码以及license.dat安装程序
进入license/tools/cli 
go build
生成cli

3.链接库
进入license/tools/linklib 
运行make
生成libplugin.so  libshared.go  libshared.h

```

### licensemgr是生成激活码程序(目录位置licensemgr)
```
licensemgr是生成激活码程序,需要输入产品,过期时间,机器码
```

### cli获取机器码以及license.dat安装程序(目录位置licensemgr/tools/cli)
```
cli获取机器码以及license.dat安装程序,其中第一步获取机器码,下一步输入激活码
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

### go程序对接(目录链,switch)
```
第一步:
进入license/tools/linklib目录

第二步:
make plugin
生成libplugin.so插件

第三步:
在go程序中使用libplugin.so
请查看:
license/tools/chkLicense/plugin/main.go 
```


### java,C程序对接(数易通)
```
第一步:
进入license/tools/linklib目录

第二步:
make shard
生成libshared.h libshared.so

第三步:
在C程序中使用libshared.so
请查看:
license/tools/chkLicense/shard/main.c

```

### 库中接口(具体数值类型根据各语言设定而不同)
```
libplugin.so
libshared.so


1.验证license
入参:
license.dat文件所在的文件夹(建议创建独立存放license的文件夹)
产品名
出参:
成功0
失败-1
int VerifyLicense(string,string)


2.读取license.dat配置
入参:
license.dat文件所在的文件夹(建议创建独立存放license的文件夹)
产品名
出参:
成功json格式license文件内容
失败"FALT"
string ReadLicnese(string)

3.获取过期时间
入参:
license.dat文件所在的文件夹(建议创建独立存放license的文件夹)
产品名
出参:
0已过期
-1失败
>0 剩余时间(剩余未过期的秒数)
int64 GetExpireSec(string)
```