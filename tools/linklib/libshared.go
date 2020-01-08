package main

import "C"
import (
	"encoding/json"
	"encoding/pem"
	"strings"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"code.uni-ledger.com/switch/license/public"
	"code.uni-ledger.com/switch/license/public/deplib/fsnotify"
)

const (
	LinkLibVersion      = 1            //eq public.CONFIGVERSION
	LicenseConfigVerTag = "LiConfigV1" //LiConfigV1,LiConfigV2,...  #strings -a libplugin.so | grep LiConfig #查看版本
	LiceseFileName      = "license.dat"
)

var (
	once           sync.Once
	licenseContent []byte
	isValidLicense bool
	errLog         *log.Logger
)

func init() {
	logFile, err := os.OpenFile("./license.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		os.Exit(0)
	}
	errLog = log.New(logFile, "[shared]", log.LstdFlags|log.Lshortfile|log.LstdFlags)
	return
}

func startWatcher(dir string, productName string, isVerifySign bool) error {
	var (
		mErr error
	)

	once.Do(func() {
		var (
			watcher *fsnotify.Watcher
		)

		//启动目录监听
		watcher, mErr = fsnotify.NewWatcher()
		if mErr != nil {
			errLog.Println(mErr.Error())
			return
		}

		//构造验证签名对象
		alg, err := public.GetNonEquAlgorthm(nil, []byte(public.ECDSA_PUBLICKEY))
		if err != nil {
			mErr = err
			errLog.Println(mErr.Error())
			return
		}

		verifyLicense := func(dir string, alg *public.NonEquAlgorthm, proName string, isVerifySign bool) ([]byte, error) {
			var (
				err      error
				content  []byte
				pemBlock *pem.Block
			)

			//读取license.dat文件
			licenseFilePath := path.Join(dir, LiceseFileName)
			licenseBytes, err := public.ReadLicensePem(licenseFilePath)
			if err != nil {
				return nil, err
			}

			if isVerifySign {
				//验证pem格式
				if pemBlock, _ = pem.Decode(licenseBytes); pemBlock == nil {
					return nil, fmt.Errorf("%s", "license must be PEM")
				}

				//验证签名
				content, err = alg.VerifySign(string(pemBlock.Bytes))
				if err != nil {
					return nil, err
				}

				//获取license对象
				license, err := public.ToLicense(content)
				if err != nil {
					return nil, err
				}

				//获取机器ID
				machineID, err := public.GetMachineID()
				if err != nil {
					return nil, err
				}

				//对比过期时间以及机器ID
				err = license.Valid(LinkLibVersion, proName, machineID)
				if err != nil {
					return nil, err
				}
			}

			return licenseBytes, nil
		}

		licenseBytes, err := verifyLicense(dir, alg, productName, isVerifySign)
		if err != nil {
			mErr = err
			errLog.Println(mErr.Error())
			return
		} else {
			licenseContent = licenseBytes
			isValidLicense = true
		}

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}

					if strings.HasSuffix(event.Name, LiceseFileName) == false {
						errLog.Println("event:", event, event.Name)
						continue
					}

					licenseBytes, err := verifyLicense(dir, alg, productName, isVerifySign)
					if err != nil {
						mErr = err
						errLog.Println(mErr.Error())
						return
					} else {
						licenseContent = licenseBytes
						isValidLicense = true
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						errLog.Println("close watcher")
						return
					}
					errLog.Println(err.Error())
				}
			}
		}()

		//监控目录
		mErr = watcher.Add(dir)
		if mErr != nil {
			errLog.Println(mErr.Error())
			return
		}
	})

	if mErr != nil {
		isValidLicense = false
	} else {
		isValidLicense = true
	}
	return mErr
}

//export LicenseConfigVer
func LicenseConfigVer() string {
	fmt.Printf("license config version", LicenseConfigVerTag)
	return LicenseConfigVerTag
}

//export VerifyLicense
func VerifyLicense(licenseDir string, productName string) C.int {
	mErr := startWatcher(licenseDir, productName, true)
	if mErr != nil {
		errLog.Println(mErr.Error())
		return -1
	}

	if isValidLicense {
		return 0
	}

	return -1
}

//export ReadLicnese
func ReadLicnese(licenseDir string, productName string) *C.char {
	var (
		pemBlock     *pem.Block
		lbytes       []byte
		err          error
		licenseBytes []byte
		licenseStr   string
	)

	err = startWatcher(licenseDir, productName, false)
	if err != nil {
		errLog.Println(err.Error())
		return C.CString("FAIL")
	}

	for {
		if len(licenseContent) == 0 {
			//读取license.dat文件
			licenseFilePath := path.Join(licenseDir, LiceseFileName)
			licenseBytes, err = public.ReadLicensePem(licenseFilePath)
			if err != nil {
				break
			}
			licenseContent = licenseBytes
		} else {
			licenseBytes = licenseContent
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
		errLog.Println(err.Error())
		return C.CString("FAIL")
	}

	return C.CString(licenseStr)
}

//export GetExpireSec
func GetExpireSec(licenseDir string, productName string) C.longlong {
	var (
		pemBlock     *pem.Block
		err          error
		licenseBytes []byte
		seconds      int64
	)

	err = startWatcher(licenseDir, productName, false)
	if err != nil {
		errLog.Println(err.Error())
		return C.longlong(-1)
	}

	for {
		if len(licenseContent) == 0 {
			//读取license.dat文件
			licenseFilePath := path.Join(licenseDir, LiceseFileName)
			licenseBytes, err = public.ReadLicensePem(licenseFilePath)
			if err != nil {
				errLog.Println(err.Error())
				break
			}
			licenseContent = licenseBytes
		} else {
			licenseBytes = licenseContent
		}

		//验证pem格式
		if pemBlock, _ = pem.Decode(licenseBytes); pemBlock == nil {
			err = fmt.Errorf("%s", "license must be PEM")
			errLog.Println(err.Error())
			break
		}

		license, err := public.BytesToLicense(string(pemBlock.Bytes))
		if err != nil {
			errLog.Println(err.Error())
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
