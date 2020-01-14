package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
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
	OneYearSeconds   = 31536000
	OneDaySeconds    = 86400
	OneMinuteSeconds = 60
	DataDir          = "data"
)

var (
	LicenseID                   = "LicenseID"
	ProductName                 = "ProductName"
	EndTime                     = "EndTime"
	kvMap       map[string]bool = map[string]bool{
		LicenseID:   true,
		ProductName: true,
		EndTime:     true,
	}
)

type Products struct {
	ProductExplan string
	ProductName   string
}

type AttrKV struct {
	Desc string
	Key  string
	Val  string
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

func SelectProduct() (string, string, error) {
	var (
		index         int
		err           error
		productExplan string
		productName   string
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
			return "", "", err
		}

		if index <= 0 || index > len(products) {
			return "", "", fmt.Errorf("invalid input, please input 1 ~ %d number", len(products))
		}

		fmt.Println()
		fmt.Printf("你选择的产品是: %s\n", products[index-1].ProductExplan)
		fmt.Println()
		productName = products[index-1].ProductName
		productExplan = products[index-1].ProductExplan
	}

	return productName, productExplan, nil
}

func InputExpiresTime() (int64, error) {
	var (
		err       error
		expiresAt int64
	)

	fmt.Printf("%s\n", "请输入过期时间,例如12天:12d (单位:天[d] 分钟[m] 秒[s] 年[y]):")
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

		if strings.HasSuffix(inputString, "y") {
			inputString = inputString[:len(inputString)-1]
			years, err := strconv.ParseInt(inputString, 10, 64)
			// days, err := strconv.ParseFloat(inputString, 64)
			if err != nil {
				return 0, err
			}

			if years <= 0 || years > 100 {
				return 0, fmt.Errorf("输入年数不能小于0,大于%d年", 100)
			}

			expiresAt = years * OneYearSeconds
			duration := time.Duration(expiresAt) * time.Second

			fmt.Println()
			fmt.Printf("过期年数: %d years, 过期日期:%s \n", years, time.Now().Add(duration).Format("2006-01-02 15:04:05"))
			fmt.Println()
		} else if strings.HasSuffix(inputString, "d") {
			inputString = inputString[:len(inputString)-1]
			days, err := strconv.ParseInt(inputString, 10, 64)
			// days, err := strconv.ParseFloat(inputString, 64)
			if err != nil {
				return 0, err
			}

			if days <= 0 || days > 100*356 {
				return 0, fmt.Errorf("输入天数不能小于0,大于%d天", 100*356)
			}

			expiresAt = int64(days * OneDaySeconds)
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

func ReadCustomKV(productName string) ([]AttrKV, error) {
	type Option struct {
		XMLName xml.Name `xml:"option"`
		Desc    string   `xml:"desc"`
		Key     string   `xml:"key"`
		Value   string   `xml:"val"`
	}

	type XMLProduct struct {
		XMLName xml.Name `xml:"options"`
		Version string   `xml:"version,attr"`
		Options []Option `xml:"option"`
	}

	filePath := filepath.Join(DataDir, strings.Join([]string{productName, ".xml"}, ""))
	if public.Exists(filePath) {
		var (
			attr   = XMLProduct{}
			attrKV []AttrKV
		)

		fd, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}
		defer fd.Close()

		attrKVBytes, err := ioutil.ReadAll(fd)
		if err != nil {
			return nil, err
		}

		err = xml.Unmarshal(attrKVBytes, &attr)
		if err != nil {
			return nil, err
		}

		// fmt.Printf("%s特性选择:\n", productExplan)
		for i := 0; i < len(attr.Options); i++ {
			// fmt.Println(i+1, attr.Options[i].Desc, attr.Options[i].Key, attr.Options[i].Value)
			attrKV = append(attrKV, AttrKV{Desc: attr.Options[i].Desc, Key: attr.Options[i].Key, Val: attr.Options[i].Value})
		}
		// fmt.Println("请输入数字序号,以分号间隔:")
		return attrKV, nil
	} else {
		fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		defer fd.Close()

		attr := &XMLProduct{
			Version: "1",
		}

		attr.Options = append(attr.Options, Option{
			Desc:  "示例配置1",
			Key:   "key1",
			Value: "val1",
		})

		output, err := xml.MarshalIndent(attr, "", "    ")
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		_, err = fd.Write([]byte(xml.Header))
		if err != nil {
			return nil, err
		}
		_, err = fd.Write(output)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func SelectCustomKV(productExplan string, kv []AttrKV) ([]AttrKV, error) {
	var (
		arrayIdx []int
		kvs      []AttrKV
	)

	if kv != nil {
		fmt.Printf("%s启用配置选择(请输入数字序号,以分号间隔,跳过按回车):\n", productExplan)
		for i := 0; i < len(kv); i++ {
			fmt.Println(i+1, kv[i].Desc)
		}

		input, err := inputReader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		defer inputReader.Reset(os.Stdin)

		inputString := strings.TrimSpace(input)

		if inputString != "" {
			if strings.HasPrefix(inputString, "q") {
				os.Exit(0)
			}

			arrayIndx := strings.Split(inputString, ",")
			for i := 0; i < len(arrayIndx); i++ {
				num := strings.TrimSpace(arrayIndx[i])
				if num == "" {
					continue
				}

				idx, err := strconv.Atoi(num)
				if err != nil {
					return nil, err
				}
				if idx <= 0 || idx > len(kv) {
					return nil, fmt.Errorf("输入不能小于等于0或大于%d", len(kv))
				}

				arrayIdx = append(arrayIdx, idx)
			}

			arrayIdx = public.RemoveDuplicate(arrayIdx)

			fmt.Printf("\n你选择的是%v,启用的配置是:\n", arrayIdx)
			for _, indx := range arrayIdx {
				fmt.Printf("%d %s\n", indx, kv[indx-1].Desc)
				kvs = append(kvs, kv[indx-1])
			}
			fmt.Println()
		}

		return kvs, nil
	}

	fmt.Println()
	return nil, nil
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
		err           error
		productName   string
		productExplan string
		expiresAt     int64
		machineID     string
		kv            []AttrKV
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
		productName, productExplan, err = SelectProduct()
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

	attrKV, err := ReadCustomKV(productName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		if len(attrKV) == 1 && attrKV[0].Desc == "示例配置1" {
			break
		}

		kv, err = SelectCustomKV(productExplan, attrKV)
		if err != nil {
			fmt.Println(err.Error())
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
	customKV := map[string]string{LicenseID: uuid, ProductName: productName, EndTime: expiresTime.Format(time.RFC3339)}

	for _, v := range kv {
		if _, ok := kvMap[v.Key]; ok {
			fmt.Printf("模板定义字段%s与系统定义字段冲突\n", v.Key)
			return
		}
		customKV[v.Key] = v.Val
	}

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
