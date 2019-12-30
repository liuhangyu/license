package public

import (
	"bufio"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func SaveLicensePem(dir string, filename string, licenseString string, licenseID string, productName string) error {
	isExist := Exists(dir)
	if isExist == false && dir != "." {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	savePath := filepath.Join(dir, filename)

	fd, err := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()

	w := bufio.NewWriter(fd)
	block := &pem.Block{
		Type:    "LICENSE",
		Headers: map[string]string{"LicenseID": licenseID, "ProductName": productName},
		Bytes:   []byte(licenseString),
	}

	if err := pem.Encode(w, block); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

func CheckPemAndSave(filePath string, licenseBytes []byte) (*pem.Block, error) {
	var (
		block *pem.Block
	)

	saveDir := filepath.Dir(filePath)
	fmt.Println(saveDir)

	isExist := Exists(saveDir)
	if isExist == false {
		err := os.MkdirAll(saveDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	if block, _ = pem.Decode(licenseBytes); block == nil {
		return nil, fmt.Errorf("%s", "license must be PEM")
	}

	fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	_, err = fd.Write(licenseBytes)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func ReadLicensePem(filePath string) ([]byte, error) {
	isExist := Exists(filePath)
	if isExist == false {
		return nil, fmt.Errorf("%s", "license file does not exist")
	}

	fd, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	licenseBytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	// block, _ := pem.Decode(licenseBytes)

	return licenseBytes, nil
}
