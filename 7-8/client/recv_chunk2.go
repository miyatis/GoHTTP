package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"log"
)

func main() {
	// CA証明書を読み込む
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		panic(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	// TLS設定
	tlsConfig := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	// TLSソケットオープン
	dialer := &net.Dialer{
		Timeout: 30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", "localhost:18443", tlsConfig)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// リクエスト送信
	request, err := http.NewRequest("GET", "https://localhost:18443/chunked", nil)
	err = request.Write(conn)
	if err != nil {
		panic(err)
	}

	// 読み込み
	reader := bufio.NewReader(conn)
	// フィールドを読む
	resp, err := http.ReadResponse(reader, request)
	if err != nil {
		panic(err)
	}
	if resp.TransferEncoding[0] != "chunked" {
		panic("wrong transfer encoding")
	}
	for {
		// サイズ取得
		sizeStr, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		log.Printf("sizeStr: %s", sizeStr)

		// 16進数のサイズをパース。サイズが0なら終了
		size, err := strconv.ParseInt(string(sizeStr[:len(sizeStr)-2]), 16, 64)
		if size == 0 {
			break
		}
		if err != nil {
			panic(err)
		}
		log.Println(size)
		// サイズ数分バッファを確保して読み込み
		line := make([]byte, int(size))
		reader.Read(line)
		reader.Discard(2)
		log.Printf(" %s\n", strings.TrimSpace(string(line)))
	}
}
