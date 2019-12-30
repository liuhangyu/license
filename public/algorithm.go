package public

import "crypto/ecdsa"

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
		alg.PrivateKey, err = ParseECPrivateKeyFromPEM(prikey)
		if err != nil {
			return nil, err
		}
	}

	if pubkey != nil {
		alg.PublicKey, err = ParseECPublicKeyFromPEM(pubkey)
		if err != nil {
			return nil, err
		}
	}

	return alg, nil
}
