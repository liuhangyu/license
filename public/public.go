package public

import (
	"os"

	uuid "code.uni-ledger.com/switch/license/public/deplib/uuid"
)

func GetUUID() string {
	return uuid.NewV4().String()
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func RemoveDuplicate(list []int) []int {
	var x []int
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}
