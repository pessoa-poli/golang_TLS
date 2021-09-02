package cert_creator

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
)

var (
	pathToPrivateKeyFile = "/home/luado/projetos/go_projects/server_client_tls/creating_certs_from_private_key/certs/priv.key"
	pathToPublicKeyFile  = "/home/luado/projetos/go_projects/server_client_tls/creating_certs_from_private_key/certs/pub.key"
)

//LoadPEMEncodedFile ... Reads the bytes loaded from a PEM encoded file in disk and returns the key stored in it.
func LoadPEMEncodedFile(pathToFile string) *rsa.PrivateKey {
	//Read the file that holds our PEMencodedPrivKey
	fileBytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	//Decode the PEMencoded bytes.
	block, rest := pem.Decode(fileBytes)
	if block == nil {
		panic("failed to decode PEM block containing public key")
	}
	//With the decoded block we can parse it and generate our Private key object
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	//Print remaining bytes (there shouldn't be any)
	fmt.Printf("Got a %T, with remaining data: %q\n", key, rest)

	//Deliver the key to whomever called this function.
	return key
}

func GenerateCertificateGivenPrivateKey(priv *rsa.PrivateKey, pathToPublicKeyFile string) {
	//Our publicKey is embedded on the struct that holds the privateKey.
	pub := &priv.PublicKey

	//Create a certificate template
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{"ORGANIZATION_NAME"},
			Country:       []string{"COUNTRY_CODE"},
			Province:      []string{"PROVINCE"},
			Locality:      []string{"CITY"},
			StreetAddress: []string{"ADDRESS"},
			PostalCode:    []string{"POSTAL_CODE"},
			CommonName:    "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	//Generate our certificate using the pub and priv keys and the template above.
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		log.Println("create ca failed", err)
		return
	}

	//Store the certificate locally in disk.
	certOut, err := os.Create(pathToPublicKeyFile)
	if err != nil {
		panic(err.Error())
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: ca_b})
	certOut.Close()
	log.Print("written cert.pem\n")
}
