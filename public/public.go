package public

import (
	uuid "code.uni-ledger.com/switch/license/public/deplib/uuid"
)

func GetUUID() string {
	return uuid.NewV4().String()
}
