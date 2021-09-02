package cert_creator

import (
	"testing"
)

func TestLoadPEMEncodedFile(t *testing.T) {
	LoadPEMEncodedFile(pathToPrivateKeyFile)
}

func TestGenerateCertificateGivenPrivateKey(t *testing.T) {
	priv := LoadPEMEncodedFile(pathToPrivateKeyFile)
	GenerateCertificateGivenPrivateKey(priv, pathToPublicKeyFile)
}
