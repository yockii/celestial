package crypto

import (
	"crypto/rand"
	"encoding/hex"
	logger "github.com/sirupsen/logrus"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

var gmSm2 = new(sm2Instance)

type sm2Instance struct {
	privateKey *sm2.PrivateKey
	publicKey  string
}

func initSm2Key() {
	priv, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		logger.Fatal("国密SM2密钥无法生成")
	}
	pub := &priv.PublicKey

	gmSm2.privateKey = priv

	gmSm2.publicKey = x509.WritePublicKeyToHex(pub)

	//gmSm2.privateKey, _ = x509.ReadPrivateKeyFromHex(config.GetString("sm2.privateKey"))
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

func Sm2Encrypt(origin string) (string, error) {
	rs, err := sm2.Encrypt(&gmSm2.privateKey.PublicKey, []byte(origin), nil, sm2.C1C3C2)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rs), nil
}

func PublicKeyString() string {
	return gmSm2.publicKey
}
