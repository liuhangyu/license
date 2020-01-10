package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"code.uni-ledger.com/switch/license/public"
)

const (
	OneDaySeconds    = 86400
	OneMinuteSeconds = 60
	DataDir          = "data"
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

	fmt.Printf("%s\n", "请输入过期时间,例如: 12d (单位:天[d] 分钟[m] 秒[s]):")
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

		if strings.HasSuffix(inputString, "d") {
			inputString = inputString[:len(inputString)-1]
			days, err := strconv.ParseInt(inputString, 10, 64)
			if err != nil {
				return 0, err
			}

			if days <= 0 || days > 100*356 {
				return 0, fmt.Errorf("输入天数不能小于0,大于%d天", 100*356)
			}

			expiresAt = days * OneDaySeconds
			duration := time.Duration(expiresAt) * time.Second

			fmt.Println()
			fmt.Printf("过期天数: %d days, 过期日期:%s \n", days, time.Now().Add(duration).Format("2006-01-02 15:04:05"))
			fmt.Println()
		} else if strings.HasSuffix(inputString, "m") {
			inputString = inputString[:len(inputString)-1]
			minute, err := strconv.ParseInt(inputString, 10, 64)
			if err != nil {
				return 0, err
			}

			if minute <= 0 || minute > 100*356*24*60 {
				return 0, fmt.Errorf("输入分钟不能小于0,大于%d分钟", 100*356*24*60)
			}

			expiresAt = minute * OneMinuteSeconds
			duration := time.Duration(expiresAt) * time.Second

			fmt.Println()
			fmt.Printf("过期分钟数: %d minute, 过期日期:%s \n", minute, time.Now().Add(duration).Format("2006-01-02 15:04:05"))
			fmt.Println()
		} else if strings.HasSuffix(inputString, "s") {
			inputString = inputString[:len(inputString)-1]
			seconds, err := strconv.ParseInt(inputString, 10, 64)
			if err != nil {
				return 0, err
			}

			if seconds <= 0 || seconds > 100*356*24*60*60 {
				return 0, fmt.Errorf("输入秒不能小于0,大于%d秒", 100*356*24*60*60)
			}

			expiresAt = seconds
			duration := time.Duration(expiresAt) * time.Second

			fmt.Println()
			fmt.Printf("过期秒数: %d minute, 过期日期:%s \n", seconds, time.Now().Add(duration).Format("2006-01-02 15:04:05"))
			fmt.Println()
		} else {
			return 0, fmt.Errorf("%s", "输入不正确,请输入时间单位")
		}
	}
	// fmt.Println("\033[H\033[2J")
	return expiresAt, nil
}

func InputMachineID() (string, error) {
	fmt.Printf("%s\n", "请输入机器码:")
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
		fmt.Printf("你输入的机器码是: %s\n", inputString)
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

func IsQuit() bool {
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return false
	}
	defer inputReader.Reset(os.Stdin)

	inputString := strings.TrimSpace(input)

	if inputString != "" {
		if strings.HasPrefix(inputString, "q") {
			os.Exit(0)
		}
	}
	return true
}

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Println("input 'quit' or 'q' to exit the program")
	fmt.Println(public.GetAppInfo())
}

func LoadConfig() ([]*Products, error) {
	var (
		productList = []*Products{}
	)

	isExist := public.Exists(DataDir)
	if isExist == false && DataDir != "." {
		err := os.MkdirAll(DataDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	filePath := filepath.Join(DataDir, "products.json")

	if public.Exists(filePath) { //文件存在
		fd, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}
		defer fd.Close()

		configBytes, err := ioutil.ReadAll(fd)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(configBytes, &productList)
		if err != nil {
			return nil, err
		}
	} else {
		encByte, err := json.Marshal(products)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(filePath, encByte, 0644)
		if err != nil {
			return nil, err
		}
		productList = products
	}

	return productList, nil
}

func main() {
	var (
		err         error
		productName string
		expiresAt   int64
		machineID   string
	)
	flag.Parse()

	//load config file
	products, err = LoadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

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

	//定义License HEAD KV
	uuid := public.GetUUID()
	expiresTime := time.Now().Add(duration)
	customKV := map[string]string{"LicenseID": uuid, "ProductName": productName, "EndTime": expiresTime.Format(time.RFC3339)}

	//构造license结构
	licenseIns := public.GenerateLicense(uuid, productName, machineID, expiresTime.Unix(), customKV)
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

	dir := filepath.Join(DataDir, "db")
	fileName := strings.Join([]string{"license", licenseIns.LicenseUUID, "dat"}, ".")
	err = public.SaveLicensePem(dir, fileName, licenseString, customKV)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ShowActiveCode(dir, fileName, licenseIns.LicenseUUID)

	for {
		IsQuit()
	}
}
