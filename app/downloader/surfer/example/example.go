package main

import (
	"io"
	"log"
	"time"

	"github.com/andeya/pholcus/app/downloader/surfer"
)

func main() {
	var values = "username=123456@qq.com&password=123456&login_btn=login_btn&submit=login_btn"

	log.Println("********************************************* Surf GET download test start *********************************************")
	resp, err := surfer.Download(&surfer.DefaultRequest{
		Url: "http://www.baidu.com/",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("baidu resp.Status: %s\nresp.Header: %#v\n", resp.Status, resp.Header)

	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("baidu resp.Body: %s\nerr: %v", b, err)

	log.Println("********************************************* Surf POST download test start *********************************************")
	resp, err = surfer.Download(&surfer.DefaultRequest{
		Url:      "http://accounts.lewaos.com/",
		Method:   "POST",
		PostData: values,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("lewaos resp.Status: %s\nresp.Header: %#v\n", resp.Status, resp.Header)

	b, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("lewaos resp.Body: %s\nerr: %v", b, err)

	log.Println("********************************************* PhantomJS GET download test start *********************************************")

	resp, err = surfer.Download(&surfer.DefaultRequest{
		Url:          "http://www.baidu.com/",
		DownloaderID: 1,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("baidu resp.Status: %s\nresp.Header: %#v\n", resp.Status, resp.Header)

	b, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("baidu resp.Body: %s\nerr: %v", b, err)

	log.Println("********************************************* PhantomJS POST download test start *********************************************")

	resp, err = surfer.Download(&surfer.DefaultRequest{
		DownloaderID: 1,
		Url:          "http://accounts.lewaos.com/",
		Method:       "POST",
		PostData:     values,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("lewaos resp.Status: %s\nresp.Header: %#v\n", resp.Status, resp.Header)

	b, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("lewaos resp.Body: %s\nerr: %v", b, err)

	surfer.DestroyJsFiles()

	time.Sleep(10e9)
}
