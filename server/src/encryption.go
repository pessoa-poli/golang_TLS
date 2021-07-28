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

//Server & client configuration vars
var (
	// Server
	server_file   = "otherClient"
	workingDir, _ = os.Getwd()
	serverCert    = workingDir + fmt.Sprintf("/../certs/otherca/%s.crt", server_file)
	srvKey        = workingDir + fmt.Sprintf("/../certs/otherca/%s.key", server_file)
	caCertFile    = workingDir + "/../certs/otherca/otherCA.crt"
	certOpt       = tls.RequireAndVerifyClientCert
	server        = &http.Server{
		Addr:         ":" + "9500",
		ReadTimeout:  5 * time.Minute, // 5 min to allow for delays when 'curl' on OSx prompts for username/password
		WriteTimeout: 10 * time.Second,
		TLSConfig:    getTLSConfig(host, caCertFile, tls.ClientAuthType(certOpt)),
	}
	//Client
	host           = "RNP_CA"
	clientCertFile = workingDir + fmt.Sprintf("/../certs/%s.crt", host)
	clientKeyFile  = workingDir + fmt.Sprintf("/../certs/%s.key", host)
	caCert         = readCaCert(caCertFile)
	caCertPool     = generateCACertPool(caCert)
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

func getTLSConfig(host, caCertFile string, certOpt tls.ClientAuthType) *tls.Config {
	var caCert []byte
	var err error
	var caCertPool *x509.CertPool
	if certOpt > tls.RequestClientCert {
		caCert, err = ioutil.ReadFile(caCertFile)
		if err != nil {
			log.Fatal("Error opening cert file: ", caCertFile, ", error ", err)
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	return &tls.Config{
		ServerName: host,
		ClientAuth: certOpt,
		ClientCAs:  caCertPool,
		MinVersion: tls.VersionTLS12, // TLS versions below 1.2 are considered insecure - see https://www.rfc-editor.org/rfc/rfc7525.txt for details
	}
}
