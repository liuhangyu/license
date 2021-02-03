package public

import "strings"

var (
	ECDSA_PRIVATE = ""

// 	ECDSA_PRIVATE = `-----BEGIN EC PRIVATE KEY-----
// MIHcAgEBBEIB0pE4uFaWRx7t03BsYlYvF1YvKaBGyvoakxnodm9ou0R9wC+sJAjH
// QZZJikOg4SwNqgQ/hyrOuDK2oAVHhgVGcYmgBwYFK4EEACOhgYkDgYYABAAJXIuw
// 12MUzpHggia9POBFYXSxaOGKGbMjIyDI+6q7wi7LMw3HgbaOmgIqFG72o8JBQwYN
// 4IbXHf+f86CRY1AA2wHzbHvt6IhkCXTNxBEffa1yMUgu8n9cKKF2iLgyQKcKqW33
// 8fGOw/n3Rm2Yd/EB56u2rnD29qS+nOM9eGS+gy39OQ==
// -----END EC PRIVATE KEY-----`
)

func (nea *NonEquAlgorthm) SignedBytes(plainText []byte) (string, error) {
	var (
		sig string
		err error
	)

	encBase64String := EncodeSegment(plainText)
	if sig, err = nea.Algorithm.Sign(encBase64String, nea.PrivateKey); err != nil {
		return "", err
	}

	return strings.Join([]string{encBase64String, sig}, "."), nil
}
