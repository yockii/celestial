package crypto

import (
	"encoding/hex"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/config"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

var gmSm2 = new(sm2Instance)

type sm2Instance struct {
	privateKey *sm2.PrivateKey
}

func initSm2Key() {
	gmSm2.privateKey, _ = x509.ReadPrivateKeyFromHex(config.GetString("sm2.privateKey"))
	if gmSm2.privateKey == nil {
		logger.Fatal("国密SM2密钥无法加载")
	}
}

func Sm2Decrypt(encrypted string) (string, error) {
	bs, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	rs, err := sm2.Decrypt(gmSm2.privateKey, bs, sm2.C1C3C2)
	if err != nil {
		return "", err
	}

	return string(rs), nil
}

func Sm2EncryptHex(origin string, publickKey string) (string, error) {
	pk, _ := x509.ReadPublicKeyFromHex(publickKey)
	rs, err := sm2.Encrypt(pk, []byte(origin), nil, sm2.C1C3C2)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rs), nil
}
