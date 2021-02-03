package public

import (
	"fmt"
	"strings"
)

var (
	ECDSA_PUBLICKEY = ""

// 	ECDSA_PUBLICKEY = `-----BEGIN PUBLIC KEY-----
// MIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQACVyLsNdjFM6R4IImvTzgRWF0sWjh
// ihmzIyMgyPuqu8IuyzMNx4G2jpoCKhRu9qPCQUMGDeCG1x3/n/OgkWNQANsB82x7
// 7eiIZAl0zcQRH32tcjFILvJ/XCihdoi4MkCnCqlt9/HxjsP590ZtmHfxAeertq5w
// 9vakvpzjPXhkvoMt/Tk=
// -----END PUBLIC KEY-----`
)

func (nea *NonEquAlgorthm) VerifySign(cipherText string) ([]byte, error) {
	var (
		err error
	)

	parts := strings.Split(cipherText, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%s", "license contains an invalid number of segments")
	}

	if err = nea.Algorithm.Verify(parts[0], parts[1], nea.PublicKey); err != nil {
		return nil, err
	}

	plainBytes, err := DecodeSegment(parts[0])
	if err != nil {
		return nil, err
	}

	return plainBytes, nil
}
