package main

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"code.uni-ledger.com/switch/license/public"
	"code.uni-ledger.com/switch/license/public/deplib/fsnotify"
)

const (
	LiceseFileName = "license.dat"
)

type License struct {
	Once           sync.Once
	LicenseContent []byte
	IsValidLicense bool
	ErrLog         *log.Logger
	LogFile        *os.File
	Watcher        *fsnotify.Watcher
	Error          *public.ErrorMsg
}

var (
	gLicenseIns *License
)

func NewLicense(dir string, productName string, logPath string) string {
	var (
		mErr        error
		mLicenseIns = &License{}
	)

	verifyLicense := func(dir string, proName string) ([]byte, error) {
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

		{
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

	saveDir := filepath.Dir(logPath)
	isExist := public.Exists(saveDir)
	if isExist == false {
		mErr = os.MkdirAll(saveDir, os.ModePerm)
		if mErr != nil {
			return mErr.Error()
		}
	}

	mLicenseIns.LogFile, mErr = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if mErr != nil {
		return mErr.Error()
	}

	//初始化日志
	mLicenseIns.ErrLog = log.New(logFile, "[license]", log.LstdFlags|log.Lshortfile|log.LstdFlags|log.Lmicroseconds)
	mLicenseIns.ErrLog.Println("License", public.GetAppInfo())

	mLicenseIns.Once.Do(func() {
		var (
			err error
		)

		for {
			if public.Exists(dir) == false {
				mErr = public.New(-1, fmt.Sprintf("%s dir does not exist", dir))
				break
			}

			//启动目录监听
			mLicenseIns.Watcher, err = fsnotify.NewWatcher()
			if err != nil {
				mErr = public.New(-2, fmt.Sprintf("new watcher err %s", err.Error()))
				break
			}

			//监控目录
			err = mLicenseIns.Watcher.Add(dir)
			if err != nil {
				mErr = public.New(-3, fmt.Sprintf("watcher add dir err %s", err.Error()))
				break
			}
			break
		}
		if mErr != nil {
			mLicenseIns.ErrLog.Println(mErr.Error())
			if mLicenseIns.Watcher != nil {
				mLicenseIns.Watcher.Close()
				mLicenseIns.Watcher = nil
			}
			if mLicenseIns.LogFile != nil {
				mLicenseIns.LogFile.Close()
				mLicenseIns.LogFile = nil
			}
			return
		}

		licenseBytes, err := verifyLicense(dir, productName)
		if err != nil {
			mErr = public.New(-4, err.Error())
			mLicenseIns.ErrLog.Println(mErr.Error())
		} else {
			mLicenseIns.LicenseContent = licenseBytes
			mLicenseIns.IsValidLicense = true
		}

		go func() {
			for {
				select {
				case event, ok := <-mLicenseIns.Watcher.Events:
					if !ok {
						return
					}

					if strings.Contains(event.Name, LiceseFileName) == false {
						mLicenseIns.ErrLog.Printf("event, event:%v, name:%s\n", event, event.Name)
						continue
					} else {
						if event.Op&fsnotify.Create == fsnotify.Create { //no modified
							//mLicenseIns.ErrLog.Printf("event modified file, event:%v, name:%s\n", event, event.Name)
						} else if event.Op&fsnotify.Write == fsnotify.Write {
							//mLicenseIns.ErrLog.Printf("event write file, event:%v, name:%s\n", event, event.Name)
						} else if event.Op&fsnotify.Remove == fsnotify.Remove {
							//mLicenseIns.ErrLog.Printf("event remove file, event:%v, name:%s\n", event, event.Name)
						} else {
							errLog.Printf("event contains file, event:%v, name:%s, continue...\n", event, event.Name)
							continue
						}
					}

					licenseBytes, err := verifyLicense(dir, productName)
					if err != nil {
						mErr = public.New(-4, err.Error())
						mLicenseIns.LicenseContent = nil
						mLicenseIns.IsValidLicense = false
					} else {
						mLicenseIns.LicenseContent = licenseBytes
						mLicenseIns.IsValidLicense = true
					}
				case err, ok := <-mLicenseIns.Watcher.Errors:
					if !ok {
						mLicenseIns.ErrLog.Println("close watcher")
						mLicenseIns.IsValidLicense = false
						return
					}
					mLicenseIns.ErrLog.Println(err.Error())
				}
			}
		}()
	})

	if mErr != nil {
		return mErr.Error()
	} else {
		mLicenseIns.Error = mErr.(*public.ErrorMsg)
		gLicenseIns = mLicenseIns
	}
	return ""
}

func FreeLicense() string {
	if gLicenseIns != nil {
		if gLicenseIns.LogFile != nil {
			err := gLicenseIns.LogFile.Close()
			if err != nil {
				return err.Error()
			}
		}

		if gLicenseIns.Watcher != nil {
			err := gLicenseIns.Watcher.Close()
			if err != nil {
				return err.Error()
			}
		}
	}
	return ""
}

//export VerifyLicense
func VerifyLicense(licenseDir string, productName string) int {
	if gLicenseIns != nil {
		if gLicenseIns.Error != nil {
			return gLicenseIns.Error.Code()
		}

		if gLicenseIns.IsValidLicense {
			return 0
		}
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

	// err = startWatcher(licenseDir, productName, true)
	// if err != nil {
	// 	errLog.Println(err.Error())
	// 	return string("FAIL")
	// }

	if public.Exists(licenseDir) == false {
		errLog.Printf("%s dir does not exist\n", licenseDir)
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

	// err = startWatcher(licenseDir, productName, true)
	// if err != nil {
	// 	errLog.Println(err.Error())
	// 	return -1
	// }

	if public.Exists(licenseDir) == false {
		errLog.Printf("%s dir does not exist\n", licenseDir)
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
