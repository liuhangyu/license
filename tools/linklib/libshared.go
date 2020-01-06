package main

import "C"
import (
	"code.uni-ledger.com/switch/license/public"
	"encoding/json"
	"encoding/pem"
	"fmt"
)

const (
	LinkLibVersion = 1 //eq public.CONFIGVERSION
)

//export VerifyLicense
func VerifyLicense(licenseFilePath string, productName string) *C.char {
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
		return C.CString(mErr.Error())
	}

	return C.CString("OK")
}

//export ReadLicnese
func ReadLicnese(licenseFilePath string) *C.char {
	var (
		pemBlock     *pem.Block
		lbytes       []byte
		err          error
		licenseBytes []byte
		licenseStr   string
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
		licenseStr = string(lbytes)
		break
	}

	if err != nil {
		return C.CString("FAIL")
	}

	return C.CString(licenseStr)
}

//export GetExpireSec
func GetExpireSec(licenseFilePath string) C.longlong {
	var (
		pemBlock     *pem.Block
		err          error
		licenseBytes []byte
		seconds      int64
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

		seconds = license.GetExpiresAt()
		return C.longlong(seconds)
	}
	return C.longlong(-1)
}

func main() {

}

//https://studygolang.com/articles/13646
