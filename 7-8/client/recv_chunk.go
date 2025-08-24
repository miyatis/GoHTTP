package main

import(
  "crypto/tls"
	"crypto/x509"
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
)

func main(){
	// 証明書を読み込む
	cert, err := os.ReadFile("ca.crt")
	if err != nil {
		panic(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)
	tlsConfig := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()

	// クライアントを作成
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// 通信を行う
	resp, err := client.Get("https://localhost:18443/chunked")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		log.Println(string(bytes.TrimSpace(line)))
	}
}


