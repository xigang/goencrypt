package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

//GetFileContent return file content from filepath
func GetFileContent(certFile string) (data []byte, err error) {
	f, err := os.Open(certFile)
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	data = b
	return
}

//GenRsaKey return generate public and public key file
func GenRsaKey(bits int) error {
	//generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("rsa_private_key.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	//generate pubic key
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("rsa_public_key.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

//GetPrivateKey return private key
func GetPrivateKey(keyData []byte) (priv *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(keyData)
	if block == nil {
		err = errors.New("private key error.")
		return
	}

	priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	return
}

//GetPublicKey return public key
func GetPublicKey(certData []byte) (pub *rsa.PublicKey, err error) {
	block, _ := pem.Decode(certData)
	if block == nil {
		err = errors.New("public key error.")
		return
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pub = pubInterface.(*rsa.PublicKey)
	return
}

//ZeroPadding return padding data
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

//ZeroUnPadding return unpadding data
func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

//PKCS5Padding return  padding data
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS5UnPadding return unpadding data
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
