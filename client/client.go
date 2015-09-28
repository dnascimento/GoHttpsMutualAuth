package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	var endpoint string
	var root string
	var cert string
	var key string
	var client *http.Client

	switch len(os.Args) {
	case 2:
		endpoint = os.Args[1]
		client = HTTPClient()
	case 3:
		endpoint = os.Args[1]
		root = os.Args[2]
		client = HTTPSClient(root)
	case 5:
		endpoint = os.Args[1]
		root = os.Args[2]
		cert = os.Args[3]
		key = os.Args[4]
		client = HTTPSMutual(root, cert, key)
	default:
		fmt.Println("HTTP server with 3 modes: No Authentication, HTTPS server-side only, HTTPS mutual")
		fmt.Println("no authentication usage: http-client  <endpoint>")
		fmt.Println("https server-side usage: http-client  <endpoint> <ca.pem>")
		fmt.Println("https mutual usage     : http-client  <endpoint> <ca.pem> <cert.pem> <key.pem>")
		os.Exit(2)
	}

	get(endpoint, client)
}

//HTTPClient is a simple no auth http client
func HTTPClient() *http.Client {
	return &http.Client{}
}

//HTTPSClient is a server-side auth https client
func HTTPSClient(caCertFile string) *http.Client {
	// Load CA cert
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	return client
}

//HTTPSMutual is a https client that requires mutual authentication
func HTTPSMutual(caCertFile, certFile, keyFile string) *http.Client {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	return client
}

func get(endpoint string, client *http.Client) {
	res, err := client.Get(fmt.Sprintf("%s/hello", endpoint))
	if err != nil {
		log.Fatalf("Unable to form the HTTPS request. %s", err)
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	fmt.Println(res.Status)
	response := string(data)
	if res.Status != "200 OK" || response != "Hello Buddy!\n" {
		log.Fatal(response)
	}
}
