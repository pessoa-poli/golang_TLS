package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"
)

var pathToPrivateKeyFile = "/home/luado/projetos/go_projects/server_client_tls/creating_certs_from_private_key/certs/root.key"

func LoadFile(pathToFile string) []byte {
	fileBytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	return fileBytes
}

func ParseRsaPrivateKeyFromPemStr(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func main() {
	//pubPrivKeyPair, _ := rsa.GenerateKey(rand.Reader,2048)
	privateKeyString := LoadFile(pathToPrivateKeyFile)
	rsaPrivateKey, err := ParseRsaPrivateKeyFromPemStr(privateKeyString)
	if err != nil {
		panic(err.Error())
	}
	caTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Country:            []string{"Brazil_Camtec"},
			Organization:       []string{"Camtec"},
			OrganizationalUnit: []string{"Main"},
			Locality:           []string{"Rio de Janeiro"},
			Province:           []string{"Nova América"},
			StreetAddress:      []string{"num"},
			PostalCode:         []string{"123"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 1, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	ca, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, rsaPrivateKey.PublicKey, rsaPrivateKey)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("CA is: %v\n", ca)
}
