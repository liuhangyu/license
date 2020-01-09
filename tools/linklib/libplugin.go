package main

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
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
	logFile        *os.File
)

func init() {
	var err error
	logFile, err = os.OpenFile("./license.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		os.Exit(0)
	}
	errLog = log.New(logFile, "[plugin]", log.LstdFlags|log.Lshortfile|log.LstdFlags|log.Lmicroseconds)
	errLog.Println("license init...")
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

		verifyLicense := func(dir string, proName string, isVerifySign bool) ([]byte, error) {
			var (
				err      error
				content  []byte
				pemBlock *pem.Block
			)

			//构造验证签名对象
			alg, err := public.GetNonEquAlgorthm(nil, []byte(public.ECDSA_PUBLICKEY))
			if err != nil {
				return nil, err
			}

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

		for {
			//启动目录监听
			watcher, mErr = fsnotify.NewWatcher()
			if mErr != nil {
				break
			}

			//监控目录
			mErr = watcher.Add(dir)
			if mErr != nil {
				break
			}
			break
		}
		if mErr != nil {
			errLog.Println(mErr.Error())
			if watcher != nil {
				watcher.Close()
			}
			if logFile != nil {
				logFile.Close()
			}
			os.Exit(0)
			return
		}

		licenseBytes, err := verifyLicense(dir, productName, isVerifySign)
		if err != nil {
			mErr = err
			errLog.Println(mErr.Error())
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

					if strings.Contains(event.Name, LiceseFileName) == false {
						errLog.Printf("event, event:%v, name:%s\n", event, event.Name)
						continue
					} else {
						if event.Op&fsnotify.Create == fsnotify.Create { //no modified
							errLog.Printf("event modified file, event:%v, name:%s\n", event, event.Name)
						} else if event.Op&fsnotify.Write == fsnotify.Write {
							errLog.Printf("event write file, event:%v, name:%s\n", event, event.Name)
						} else if event.Op&fsnotify.Remove == fsnotify.Remove {
							errLog.Printf("event remove file, event:%v, name:%s\n", event, event.Name)
						} else {
							errLog.Printf("event contains file, event:%v, name:%s, continue...\n", event, event.Name)
							continue
						}
					}

					licenseBytes, err := verifyLicense(dir, productName, isVerifySign)
					if err != nil {
						mErr = err
						licenseContent = nil
						isValidLicense = false
					} else {
						licenseContent = licenseBytes
						isValidLicense = true
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						errLog.Println("close watcher")
						isValidLicense = false
						return
					}
					errLog.Println(err.Error())
				}
			}
		}()
	})

	return mErr
}

//export LicenseConfigVer
func LicenseConfigVer() string {
	fmt.Printf("license config version", LicenseConfigVerTag)
	return LicenseConfigVerTag
}

//export VerifyLicense
func VerifyLicense(licenseDir string, productName string) int {
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
func ReadLicnese(licenseDir string, productName string) string {
	var (
		pemBlock     *pem.Block
		lbytes       []byte
		err          error
		licenseBytes []byte
	)

	err = startWatcher(licenseDir, productName, true)
	if err != nil {
		errLog.Println(err.Error())
		return string("FAIL")
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
		break
	}

	if err != nil {
		errLog.Println(err.Error())
		return string("FAIL")
	}

	return string(lbytes)
}

//export GetExpireSec
func GetExpireSec(licenseDir string, productName string) int64 {
	var (
		pemBlock     *pem.Block
		err          error
		licenseBytes []byte
	)

	err = startWatcher(licenseDir, productName, true)
	if err != nil {
		errLog.Println(err.Error())
		return -1
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

		return license.GetExpiresAt()
	}
	return -1
}

func main() {
}

//https://studygolang.com/articles/13646
