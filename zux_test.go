package zephyr

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func Test_Zuxer(t *testing.T) {
	go func() { /* runMuxer() */ }()
	time.Sleep(time.Second)
	base := "http://localhost:8080"
	printRes(base + "/hello/world/PARAM/wildy")
	printRes(base + "/hello/world/regex/anythingiwant/wow")
	printRes(base + "/hello/WHATSUP")

	select {}
}

func printRes(url string) {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	data, _ := io.ReadAll(res.Body)

	fmt.Println(string(data))
}
