package main

import (
	"flag"
	"fmt"
	"plugin"
)

var (
	l   string
	lib string
	p string
)

func init() {
	flag.StringVar(&l, "l", "../../cli/license.dat", "license.dat file path")
	flag.StringVar(&lib, "lib", "../../linklib/libplugin.so", "libauth.so file path")
	flag.StringVar(&p, "p", "switch-directory-chain", "product name ")
	flag.Usage = usage
}

func usage() {
	flag.PrintDefaults()
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

		ret := VerifyLicenseFunc.(func(string, string) string)(l, p)
		if ret != "OK" {
			fmt.Println(ret)
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

		ret := ReadLicneseFunc.(func(string) string)(l)
		if ret == "FAIL" {
			fmt.Println(ret)
			return
		}
		fmt.Println(ret)
	}

	//查询过期时间
	{
		GetExpireSecFunc, err := plugin.Lookup("GetExpireSec")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		willExpireSec := GetExpireSecFunc.(func(string) int64)(l)
		if willExpireSec == -1 {
			fmt.Println("fail")
			return
		}
		fmt.Println(willExpireSec)
	}
}
