/*
*  Symbios Certificate Authority
*  Author: Dario Nascimento
 */
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	var port string
	var cert string
	var key string
	var root string

	switch len(os.Args) {
	case 2:
		port := os.Args[1]
		HTTP(port)
	case 4:
		port = os.Args[1]
		key = os.Args[2]
		cert = os.Args[3]
		HTTPS(port, key, cert)
	case 5:
		port = os.Args[1]
		key = os.Args[2]
		cert = os.Args[3]
		root = os.Args[4]
		HTTPSMutual(port, key, cert, root)
	default:
		fmt.Println("HTTP server with 3 modes: No Authentication, HTTPS server-side only, HTTPS mutual with CA root cert")
		fmt.Println("no authentication usage: http-server  <port>")
		fmt.Println("https server-side usage: http-server  <port> <key.pem> <cert.pem> ")
		fmt.Println("https mutual usage     : http-server  <port> <key.pem> <cert.pem> <cacert.pem>")
		os.Exit(2)
	}
}

//HTTP is a simple no auth http server
func HTTP(port string) {
	http.HandleFunc("/", HandleHello)
	log.Printf("HTTP server listen on port: %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
}

//HTTPS is a server-side auth https server
func HTTPS(port, key, cert string) {
	http.HandleFunc("/", HandleHello)
	log.Printf("HTTPS server listen on port: %s", port)
	err := http.ListenAndServeTLS(":"+port, cert, key, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
}

//HTTPSMutual is a https that requires mutual authentication
func HTTPSMutual(port, key, cert, root string) {
	rootCert, err := ioutil.ReadFile(root)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(rootCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	tlsConfig.BuildNameToCertificate()

	http.HandleFunc("/", HandleHello)
	log.Printf("HTTPS mutual auth server listen on port: %s", port)

	server := &http.Server{
		Addr:      ":" + port,
		TLSConfig: tlsConfig,
	}
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
}

func HandleHello(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello Buddy!\n")
}
