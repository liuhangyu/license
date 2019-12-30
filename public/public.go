package public

import uuid "github.com/satori/go.uuid"

func GetUUID() string {
	return uuid.NewV4().String()
}
