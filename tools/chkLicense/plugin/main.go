package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"plugin"
	"time"
)

var (
	l   string
	lib string
	p   string
)

func init() {
	flag.StringVar(&l, "l", "../cli", "license.dat directory path")
	flag.StringVar(&lib, "lib", "../linklib/liblicense.so", "libauth.so file path")
	flag.StringVar(&p, "p", "switch-directory-chain", "product name ")
	flag.Usage = usage
}

func usage() {
	flag.PrintDefaults()
}

var (
	ErrList = map[int]string{
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
	}
)

func main() {
	flag.Parse()
	if l == "" {
		fmt.Println("please input license.dat file path")
		return
	}

	if lib == "" {
		fmt.Println("please input liblicense.so file path")
		return
	}

	plugin, err := plugin.Open(lib)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//创建license对象
	{
		NewLicenseFunc, err := plugin.Lookup("NewLicense")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		retSult := NewLicenseFunc.(func(string, string, string) string)(l, p, "./license.log")
		if retSult != "" {
			fmt.Printf("创建License对象失败,%s", retSult)
			return
		}
		fmt.Println("创建License对象成功")
	}

	//验证license
	{
		VerifyLicenseFunc, err := plugin.Lookup("VerifyLicense")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		errCode := VerifyLicenseFunc.(func(string, string) int)(l, p)
		if errCode <= 0 {
			fmt.Printf("验证失败,%d, %s\n", errCode, ErrList[errCode])
			return
		}
		fmt.Println("验证成功")
	}

	//读取license配置文件
	{
		ReadLicneseFunc, err := plugin.Lookup("ReadLicnese")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		ret := ReadLicneseFunc.(func(string, string) string)(l, p)
		if ret == "FAIL" {
			fmt.Println(ret)
			return
		}

		kvs := new(map[string]interface{})
		if err := json.Unmarshal([]byte(ret), kvs); err != nil {
			fmt.Println(err.Error())
			return
		}

		for k, v := range *kvs {
			fmt.Printf("k:%s, v:%s\n", k, v.(string))
		}
	}

	//查询过期时间
	{
		GetExpireSecFunc, err := plugin.Lookup("GetExpireSec")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		willExpireSec := GetExpireSecFunc.(func(string, string) int64)(l, p)
		if willExpireSec == -1 {
			fmt.Println("fail")
			return
		}
		fmt.Println("GetExpireSecFunc, 剩余秒数:", willExpireSec, time.Unix(willExpireSec+time.Now().Unix(), 0).Format("2006-01-02 15:04:05"))
	}

	//销毁license对象
	{
		FreeLicenseFunc, err := plugin.Lookup("FreeLicense")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		retSult := FreeLicenseFunc.(func() string)()
		if retSult != "" {
			fmt.Printf("销毁License对象失败,%s", retSult)
			return
		}
		fmt.Println("销毁License对象成功")
	}
	fmt.Println()
}

/*
2020年1月10号
sudo date -s "01/10/2020 13:30:00"
*/

// func main() {
// 	flag.Parse()
// 	if l == "" {
// 		fmt.Println("please input license.dat file path")
// 		return
// 	}

// 	if lib == "" {
// 		fmt.Println("please input libplugin.so file path")
// 		return
// 	}

// 	plugin, err := plugin.Open(lib)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	//创建license对象
// 	{
// 		NewLicenseFunc, err := plugin.Lookup("NewLicense")
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return
// 		}

// 		retSult := NewLicenseFunc.(func(string, string, string) string)(l, p, "./license.log")
// 		if retSult != "" {
// 			fmt.Printf("创建License对象失败,%s", retSult)
// 			return
// 		}
// 		fmt.Println("创建License对象成功")
// 	}

// 	//验证license while-test
// 	{
// 		VerifyLicenseFunc, err := plugin.Lookup("VerifyLicense")
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return
// 		}
// 		for {
// 			time.Sleep(time.Second * 1)
// 			errCode := VerifyLicenseFunc.(func(string, string) int)(l, p)
// 			if errCode <= 0 {
// 				fmt.Printf("验证失败,%d, %s\n", errCode, ErrList[errCode])
// 				continue
// 			}
// 			fmt.Println("验证成功")
// 		}
// 	}

// 	// 读取license配置文件
// 	{
// 		ReadLicneseFunc, err := plugin.Lookup("ReadLicnese")
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return
// 		}
// 		for {
// 			time.Sleep(time.Second * 1)
// 			ret := ReadLicneseFunc.(func(string, string) string)(l, p)
// 			if ret == "FAIL" {
// 				fmt.Println(ret)
// 				continue
// 			}
// 			kvs := new(map[string]interface{})
// 			if err := json.Unmarshal([]byte(ret), kvs); err != nil {
// 				fmt.Println(err.Error())
// 				continue
// 			}

// 			for k, v := range *kvs {
// 				fmt.Printf("k:%s, v:%s\n", k, v.(string))
// 			}
// 		}
// 	}

// 	//查询过期时间
// 	{
// 		GetExpireSecFunc, err := plugin.Lookup("GetExpireSec")
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return
// 		}
// 		for {
// 			time.Sleep(time.Second * 1)
// 			willExpireSec := GetExpireSecFunc.(func(string, string) int64)(l, p)
// 			if willExpireSec == -1 {
// 				fmt.Println("FAIL")
// 				continue
// 			}
// 			fmt.Println(GetExpireSecFunc, willExpireSec)
// 		}
// 	}

// 	//销毁license对象
// 	{
// 		FreeLicenseFunc, err := plugin.Lookup("FreeLicense")
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return
// 		}

// 		retSult := FreeLicenseFunc.(func() string)()
// 		if retSult != "" {
// 			fmt.Printf("销毁License对象失败,%s", retSult)
// 			return
// 		}
// 		fmt.Println("销毁License对象成功")
// 	}
// 	fmt.Println()
// }

// /*
// 2020年1月10号
// sudo date -s "01/10/2020 13:30:00"
// */
