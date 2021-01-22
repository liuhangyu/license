package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"license/public"
	"os"
	"strings"
)

//go build直接编译
var (
	inputReader *bufio.Reader
)

func ShowMachine() (string, error) {
	fmt.Printf("%s\n", "机器码是:")
	machineID, err := public.GetMachineID()
	if err != nil {
		return "", err
	}

	fmt.Println(machineID)
	fmt.Println()
	return machineID, nil
}

func ActiveProduct(alg *public.NonEquAlgorthm, savePath string, machineID string) error {
	var (
		err error
	)

	fmt.Printf("%s\n", "请输入激活码:")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return err
	}
	defer inputReader.Reset(os.Stdin)

	inputString := strings.TrimSpace(input)

	if inputString != "" {
		fmt.Println()
		fmt.Printf("你输入的激活码是: %s\n", inputString)
		fmt.Println()

		activeCodeBytes, err := base64.URLEncoding.DecodeString(inputString)
		if err != nil {
			return err
		}

		//验证激活码格式
		block, err := public.CheckPemAndSave(savePath, activeCodeBytes)
		if err != nil {
			return err
		}

		l, err := public.BytesToLicense(string(block.Bytes))
		if err != nil {
			return err
		}

		//验证机密是否对应
		if l.MachineID != machineID {
			return fmt.Errorf("%s", "输入激活码无效,无法匹配机器,请重新输入")
		}

		//验证签名
		_, err = alg.VerifySign(string(block.Bytes))
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("%s", "cannot enter empty activation code")
}

var (
	o string
)

func init() {
	flag.StringVar(&o, "o", "./license.dat", "save activation code file path")
	flag.Usage = usage
}

func usage() {
	flag.PrintDefaults()
	fmt.Println(public.GetAppInfo())
}

func main() {
	var (
		err       error
		machineID string
	)

	flag.Parse()

	//验证签名实例
	alg, err := public.GetNonEquAlgorthm(nil, []byte(public.ECDSA_PUBLICKEY))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	inputReader = bufio.NewReader(os.Stdin)

	//获取机器ID
	machineID, err = ShowMachine()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		//激活程序
		err := ActiveProduct(alg, o, machineID)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println("激活成功")
		break
	}
}
