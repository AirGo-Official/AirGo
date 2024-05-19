package encrypt_plugin

import (
	"bytes"
	"crypto/md5"
	crypt_rand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/curve25519"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	PEM_BEGIN_RSA_PUBLIC  = "-----BEGIN RSA PUBLIC KEY-----\n"
	PEM_END_RSA_PUBLIC    = "\n-----END RSA PUBLIC KEY-----"
	PEM_BEGIN_RSA_PRIVATE = "-----BEGIN RSA PRIVATE KEY-----\n"
	PEM_END_RSA_PPRIVATE  = "\n-----END RSA PRIVATE KEY-----"
)

func FormatPublicKey(key string) string {
	if !strings.HasPrefix(key, PEM_BEGIN_RSA_PUBLIC) {
		key = PEM_BEGIN_RSA_PUBLIC + key
	}
	if !strings.HasSuffix(key, PEM_END_RSA_PUBLIC) {
		key = key + PEM_END_RSA_PUBLIC
	}
	return key
}
func FormatPrivateKey(key string) string {
	if !strings.HasPrefix(key, PEM_BEGIN_RSA_PRIVATE) {
		key = PEM_BEGIN_RSA_PRIVATE + key
	}
	if !strings.HasSuffix(key, PEM_END_RSA_PPRIVATE) {
		key = key + PEM_END_RSA_PPRIVATE
	}
	return key
}

func RSAEnCrypt(data, publicKey string) (string, error) {
	key := FormatPublicKey(publicKey)
	block, _ := pem.Decode([]byte(key))
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	encryptedData, err := rsa.EncryptPKCS1v15(crypt_rand.Reader, pubKey.(*rsa.PublicKey), []byte(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedData), err
}

func RSADecrypt(encryptedData, privateKey string) (string, error) {
	key := FormatPrivateKey(privateKey)

	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode([]byte(key))
	priKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	originalData, err := rsa.DecryptPKCS1v15(crypt_rand.Reader, priKey.(*rsa.PrivateKey), encryptedDecodeBytes)
	return string(originalData), err
}

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptEncode(password string) string {
	// Go 中的 bcrypt.DefaultCost 是 10
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptDecode(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// md5 encode
func Md5Encode(str string, isUpper bool) string {
	sum := md5.Sum([]byte(str))
	res := hex.EncodeToString(sum[:])
	//转大写，strings.ToUpper(res)
	if isUpper {
		res = strings.ToUpper(res)
	}
	return res
}

// sha256 encode
func Sha256Encode(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	res := hex.EncodeToString(hash.Sum(nil))
	return res
}

// 随机数，n为 位数
func RandomString(n int) string {
	var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomStr := make([]rune, n)
	for i := range randomStr {
		randomStr[i] = defaultLetters[r.Intn(len(defaultLetters))]
	}
	return string(randomStr)
}

// 随机数字
func RandomNumber(start, end int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(end-start) + start
}

func StrToUnicode(str string) string {
	DD := []rune(str) //需要分割的字符串内容，将它转为字符，然后取长度。
	finallStr := ""
	for i := 0; i < len(DD); i++ {
		if unicode.Is(unicode.Scripts["Han"], DD[i]) {
			textQuoted := strconv.QuoteToASCII(string(DD[i]))
			finallStr += textQuoted[1 : len(textQuoted)-1]
		} else {
			h := fmt.Sprintf("%x", DD[i])
			finallStr += "\\u" + isFullFour(h)
		}
	}
	return finallStr
}

func isFullFour(str string) string {
	if len(str) == 1 {
		str = "000" + str
	} else if len(str) == 2 {
		str = "00" + str
	} else if len(str) == 3 {
		str = "0" + str
	}
	return str
}

// unicode 转字符
func UnicodeToStr(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// 订阅base64 解码
func SubBase64Decode(str string) string {
	i := len(str) % 4
	switch i {
	case 1:
		str = str[:len(str)-1]
	case 2:
		str += "=="
	case 3:
		str += "="
	}
	//str = strings.Split(str, "//")[1]
	var data []byte
	var err error
	if strings.Contains(str, "-") || strings.Contains(str, "_") {
		data, err = base64.URLEncoding.DecodeString(str)

	} else {
		data, err = base64.StdEncoding.DecodeString(str)
		//data, err = base64.RawURLEncoding.DecodeString(str)
	}
	if err != nil {
		fmt.Println(err)
	}
	return string(data)
}

// Private key: sJxwD9sEodPf97oNG872idTkFhxlkFXLsTmRxVWvx2g
// Public key: HC5OmUS72bUGoVH1ONQjDUIIskr5LqHuUG4nGsg7dhc
func ExecuteX25519(str string) (string, string, error) {
	var err error
	var privateKey []byte
	var publicKey []byte
	var privateKeyStr string
	var publicKeyStr string

	privateKey, err = base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		goto out
	}
	if len(privateKey) != curve25519.ScalarSize {
		err = errors.New("Invalid length of private key.")
		goto out
	}
	// Modify random bytes using algorithm described at:
	// https://cr.yp.to/ecdh.html.
	privateKey[0] &= 248
	privateKey[31] &= 127
	privateKey[31] |= 64

	if publicKey, err = curve25519.X25519(privateKey, curve25519.Basepoint); err != nil {
		goto out
	}
	//output = fmt.Sprintf("Private key: %v\nPublic key: %v",
	//	base64.RawURLEncoding.EncodeToString(privateKey),
	//	base64.RawURLEncoding.EncodeToString(publicKey))
	publicKeyStr = base64.RawURLEncoding.EncodeToString(publicKey)
	privateKeyStr = base64.RawURLEncoding.EncodeToString(privateKey)
out:
	return publicKeyStr, privateKeyStr, err
}

func JsonMarshal(data any) (string, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	err := jsonEncoder.Encode(data)
	if err != nil {
		return "", err
	}
	return bf.String(), nil
}
