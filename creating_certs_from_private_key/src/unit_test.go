package main

import (
	"fmt"
	"testing"
)

func TestParseRsaPrivateKeyFromPemStr(t *testing.T) {
	fileBytes := LoadFile(pathToPrivateKeyFile)
	rsaPrivateKey, err := ParseRsaPrivateKeyFromPemStr(fileBytes)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rsaPrivateKey)
}
