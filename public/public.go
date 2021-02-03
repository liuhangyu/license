package public

import (
	"bytes"
	"fmt"
	"os"

	uuid "license/public/deplib/uuid"
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

func FormatPem(content []byte, isPriKey bool) ([]byte, error) {
	var (
		pemPubKeyStart = []byte("-----BEGIN PUBLIC KEY-----")
		pemPubKeyEnd   = []byte("-----END PUBLIC KEY-----")

		pemPriKeyStart = []byte("-----BEGIN EC PRIVATE KEY-----")
		pemPriKeyEnd   = []byte("-----END EC PRIVATE KEY-----")

		orgContent []byte
		pemContent []byte
		slen       = len(content)
	)

	if slen == 0 {
		if isPriKey && slen < (len(pemPriKeyStart)+len(pemPriKeyEnd)) {
			return nil, fmt.Errorf("%s", "prikey centent length not enough")
		} else if slen < (len(pemPubKeyStart) + len(pemPubKeyEnd)) {
			return nil, fmt.Errorf("%s", "pubkey centent length not enough")
		}
		return nil, fmt.Errorf("%s", "content not null")
	} else {
		orgContent = make([]byte, slen)
		copy(orgContent, content)
		orgContent = bytes.Join(bytes.Fields(orgContent), []byte(" "))
		slen = len(orgContent)
	}

	if isPriKey {
		//私钥
		if bytes.HasPrefix(orgContent, pemPriKeyStart) && bytes.HasSuffix(orgContent, pemPriKeyEnd) {
			pemContent = append(pemContent, pemPriKeyStart...)
			for _, v := range orgContent[len(pemPriKeyStart) : slen-len(pemPriKeyEnd)] {
				if v == ' ' {
					pemContent = append(pemContent, []byte("\n")...)
				} else {
					pemContent = append(pemContent, v)
				}
			}
			pemContent = append(pemContent, pemPriKeyEnd...)
		} else {
			return nil, fmt.Errorf("%s", "pem file prefix must be 'BEGIN/END EC PRIVATE KEY' format")
		}
	} else {
		if bytes.HasPrefix(orgContent, pemPubKeyStart) && bytes.HasSuffix(orgContent, pemPubKeyEnd) {
			pemContent = append(pemContent, pemPubKeyStart...)
			for _, v := range orgContent[len(pemPubKeyStart) : slen-len(pemPubKeyEnd)] {
				if v == ' ' {
					pemContent = append(pemContent, []byte("\n")...)
				} else {
					pemContent = append(pemContent, v)
				}
			}
			pemContent = append(pemContent, pemPubKeyEnd...)
		} else {
			return nil, fmt.Errorf("%s", "pem file prefix must be 'BEGIN/END PUBLIC KEY' format")
		}
	}
	return pemContent, nil
}
