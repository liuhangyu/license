package public

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const (
	FSTabCommand = "blkid"
	FSTabFile    = "/etc/fstab"
)

func CheckCmdExists(command string) (string, error) {
	path, err := exec.LookPath(command)
	if err != nil {
		fmt.Printf("didn't find 'blkid' executable\n")
		return "", err
	}
	return path, nil
}

func IsExistFStab(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func ReadFStabFile(filePath string) ([]byte, error) {
	isExist := IsExistFStab(filePath)
	if isExist == false {
		return nil, fmt.Errorf("%s file does not exist", filePath)
	}

	fd, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	fstabText, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	return fstabText, nil
}

/*
GetUniqueMachineID 获取机器ID
获取硬盘分区UUID
*/
func GetUniqueMachineID() (string, error) {
	cmdPath, err := CheckCmdExists(FSTabCommand)
	if err == nil {
		//存在blkid命令
		//获取硬盘分区UUID
		cmd := exec.Command(cmdPath)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return "", err
		}

		defer stdout.Close()
		if err := cmd.Start(); err != nil {
			return "", err
		}

		opBytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			return "", err
		}

		lines := strings.Split(string(opBytes), "\n")
		if len(lines) > 0 {
			var preCodeList []string

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
						preCodeList = append(preCodeList, oneLine)
					}
				} else if strings.Count(oneLine, "ext3") >= 1 {
					ind := strings.Index(oneLine, "UUID=")
					spaceInd := strings.Index(oneLine[ind:], " ")
					if ind != -1 && spaceInd != -1 && ind+6 < len(oneLine) {
						oneLine = oneLine[ind+6 : spaceInd+ind-1]
						preCodeList = append(preCodeList, oneLine)
					}
				} else if strings.Count(oneLine, "ext4") >= 1 {
					ind := strings.Index(oneLine, "UUID=")
					spaceInd := strings.Index(oneLine[ind:], " ")
					if ind != -1 && spaceInd != -1 && ind+6 < len(oneLine) {
						oneLine = oneLine[ind+6 : spaceInd+ind-1]
						preCodeList = append(preCodeList, oneLine)
					}
				}
			}

			// // // test
			// fmt.Println(len(preCodeList))
			// for n := 0; n < len(preCodeList); n++ {
			// 	fmt.Println(n, preCodeList[n])
			// }

			if len(preCodeList) > 0 {
				//排序字符串
				sort.Strings(preCodeList)
				//编码
				encByte, err := json.Marshal(preCodeList)
				if err != nil {
					return "", err
				}
				return base64.StdEncoding.EncodeToString(encByte), nil
			}

			return "", fmt.Errorf("%s", "failed to get machine id")
		}
	}

	fsContent, err := ReadFStabFile(FSTabFile)
	if err != nil {
		return "", err
	}

	var fsUUIDList []string
	lines := strings.Split(string(fsContent), "\n")
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" || (len(lines[i]) >= 1 && lines[i][0] == '#') {
			continue
		}

		// // test
		// for i, v := range strings.Fields(lines[i]) {
		// 	fmt.Println(i, v)
		// }
		if len(strings.Fields(lines[i])) == 6 {
			if strings.Fields(lines[i])[2] == "xfs" {
				if len(strings.Fields(lines[i])[0]) > 5 {
					fsUUIDList = append(fsUUIDList, strings.Fields(lines[i])[0][5:])
				}
			} else if strings.Fields(lines[i])[2] == "ext3" {
				if len(strings.Fields(lines[i])[0]) > 5 {
					fsUUIDList = append(fsUUIDList, strings.Fields(lines[i])[0][5:])
				}
			} else if strings.Fields(lines[i])[2] == "ext4" {
				if len(strings.Fields(lines[i])[0]) > 5 {
					fsUUIDList = append(fsUUIDList, strings.Fields(lines[i])[0][5:])
				}
			}
		}
	}
	// // //test
	// for m := 0; m < len(fsUUIDList); m++ {
	// 	fmt.Println(fsUUIDList[m])
	// }

	if len(fsUUIDList) > 0 {
		//排序字符串
		sort.Strings(fsUUIDList)
		//编码
		encByte, err := json.Marshal(fsUUIDList)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(encByte), nil
	}

	return "", fmt.Errorf("%s", "get machine id abort")
}
