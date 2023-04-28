package dingtalk

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/module/uc/model"
	r "math/rand"
	"sort"
	"strings"
	"time"
)

var DingTalkCryptoMap = make(map[uint64]*dingTalkCrypto)

type dingTalkCrypto struct {
	Token          string
	EncodingAESKey string
	SuiteKey       string
	BKey           []byte
	Block          cipher.Block
}

func initDingtalkCrypto(source *model.ThirdSource) (*dingTalkCrypto, error) {
	configuration := source.Configuration
	conf := gjson.Parse(configuration)

	token := conf.Get("callbackToken").String()
	encodingAESKey := conf.Get("callbackAesKey").String()
	suiteKey := conf.Get("appKey").String()

	if len(encodingAESKey) != int(43) {
		return nil, errors.New("不合法的EncodingAESKey")
	}
	bkey, err := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(bkey)
	if err != nil {
		return nil, err
	}
	return &dingTalkCrypto{
		Token:          token,
		EncodingAESKey: encodingAESKey,
		SuiteKey:       suiteKey,
		BKey:           bkey,
		Block:          block,
	}, nil
}

func checkCryptoSource(source *model.ThirdSource) (c *dingTalkCrypto, err error) {
	if source.ID == 0 {
		err = errors.New("sourceId is 0")
		logger.Errorln(err)
		return
	}
	var ok bool
	c, ok = DingTalkCryptoMap[source.ID]
	if !ok {
		c, err = initDingtalkCrypto(source)
		if err != nil {
			return
		}
		DingTalkCryptoMap[source.ID] = c
	}
	return
}

func GetDecryptMsg(source *model.ThirdSource, signature, timestamp, nonce, secretMsg string) (msg string, err error) {
	var c *dingTalkCrypto
	c, err = checkCryptoSource(source)
	if err != nil {
		return
	}
	if c == nil {
		return
	}

	if !c.VerificationSignature(c.Token, timestamp, nonce, secretMsg, signature) {
		return "", errors.New("ERROR: 签名不匹配")
	}
	decode, err := base64.StdEncoding.DecodeString(secretMsg)
	if err != nil {
		return "", err
	}
	if len(decode) < aes.BlockSize {
		return "", errors.New("ERROR: 密文太短")
	}
	blockMode := cipher.NewCBCDecrypter(c.Block, c.BKey[:c.Block.BlockSize()])
	plantText := make([]byte, len(decode))
	blockMode.CryptBlocks(plantText, decode)
	plantText = pkCS7UnPadding(plantText)
	size := binary.BigEndian.Uint32(plantText[16:20])
	plantText = plantText[20:]
	corpID := plantText[size:]
	if string(corpID) != c.SuiteKey {
		logger.Info(corpID, "==============", c.SuiteKey)
		return "", errors.New("ERROR: CorpID匹配不正确")
	}
	return string(plantText[:size]), nil
}

func GetEncryptMsg(source *model.ThirdSource, msg string) (result map[string]string, err error) {
	var c *dingTalkCrypto
	c, err = checkCryptoSource(source)
	if err != nil {
		return
	}
	if c == nil {
		return
	}

	var timestamp = time.Now().Second()
	var nonce = randomString(12)
	str, sign, err := c.GetEncryptMsgDetail(msg, fmt.Sprint(timestamp), nonce)
	if err != nil {
		logger.Errorln(err)
		return
	}

	result = map[string]string{"nonce": nonce, "timeStamp": fmt.Sprint(timestamp), "encrypt": str, "msg_signature": sign}
	return
}

func (c *dingTalkCrypto) GetEncryptMsgDetail(msg, timestamp, nonce string) (string, string, error) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(msg)))
	msg = randomString(16) + string(size) + msg + c.SuiteKey
	plantText := pkCS7Padding([]byte(msg), c.Block.BlockSize())
	if len(plantText)%aes.BlockSize != 0 {
		return "", "", errors.New("ERROR: 消息体size不为16的倍数")
	}
	blockMode := cipher.NewCBCEncrypter(c.Block, c.BKey[:c.Block.BlockSize()])
	chipherText := make([]byte, len(plantText))
	blockMode.CryptBlocks(chipherText, plantText)
	outMsg := base64.StdEncoding.EncodeToString(chipherText)
	signature := c.CreateSignature(c.Token, timestamp, nonce, string(outMsg))
	return string(outMsg), signature, nil
}

func sha1Sign(s string) string {
	// The pattern for generating a hash is `sha1.New()`,
	// `sha1.Write(bytes)`, then `sha1.Sum([]byte{})`.
	// Here we start with a new hash.
	h := sha1.New()

	// `Write` expects bytes. If you have a string `s`,
	// use `[]byte(s)` to coerce it to bytes.
	h.Write([]byte(s))

	// This gets the finalized hash result as a byte
	// slice. The argument to `Sum` can be used to append
	// to an existing byte slice: it usually isn't needed.
	bs := h.Sum(nil)

	// SHA1 values are often printed in hex, for example
	// in git commits. Use the `%x` format verb to convert
	// a hash results to a hex string.
	return fmt.Sprintf("%x", bs)
}

// 数据签名
func (c *dingTalkCrypto) CreateSignature(token, timestamp, nonce, msg string) string {
	params := make([]string, 0)
	params = append(params, token)
	params = append(params, timestamp)
	params = append(params, nonce)
	params = append(params, msg)
	sort.Strings(params)
	return sha1Sign(strings.Join(params, ""))
}

// 验证数据签名
func (c *dingTalkCrypto) VerificationSignature(token, timestamp, nonce, msg, sigture string) bool {
	return c.CreateSignature(token, timestamp, nonce, msg) == sigture
}

// 解密补位
func pkCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

// 加密补位
func pkCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 随机字符串
func randomString(n int, alphabets ...byte) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return string(bytes)
}
