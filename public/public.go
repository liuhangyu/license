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
