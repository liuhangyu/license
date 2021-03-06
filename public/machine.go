package public

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const (
	FSTabCommand  = "blkid"
	FSTab_File    = "/etc/fstab"
	DMI_UUID_FILE = "/sys/class/dmi/id/product_uuid"
)

func CheckCmdExists(command string) (string, error) {
	path, err := exec.LookPath(command)
	if err != nil {
		if runtime.GOOS != "darwin" || //macos,windows测试使用,不获取硬盘分区UUID
			runtime.GOOS != "windows" {
			fmt.Printf("didn't find 'blkid' executable,err %s\n", err.Error())
		}
		return "", err
	}
	return path, nil
}

func IsExistFile(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func ReadSysFile(filePath string) ([]byte, error) {
	isExist := IsExistFile(filePath)
	if isExist == false {
		return nil, fmt.Errorf("%s file does not exist", filePath)
	}

	fd, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("%s", "permission denied get machine id")
		}
		return nil, fmt.Errorf("%s", "get machine id failed")
	}
	defer fd.Close()

	fstabText, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	return fstabText, nil
}

//降低数字串长度
func Sum(data []byte) string {
	var (
		sum    uint64
		length int = len(data)
		index  int
	)

	//以32位求和
	for length >= 4 {
		sum += uint64(data[index])<<24 + uint64(data[index+1])<<16 + uint64(data[index+2])<<8 + uint64(data[index+3])
		index += 4
		length -= 4
	}

	switch length {
	case 3:
		sum += uint64(data[index])<<16 + uint64(data[index+1])<<8 + uint64(data[index+2])
	case 2:
		sum += uint64(data[index])<<8 + uint64(data[index+1])
	case 1:
		sum += uint64(data[index])
	case 0:
		break
	}

	return strconv.FormatUint(sum, 16)
}

/*
获取文件系统类型等信息
*/
func GetFSInfo() ([]string, error) {
	var (
		fsInfoList []string
	)

	cmdPath, err := CheckCmdExists(FSTabCommand)
	if err == nil {
		//存在blkid命令 获取硬盘中文件系统的类型以及UUID
		cmd := exec.Command(cmdPath)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}

		defer stdout.Close()
		if err := cmd.Start(); err != nil {
			return nil, err
		}

		opBytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(opBytes), "\n")

		if len(opBytes) > 0 && len(lines) > 0 {
			for i := 0; i < len(lines); i++ {
				oneLine := strings.TrimSpace(lines[i])
				if oneLine == "" {
					continue
				}

				if strings.Count(oneLine, "xfs") >= 1 {
					ind := strings.Index(oneLine, "UUID=")
					spaceInd := strings.Index(oneLine[ind:], " ")
					if ind != -1 && spaceInd != -1 && ind+6 < len(oneLine) {
						oneLine = oneLine[ind+6 : spaceInd+ind-1]
						fsInfoList = append(fsInfoList, oneLine)
					}
				} else if strings.Count(oneLine, "ext3") >= 1 {
					ind := strings.Index(oneLine, "UUID=")
					spaceInd := strings.Index(oneLine[ind:], " ")
					if ind != -1 && spaceInd != -1 && ind+6 < len(oneLine) {
						oneLine = oneLine[ind+6 : spaceInd+ind-1]
						fsInfoList = append(fsInfoList, oneLine)
					}
				} else if strings.Count(oneLine, "ext4") >= 1 {
					ind := strings.Index(oneLine, "UUID=")
					spaceInd := strings.Index(oneLine[ind:], " ")
					if ind != -1 && spaceInd != -1 && ind+6 < len(oneLine) {
						oneLine = oneLine[ind+6 : spaceInd+ind-1]
						fsInfoList = append(fsInfoList, oneLine)
					}
				}
			}

			// // //test
			// fmt.Println("preCodeList len=", len(fsInfoList))
			// for n := 0; n < len(fsInfoList); n++ {
			// 	fmt.Println(n, fsInfoList[n])
			// }

			if len(fsInfoList) > 0 {
				//排序字符串
				sort.Strings(fsInfoList)
				return fsInfoList, nil
			}
			// ubuntu can not get fsInfo UUID by blkid, can get it by /etc/fstab
			// return nil, fmt.Errorf("%s", "get system file info failed from command")
		}
	}

	fsContent, err := ReadSysFile(FSTab_File)
	if err != nil {
		if runtime.GOOS == "darwin" || //macos,windows测试使用,不获取硬盘分区UUID
			runtime.GOOS == "windows" {
			return nil, nil
		}
		return nil, err
	}

	lines := strings.Split(string(fsContent), "\n")
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" || (len(lines[i]) >= 1 && lines[i][0] == '#') {
			continue
		}

		// //test
		// for i, v := range strings.Fields(lines[i]) {
		// 	fmt.Println(i, v)
		// }
		if len(strings.Fields(lines[i])) == 6 {
			if strings.Fields(lines[i])[2] == "xfs" {
				uuid := strings.Fields(lines[i])[0]
				if len(uuid) > 5 && strings.HasPrefix(uuid, "UUID=") {
					fsInfoList = append(fsInfoList, strings.Fields(lines[i])[0][5:])
				}
			} else if strings.Fields(lines[i])[2] == "ext3" {
				uuid := strings.Fields(lines[i])[0]
				if len(uuid) > 5 && strings.HasPrefix(uuid, "UUID=") {
					fsInfoList = append(fsInfoList, strings.Fields(lines[i])[0][5:])
				}
			} else if strings.Fields(lines[i])[2] == "ext4" {
				uuid := strings.Fields(lines[i])[0]
				if len(uuid) > 5 && strings.HasPrefix(uuid, "UUID=") {
					fsInfoList = append(fsInfoList, strings.Fields(lines[i])[0][5:])
				}
			}
		}
	}
	// // //test
	// for m := 0; m < len(fsUUIDList); m++ {
	// 	fmt.Println(fsUUIDList[m])
	// }

	if len(fsInfoList) > 0 {
		//排序字符串
		sort.Strings(fsInfoList)
		return fsInfoList, nil
	}

	return nil, fmt.Errorf("%s", "get system file info failed from file")
}

//获取BIOS出厂UUID
func GetProductUUID() (string, error) {
	if runtime.GOOS == "darwin" || //macos,windows测试使用,不获取硬盘分区UUID
		runtime.GOOS == "windows" {
		return "", fmt.Errorf("%s", "the operating system not support get product uuid")
	}

	fsContent, err := ReadSysFile(DMI_UUID_FILE)
	if err != nil {
		return "", err
	}

	if len(fsContent) == 0 {
		return "", fmt.Errorf("%s", "get product uuid is empty string")
	}
	return string(fsContent), nil
}

//获取网卡信息
func GetHardwareAddr() ([]string, error) {
	var (
		hardAddrList []string
	)

	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if len(iface.HardwareAddr) >= 6 {
				if strings.HasPrefix(iface.Name, "eno") ||
					strings.HasPrefix(iface.Name, "ens") ||
					strings.HasPrefix(iface.Name, "enp") ||
					strings.HasPrefix(iface.Name, "wl") ||
					strings.HasPrefix(iface.Name, "ww") ||
					strings.HasPrefix(iface.Name, "eth") {
					////test
					// fmt.Println("HardwareAddr", iface.Name, iface.HardwareAddr.String())
					hardAddrList = append(hardAddrList, iface.HardwareAddr.String())
				}
			}
		}
	} else {
		return nil, err
	}

	if len(hardAddrList) > 0 {
		sort.Strings(hardAddrList)
	}
	return hardAddrList, nil
}

/*
GetMachineID 获取机器ID (分区UUID和网卡地址 or Product UUID)
*/
func GetMachineID() (string, error) {
	var (
		MachineIDList []string
	)

	if false {
		//获取磁盘分区UUID
		fsInfoList, err := GetFSInfo()
		if err != nil {
			return "", err
		}
		MachineIDList = append(MachineIDList, fsInfoList...)

		//获取网卡地址
		hardAddrList, err := GetHardwareAddr()
		if err != nil {
			return "", err
		}
		MachineIDList = append(MachineIDList, hardAddrList...)
	} else {
		productUUID, err := GetProductUUID()
		if err != nil {
			return "", err
		}
		MachineIDList = append(MachineIDList, productUUID)
	}

	//编码
	encByte, err := json.Marshal(MachineIDList)
	if err != nil {
		return "", err
	}

	//降低长度
	return Sum(encByte), nil
}
