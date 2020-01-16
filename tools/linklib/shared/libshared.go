package main

import "C"
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

//export NewLicense
/*
创建license对象
dir 存放license.dat的文件夹(监控文件夹使用)
productName 在二进制中每个产品的英文名称
logPath 指定路径的license.log路径

成功返回空字符
失败返回errmsg
*/
func NewLicense(dir string, productName string, logPath string) *C.char {
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
			return nil, public.ErrLoadPubKey.SetErr(err)
		}

		//读取license.dat文件
		licenseFilePath := path.Join(dir, LiceseFileName)
		licenseBytes, err := public.ReadLicensePem(licenseFilePath)
		if err != nil {
			return nil, public.ErrReadAuthFile.SetErr(err)
		}

		{
			//验证pem格式
			if pemBlock, _ = pem.Decode(licenseBytes); pemBlock == nil {
				return nil, public.ErrDecodeAuthFile
			}

			//验证签名
			content, err = alg.VerifySign(string(pemBlock.Bytes))
			if err != nil {
				return nil, public.ErrVerifySign.SetErr(err)
			}

			//获取license对象
			license, err := public.ToLicense(content)
			if err != nil {
				return nil, public.ErrUnmarshalLiObj.SetErr(err)
			}

			//获取机器ID
			machineID, err := public.GetMachineID()
			if err != nil {
				return nil, public.ErrGetMachineCode.SetErr(err)
			}

			//对比过期时间以及机器ID
			err = license.Valid(proName, machineID)
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
			return C.CString(mErr.Error())
		}
	}

	mLicenseIns.LogFile, mErr = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if mErr != nil {
		return C.CString(mErr.Error())
	}

	//初始化日志
	mLicenseIns.ErrLog = log.New(mLicenseIns.LogFile, "[License]", log.LstdFlags|log.Lshortfile|log.LstdFlags|log.Lmicroseconds)
	mLicenseIns.ErrLog.Println("License", public.GetAppInfo())
	mLicenseIns.Error = nil

	mLicenseIns.Once.Do(func() {
		var (
			err error
		)

		for {
			if public.Exists(dir) == false {
				mLicenseIns.Error = public.ErrDirNoExist.SetErrText(dir)
				break
			}

			//启动目录监听
			mLicenseIns.Watcher, err = fsnotify.NewWatcher()
			if err != nil {
				mLicenseIns.Error = public.ErrNewWatcher.SetErr(err)
				break
			}

			//监控目录
			err = mLicenseIns.Watcher.Add(dir)
			if err != nil {
				mLicenseIns.Error = public.ErrWatcherAdd.SetErr(err)
				break
			}
			break
		}
		if mLicenseIns.Error != nil {
			mLicenseIns.ErrLog.Println(mLicenseIns.Error.Error())
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
			if iErr, ok := err.(*public.ErrorMsg); ok {
				mLicenseIns.Error = iErr
			} else {
				mLicenseIns.Error = public.ErrUnKnown.SetErr(err)
			}

			mLicenseIns.ErrLog.Println(err.Error())
		}
		mLicenseIns.LicenseContent = licenseBytes
		mLicenseIns.IsValidLicense = true
		gLicenseIns = mLicenseIns

		go func() {
			for {
				select {
				case event, ok := <-gLicenseIns.Watcher.Events:
					if !ok {
						return
					}

					if strings.Contains(event.Name, LiceseFileName) == false {
						gLicenseIns.ErrLog.Printf("event, event:%v, name:%s\n", event, event.Name)
						continue
					} else {
						if event.Op&fsnotify.Create == fsnotify.Create { //no modified
							//gLicenseIns.ErrLog.Printf("event modified file, event:%v, name:%s\n", event, event.Name)
						} else if event.Op&fsnotify.Write == fsnotify.Write {
							//gLicenseIns.ErrLog.Printf("event write file, event:%v, name:%s\n", event, event.Name)
						} else if event.Op&fsnotify.Remove == fsnotify.Remove {
							//gLicenseIns.ErrLog.Printf("event remove file, event:%v, name:%s\n", event, event.Name)
						} else {
							gLicenseIns.ErrLog.Printf("event contains file, event:%v, name:%s, continue...\n", event, event.Name)
							continue
						}
					}

					licenseBytes, err := verifyLicense(dir, productName)
					if err != nil {
						gLicenseIns.ErrLog.Println(err.Error())

						if iErr, ok := err.(*public.ErrorMsg); ok {
							gLicenseIns.Error = iErr
						} else {
							gLicenseIns.Error = public.ErrUnKnown.SetErr(err)
						}

						gLicenseIns.LicenseContent = nil
						gLicenseIns.IsValidLicense = false
					} else {
						gLicenseIns.Error = nil
						gLicenseIns.LicenseContent = licenseBytes
						gLicenseIns.IsValidLicense = true
					}
				case err, ok := <-gLicenseIns.Watcher.Errors:
					if !ok {
						gLicenseIns.ErrLog.Println("close watcher")
						gLicenseIns.IsValidLicense = false
						return
					}
					gLicenseIns.ErrLog.Println(err.Error())
				}
			}
		}()
	})

	if mLicenseIns != nil && mLicenseIns.Error != nil {
		gLicenseIns.Error = mLicenseIns.Error
		return C.CString(mLicenseIns.Error.Error())
	}

	if gLicenseIns != nil && gLicenseIns.Error != nil {
		return C.CString(gLicenseIns.Error.Error())
	}

	gLicenseIns = mLicenseIns
	return C.CString("")
}

//export FreeLicense
/*
释放license对象
*/
func FreeLicense() *C.char  {
	if gLicenseIns != nil {
		if gLicenseIns.LogFile != nil {
			err := gLicenseIns.LogFile.Close()
			if err != nil {
				return C.CString(err.Error())
			}
		}

		if gLicenseIns.Watcher != nil {
			err := gLicenseIns.Watcher.Close()
			if err != nil {
				return C.CString(err.Error())
			}
		}
	}
	return C.CString("")
}

//export VerifyLicense
/*
licenseDir 存放license.dat的文件夹(监控文件夹使用)
productName 在二进制中每个产品的英文名称

失败错误信息查看:errors.go中定义
*/
func VerifyLicense(licenseDir string, productName string) C.int {
	var (
		err = public.ErrNoCreateObj
	)

	if gLicenseIns != nil {
		if gLicenseIns.Error != nil {
			gLicenseIns.ErrLog.Println(gLicenseIns.Error.Error())
			return C.int(gLicenseIns.Error.GetCode())
		}

		if gLicenseIns.IsValidLicense {
			return C.int(1)
		}
	}

	return C.int(err.GetCode())
}

//export ReadLicnese
/*
licenseDir 存放license.dat的文件夹(监控文件夹使用)
productName 在二进制中每个产品的英文名称

返回:
失败返回"FAIL"
成功返回定义的KV配置项
*/
func ReadLicnese(licenseDir string, productName string)  *C.char {
	var (
		pemBlock     *pem.Block
		lbytes       []byte
		err          error
		licenseBytes []byte
	)

	if gLicenseIns == nil {
		return C.CString("FAIL")
	}

	if public.Exists(licenseDir) == false {
		gLicenseIns.ErrLog.Printf("%s dir does not exist\n", licenseDir)
		return C.CString("FAIL")
	}

	for {
		if gLicenseIns != nil && len(gLicenseIns.LicenseContent) == 0 {
			//读取license.dat文件
			licenseFilePath := path.Join(licenseDir, LiceseFileName)
			licenseBytes, err = public.ReadLicensePem(licenseFilePath)
			if err != nil {
				break
			}
			gLicenseIns.LicenseContent = licenseBytes
		} else {
			licenseBytes = gLicenseIns.LicenseContent
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

		lbytes, err = json.Marshal(license.CustomKV)
		if err != nil {
			break
		}
		break
	}

	if err != nil {
		gLicenseIns.ErrLog.Println(err.Error())
		return C.CString("FAIL")
	}

	return  C.CString(string(lbytes))
}

//export GetExpireSec
/*
licenseDir 存放license.dat的文件夹(监控文件夹使用)
productName 在二进制中每个产品的英文名称

返回值:
0已过期
-1失败
>0 剩余时间(剩余未过期的秒数)
*/
func GetExpireSec(licenseDir string, productName string) C.longlong {
	var (
		pemBlock     *pem.Block
		err          error
		licenseBytes []byte
		seconds      int64
	)

	if gLicenseIns == nil {
		return C.longlong(-1)
	}

	if public.Exists(licenseDir) == false {
		gLicenseIns.ErrLog.Printf("%s dir does not exist\n", licenseDir)
		return C.longlong(-1)
	}

	for {
		if gLicenseIns != nil && len(gLicenseIns.LicenseContent) == 0 {
			//读取license.dat文件
			licenseFilePath := path.Join(licenseDir, LiceseFileName)
			licenseBytes, err = public.ReadLicensePem(licenseFilePath)
			if err != nil {
				gLicenseIns.ErrLog.Println(err.Error())
				break
			}
			gLicenseIns.LicenseContent = licenseBytes
		} else {
			licenseBytes = gLicenseIns.LicenseContent
		}

		//验证pem格式
		if pemBlock, _ = pem.Decode(licenseBytes); pemBlock == nil {
			err = fmt.Errorf("%s", "license must be PEM")
			gLicenseIns.ErrLog.Println(err.Error())
			break
		}

		license, err := public.BytesToLicense(string(pemBlock.Bytes))
		if err != nil {
			gLicenseIns.ErrLog.Println(err.Error())
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
