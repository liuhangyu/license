# switch license
```
1.
根目录make生成激活程序
licensemgr linux平台
licensemgr.exe win平台

2.获取机器码以及license.dat安装程序
进入license/tools/register 
运行make
生成register

3.链接库
进入license/tools/linklib 
运行make
生成
/plugin/liblicense.so  
/shared/liblicense.go  
/shared/liblicense.h

```

### licensemgr是生成激活码程序(目录位置licensemgr)
```
licensemgr是生成激活码程序,需要输入产品,过期时间,机器码
(按quit或q退出)
初次运行会在当前目录生成/data/products.json产品配置信息;

如需添加新的产品信息(产品名是"新产品1"和程序内产品名是"new-product"):
1.打开文件/data/products.json
"
[
{"ProductExplan":"新产品1","ProductName":"new-product"},
{"ProductExplan":"目录链","ProductName":"switch-directory-chain"},
{"ProductExplan":"数据交换平台","ProductName":"switch"},
{"ProductExplan":"数易通","ProductName":"tusdao-shuttle"}
]
"
2.保存,重新运行./licensemgr 就可以看见新加的产品

3.如果产品需设置自定义的KV配置,需要修改对应产品的xml模板文件.
例如目录链配置:
./data/switch-directory-chain.xml

<?xml version="1.0" encoding="UTF-8"?>
<options version="1">
    <option>
        <desc>示例配置1</desc>
        <key>key1</key>
        <val>val1</val>
    </option>
    <option>
        <desc>配置版本1</desc>
        <key>key1</key>
        <val>val1</val>
    </option>
</options>


```

### register获取机器码以及license.dat安装程序(目录位置licensemgr/tools/register)
```
register获取机器码以及license.dat安装程序,其中第一步获取机器码,下一步输入激活码
```


### lib链接库(目录位置licensemgr/tools/linklib)
```
plugin/liblicense.go是go的license链接库
shared/liblicense.go是java的链接库
```

### license测试验证程序(目录位置licensemgr/tools/chkLicense)
```
chkplugin go测试license.dat程序
输入参数如:
./chkplugin -l ../register  -lib ../linklib/plugin/liblicense.so -p switch-directory-chain

chkshard java测试license.dat程序
输入参数如:
export LD_LIBRARY_PATH=../linklib/shared
./chkshard  "../register"  "switch-directory-chain" "../linklib/shared/liblicense.so"
```


### go程序对接(目录链,switch)
```
第一步:
进入license/tools/linklib目录

第二步:
make plugin
生成go的plugin/liblicense.so插件

第三步:
在go程序中使用liblicense.so
请查看:
license/tools/chkLicense/plugin/main.go 
```

linux获取so版本:
在运行程序当前目录license.log文件中记录dll版本信息

### java,C程序对接(数易通)
```
第一步:
进入license/tools/linklib目录

第二步:
make shard
生成shard/libshared.h shard/libshared.so

第三步:
在C程序中使用shard/liblicense.so
请查看:
license/tools/chkLicense/shard/main.c

```

### 库中接口(具体数值类型根据各语言设定而不同)
```
liblicense.so

1.创建license对象
入参:
license.dat文件所在的文件夹(建议创建独立存放license的文件夹)
产品名
license log文件路径
出参:
成功""空字符
失败,则errMsg
string NewLicense(string,string,string)  

2.销毁license对象
入参:
无
出参:
成功""空字符
失败,则errMsg
string FreeLicense() 

3.验证license(验证签名)
入参:
license.dat文件所在的文件夹(建议创建独立存放license的文件夹)
产品名
出参:
成功返回1
失败返回errCode,且小于等于0
错误errCode:
0:   "uninitialized object",
-1:  "unknown error",
-2:  "dir does not exist",
-3:  "new watcher object failed",
-4:  "watcher add dir failed",
-5:  "failed to load public key",
-6:  "reading authorization file failed",
-7:  "decode authorization file failed",
-8:  "failed to verify signature",
-9:  "unmarshal license object failed",
-10: "failed to get machine code",
-11: "product name does not match",
-12: "license is expired",
-13: "license used before issued",
-14: "machine id does not match",

int VerifyLicense(string,string)


4.读取license.dat配置(不验证签名)
入参:
license.dat文件所在的文件夹(建议创建独立存放license的文件夹)
产品名
出参:
成功json格式license文件内容
失败"FALT"
string ReadLicnese(string,string)

5.获取过期时间(不验证签名)
入参:
license.dat文件所在的文件夹(建议创建独立存放license的文件夹)
产品名
出参:
0已过期
-1失败
>0 剩余时间(剩余未过期的秒数)
int64 GetExpireSec(string,string)
```