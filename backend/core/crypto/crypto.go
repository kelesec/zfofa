package crypto

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

/*
----------------------------------------------------------------------------------------------------------
Rsa: RSA加密
----------------------------------------------------------------------------------------------------------
*/

type Rsa struct {
}

// LoadPrivateKey 加载私钥
func (r *Rsa) LoadPrivateKey(key string) (*rsa.PrivateKey, error) {
	privateKeyBytes := bytes.NewBufferString(key).Bytes()
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("error: pem decode fail")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// SHA256WithRSA 进行SHA256签名，返回Base64加密后的字符串
func (r *Rsa) SHA256WithRSA(key *rsa.PrivateKey, plainText string) (string, error) {
	hashed := sha256.Sum256(bytes.NewBufferString(plainText).Bytes())
	signBytes, err := rsa.SignPKCS1v15(nil, key, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signBytes), nil
}

// GenSignWithRsaSha256 SHA256签名生成
func GenSignWithRsaSha256(s, privateKey string) string {
	r := Rsa{}
	priKey, err := r.LoadPrivateKey(privateKey)
	if err != nil {
		return ""
	}

	sign, err := r.SHA256WithRSA(priKey, s)
	if err != nil {
		return ""
	}

	return sign
}
