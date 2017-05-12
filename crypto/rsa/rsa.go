package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

func ReadPEM(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return ioutil.ReadAll(f)
}

func RSAEncrypt(data []byte, pub string) ([]byte, error) {
	pubkey, err := ReadPEM(pub)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pubkey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, pubInterface.(*rsa.PublicKey), data)
}

func RSADecrypt(data []byte, privkey string) ([]byte, error) {
	pkey, err := ReadPEM(privkey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pkey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, data)
}
