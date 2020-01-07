package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"code.uni-ledger.com/switch/license/public"
)

const (
	OneDaySeconds = 86400
)

type Products struct {
	ProductExplan string
	ProductName   string
}

var (
	inputReader *bufio.Reader
	products    = []*Products{
		&Products{
			ProductExplan: "目录链",
			ProductName:   "switch-directory-chain",
		},

		&Products{
			ProductExplan: "数据交换平台",
			ProductName:   "switch",
		},

		&Products{
			ProductExplan: "数易通",
			ProductName:   "tusdao-shuttle",
		},
	}
)

func SelectProduct() (string, error) {
	var (
		index      int
		err        error
		retProName string
	)

	fmt.Printf("%s\n", "请选择需要激活的产品(输入数字):")
	for i := 1; i <= len(products); i++ {
		fmt.Printf("%d %s\n", i, products[i-1].ProductExplan)
	}

	input, err := inputReader.ReadString('\n')
	if err != nil {
		os.Exit(0)
	}
	defer inputReader.Reset(os.Stdin)

	inputString := strings.TrimSpace(input)

	if inputString != "" {
		if strings.HasPrefix(inputString, "q") {
			os.Exit(0)
		}

		index, err = strconv.Atoi(inputString)
		if err != nil {
			return "", err
		}

		if index <= 0 || index > len(products) {
			return "", fmt.Errorf("invalid input, please input 1 ~ %d number", len(products))
		}

		fmt.Println()
		fmt.Printf("你选择的产品是: %s\n", products[index-1].ProductExplan)
		fmt.Println()
		retProName = products[index-1].ProductName
	}

	return retProName, nil
}

func InputExpiresTime() (int64, error) {
	var (
		err       error
		expiresAt int64
	)

	fmt.Printf("%s\n", "请输入过期时间(单位天):")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	defer inputReader.Reset(os.Stdin)

	inputString := strings.TrimSpace(input)

	if inputString != "" {
		if strings.HasPrefix(inputString, "q") {
			os.Exit(0)
		}

		days, err := strconv.ParseInt(inputString, 10, 64)
		if err != nil {
			return 0, err
		}

		if days <= 0 || days > 100*356 {
			return 0, fmt.Errorf("%s", "invalid input number")
		}

		expiresAt = days * OneDaySeconds
		// expiresAt = days //test second
		duration := time.Duration(expiresAt) * time.Second
		fmt.Println()
		fmt.Printf("过期天数: %d days, 过期日期:%s \n", days, time.Now().Add(duration).Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	// fmt.Println("\033[H\033[2J")
	return expiresAt, nil
}

func InputMachineID() (string, error) {
	fmt.Printf("%s\n", "请输入机器ID:")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	defer inputReader.Reset(os.Stdin)

	inputString := strings.TrimSpace(input)

	if inputString != "" {
		if strings.HasPrefix(inputString, "q") {
			os.Exit(0)
		}

		fmt.Println()
		fmt.Printf("你输入的机器ID是: %s\n", inputString)
		fmt.Println()
	}

	return inputString, nil
}

func ShowActiveCode(dir, fileName, uuid string) {
	fmt.Printf("序号:%s \n", uuid)
	fmt.Printf("\n%s\n", "激活码是:")
	readPath := filepath.Join(dir, fileName)
	licenseActive, err := public.ReadLicensePem(readPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(base64.URLEncoding.EncodeToString(licenseActive))
	// fmt.Println(string(licenseActive))
}

func main() {
	var (
		err         error
		productName string
		expiresAt   int64
		machineID   string
	)

	flag.Parse()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	inputReader = bufio.NewReader(os.Stdin)

	for {
		productName, err = SelectProduct()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if productName == "" {
			continue
		}
		break
	}

	for {
		expiresAt, err = InputExpiresTime()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if expiresAt <= 0 {
			continue
		}
		break
	}

	for {
		machineID, err = InputMachineID()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if machineID == "" {
			continue
		}
		break
	}

	alg, err := public.GetNonEquAlgorthm([]byte(public.ECDSA_PRIVATE), []byte(public.ECDSA_PUBLICKEY))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	duration := time.Duration(expiresAt) * time.Second

	//构造license结构
	licenseIns := public.GenerateLicense(productName, machineID, duration)
	enCodeBytes, err := licenseIns.ToBytes()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//签名license
	licenseString, err := alg.SignedBytes(enCodeBytes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dir := "./db"
	// fileName := "license.dat"
	fileName := strings.Join([]string{"license", licenseIns.LicenseUUID, "dat"}, ".")
	err = public.SaveLicensePem(dir, fileName, licenseString, licenseIns.LicenseUUID, productName, licenseIns.GetEndTime())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ShowActiveCode(dir, fileName, licenseIns.LicenseUUID)
	<-c //quit
}
