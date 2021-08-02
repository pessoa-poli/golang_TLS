package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	c "server_client_tls/cert_creator/src"
	"time"
)

//Server & client configuration vars
var (
	//Globals
	PATH_TO_PRIVATE_KEY = "/home/luado/projetos/go_projects/server_client_tls/server/certs/priv.key"
	PATH_TO_PUBLIC_CERT = "/home/luado/projetos/go_projects/server_client_tls/server/certs/cert.key"

	// Server
	HOST           = "localhost"
	server         *http.Server
	rsaPriv        = c.LoadPEMEncodedFile(PATH_TO_PRIVATE_KEY)
	serverCertPath = PATH_TO_PUBLIC_CERT
	srvKeyPath     = PATH_TO_PRIVATE_KEY
	certOpt        = tls.RequireAndVerifyClientCert

	//Client
	clientCertFile = PATH_TO_PUBLIC_CERT
	clientKeyFile  = PATH_TO_PRIVATE_KEY
	caCert         []byte
	caCertPool     *x509.CertPool
	cert           tls.Certificate
	t              = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}
	globalHTTPClient = http.Client{Transport: t, Timeout: 15 * time.Second}
)

func init() {
	c.GenerateCertificateGivenPrivateKey(rsaPriv, PATH_TO_PUBLIC_CERT)
	caCert = readCaCert(PATH_TO_PUBLIC_CERT)
	caCertPool = generateCACertPool(caCert)
	cert = x509KeyPairLoader(clientCertFile, clientKeyFile)
	server = &http.Server{
		Addr:         ":" + "9500",
		ReadTimeout:  5 * time.Minute, // 5 min to allow for delays when 'curl' on OSx prompts for username/password
		WriteTimeout: 10 * time.Second,
		TLSConfig:    getTLSConfig(HOST, PATH_TO_PUBLIC_CERT, tls.ClientAuthType(certOpt)),
	}
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
