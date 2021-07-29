#!/bin/sh
#openssl genrsa -out otherCA.key 2048
openssl req -x509 -new -nodes -key otherCA.key -sha256 -days 1825 -out otherCA.crt

openssl genrsa -out otherClient.key 2048
openssl req -new -key otherClient.key -out otherClient.csr
openssl x509 -req -in otherClient.csr -CA otherCA.crt -CAkey otherCA.key -CAcreateserial -out otherClient.crt -days 825 -sha256