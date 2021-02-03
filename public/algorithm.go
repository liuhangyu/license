package public

import (
	"crypto/ecdsa"
	"encoding/pem"
)

type NonEquAlgorthm struct {
	Algorithm  *SigningMethodECDSA
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func GetNonEquAlgorthm(prikey []byte, pubkey []byte) (*NonEquAlgorthm, error) {
	var err error
	alg := &NonEquAlgorthm{
		Algorithm: GetSignVerifyMgr(),
	}

	if prikey != nil {
		if block, _ := pem.Decode(prikey); block == nil {
			prikey, err = FormatPem(prikey, true)
			if err != nil {
				return nil, err
			}
		}

		alg.PrivateKey, err = ParseECPrivateKeyFromPEM(prikey)
		if err != nil {
			return nil, err
		}
	}

	if pubkey != nil {
		if block, _ := pem.Decode(pubkey); block == nil {
			pubkey, err = FormatPem(pubkey, false)
			if err != nil {
				return nil, err
			}
		}

		alg.PublicKey, err = ParseECPublicKeyFromPEM(pubkey)
		if err != nil {
			return nil, err
		}
	}

	return alg, nil
}
