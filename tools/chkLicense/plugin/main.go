package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"plugin"
)

var (
	l   string
	lib string
	p   string
)

func init() {
	flag.StringVar(&l, "l", "../cli", "license.dat directory path")
	flag.StringVar(&lib, "lib", "../linklib/libplugin.so", "libauth.so file path")
	flag.StringVar(&p, "p", "switch-directory-chain", "product name ")
	flag.Usage = usage
}

func usage() {
	flag.PrintDefaults()
}

type License struct {
	LicenseUUID string `json:"licensever"`  //license 唯一编号
	ConfigVer   uint32 `json:"configver"`   //配置版本
	ProductName string `json:"productname"` //产品名称
	MachineID   string `json:"machineid"`   //机器ID
	ExpiresAt   int64  `json:"expiresat"`   //过期时间
	IssuedAt    int64  `json:"issuedat"`    //签发时间
}

func main() {
	flag.Parse()
	if l == "" {
		fmt.Println("please input license.dat file path")
		return
	}

	if lib == "" {
		fmt.Println("please input libplugin.so file path")
		return
	}

	plugin, err := plugin.Open(lib)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//验证license
	{
		VerifyLicenseFunc, err := plugin.Lookup("VerifyLicense")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		ret := VerifyLicenseFunc.(func(string, string) int)(l, p)
		if ret == -1 {
			fmt.Println("验证失败")
			return
		}
		fmt.Println("验证成功")
	}

	// //验证license while-test
	// {
	// 	VerifyLicenseFunc, err := plugin.Lookup("VerifyLicense")
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		return
	// 	}

	// 	for {
	// 		time.Sleep(time.Second * 1)
	// 		ret := VerifyLicenseFunc.(func(string, string) int)(l, p)
	// 		if ret == -1 {
	// 			fmt.Println("验证失败")
	// 			continue
	// 		}
	// 		fmt.Println("验证成功")
	// 	}
	// }

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

		l := new(License)
		if err := json.Unmarshal([]byte(ret), l); err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("%+v", *l)
		fmt.Println()
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
		fmt.Println(GetExpireSecFunc, willExpireSec)
	}
}
