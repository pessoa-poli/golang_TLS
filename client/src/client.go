package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	c "server_client_tls/cert_creator/src"
	"time"
)

var (
	//Globals
	PATH_TO_PRIVATE_KEY = "/home/luado/projetos/go_projects/server_client_tls/client/certs/priv.key"
	PATH_TO_PUBLIC_CERT = "/home/luado/projetos/go_projects/server_client_tls/client/certs/cert.key"

	rsaPriv *rsa.PrivateKey

	//Client
	clientCertFile   = PATH_TO_PUBLIC_CERT
	clientKeyFile    = PATH_TO_PRIVATE_KEY
	caCertFile       = PATH_TO_PUBLIC_CERT
	caCertBytes      []byte
	caCertPool       *x509.CertPool
	cert             tls.Certificate
	t                *http.Transport
	globalHTTPClient http.Client
)

//Initializes all variables above.
func init() {
	rsaPriv = c.LoadPEMEncodedFile(PATH_TO_PRIVATE_KEY)
	c.GenerateCertificateGivenPrivateKey(rsaPriv, PATH_TO_PUBLIC_CERT)
	caCertBytes = readCaCert(caCertFile)
	caCertPool = generateCACertPool(caCertBytes)
	cert = x509KeyPairLoader(clientCertFile, clientKeyFile)
	t = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		},
	}
	globalHTTPClient = http.Client{Transport: t, Timeout: 15 * time.Second}
}

func x509KeyPairLoader(clientCertFile, clientKeyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", clientCertFile, clientKeyFile)
	}
	return cert
}

func generateCACertPool(caCert []byte) (caCertPool *x509.CertPool) {
	caCertPool = x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool
}

func readCaCert(caCertFile string) []byte {
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("Error opening cert file %s, Error: %s", caCertFile, err)
	}
	return caCert
}

func main() {
	//Get Request
	resp, err := globalHTTPClient.Get("https://localhost:9500/api/v1")
	if err != nil {
		panic(err.Error())
	}
	rspFull, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("RespoFUll: " + string(rspFull))
	fmt.Println(string(rspFull))

	//Post request
	postBody, _ := json.Marshal(map[string]string{
		"name":     "Tobe",
		"password": "tobebadtobegood",
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err = globalHTTPClient.Post("https://localhost:9500/api/v1", "application/json", responseBody)
	if err != nil {
		panic(err.Error())
	}
	rspFull, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("RespoFUll: " + string(rspFull))
	fmt.Println(string(rspFull))
}
