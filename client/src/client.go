package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	//Client
	host           = "otherClient"
	workingDir, _  = os.Getwd()
	clientCertFile = workingDir + fmt.Sprintf("/../certs/otherca/%s.crt", host)
	clientKeyFile  = workingDir + fmt.Sprintf("/../certs/otherca/%s.key", host)
	caCertFile     = workingDir + "/../certs/otherca/otherCA.crt"
	caCertBytes    = readCaCert(caCertFile)
	caCertPool     = generateCACertPool(caCertBytes)
	cert           = x509KeyPairLoader(clientCertFile, clientKeyFile)
	t              = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}
	globalHTTPClient = http.Client{Transport: t, Timeout: 15 * time.Second}
)

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
}
