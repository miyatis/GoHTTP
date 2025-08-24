package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

// 単純なHTMLレスポンスを返すハンドラ
func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, "Error dumping request", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>")
}

// レスポンスをチャンクで返すハンドラ
func handleChunkResponce(w http.ResponseWriter, r *http.Request) {
	c := http.NewResponseController(w)

	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "Chunk %d\n", i)
		c.Flush()
		time.Sleep(500 * time.Millisecond)
	}
	c.Flush()
}

// 実行方法: 7-8/serverディレクトリで$go run server.go
func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/chunked", handleChunkResponce)

	log.Println("start http listening :18443")
	err := http.ListenAndServeTLS(":18443", "server.crt", "server.key", nil)
	log.Println(err)
}
