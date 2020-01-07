package main

import "C"
import (
	"code.uni-ledger.com/switch/license/public"
	"encoding/json"
	"encoding/pem"
	"fmt"
)

const (
	LinkLibVersion      = 1            //eq public.CONFIGVERSION
	LicenseConfigVerTag = "LiConfigV1" //strings -a libplugin.so | grep LiConfig #查看版本
)

//export LicenseConfigVer
func LicenseConfigVer() string {
	fmt.Printf("license config version", LicenseConfigVerTag)
	return LicenseConfigVerTag
}

//export VerifyLicense
func VerifyLicense(licenseFilePath string, productName string) string {
	var (
		pemBlock *pem.Block
		mErr     error
	)

	for {
		//构造验证签名对象
		alg, err := public.GetNonEquAlgorthm(nil, []byte(public.ECDSA_PUBLICKEY))
		if err != nil {
			mErr = err
			break
		}

		//读取license.dat文件
		licenseBytes, err := public.ReadLicensePem(licenseFilePath)
		if err != nil {
			mErr = err
			break
		}

		//验证pem格式
		if pemBlock, _ = pem.Decode(licenseBytes); pemBlock == nil {
			mErr = fmt.Errorf("%s", "license must be PEM")
			break
		}

		//验证签名
		licenseBytes, err = alg.VerifySign(string(pemBlock.Bytes))
		if err != nil {
			mErr = err
			break
		}

		//获取license对象
		license, err := public.ToLicense(licenseBytes)
		if err != nil {
			mErr = err
			break
		}

		//获取机器ID
		machineID, err := public.GetMachineID()
		if err != nil {
			mErr = err
			break
		}

		//对比过期时间以及机器ID
		err = license.Valid(LinkLibVersion, productName, machineID)
		if err != nil {
			mErr = err
			break
		}

		break
	}

	if mErr != nil {
		return mErr.Error()
	}

	return string("OK")
}

//export ReadLicnese
func ReadLicnese(licenseFilePath string) string {
	var (
		pemBlock     *pem.Block
		lbytes       []byte
		err          error
		licenseBytes []byte
	)

	for {
		//读取license.dat文件
		licenseBytes, err = public.ReadLicensePem(licenseFilePath)
		if err != nil {
			break
		}

		//验证pem格式
		if pemBlock, _ = pem.Decode(licenseBytes); pemBlock == nil {
			err = fmt.Errorf("%s", "license must be PEM")
			break
		}

		license, err := public.BytesToLicense(string(pemBlock.Bytes))
		if err != nil {
			break
		}

		lbytes, err = json.Marshal(license)
		if err != nil {
			break
		}
		break
	}

	if err != nil {
		return string("FAIL")
	}

	return string(lbytes)
}

//export GetExpireSec
func GetExpireSec(licenseFilePath string) int64 {
	var (
		pemBlock     *pem.Block
		err          error
		licenseBytes []byte
	)

	for {
		//读取license.dat文件
		licenseBytes, err = public.ReadLicensePem(licenseFilePath)
		if err != nil {
			break
		}

		//验证pem格式
		if pemBlock, _ = pem.Decode(licenseBytes); pemBlock == nil {
			err = fmt.Errorf("%s", "license must be PEM")
			break
		}

		license, err := public.BytesToLicense(string(pemBlock.Bytes))
		if err != nil {
			break
		}

		return license.GetExpiresAt()
	}
	return -1
}

func main() {
	return
}

//https://studygolang.com/articles/13646
