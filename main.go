package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type MEMES struct {
	// Image string
	Total    int
	ImageUrl string
}

func ResponseURL(url string) (*http.Response, error) {
	tr := &http.Transport{
		MaxIdleConns:        1,
		MaxIdleConnsPerHost: 1,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  true,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	time.Sleep(time.Second * 5)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}

func MEMEPAGE(response *http.Response) ([]MEMES, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []MEMES{}
	sel := doc.Find("div._3Oa0THmZ3f5iZXAQ0hBJ0k ")
	rank := 1
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("div > img.ImageBox-image")
		link, _ := linkTag.Attr("src")

		link = strings.Trim(link, " ")

		if link != "" && link != "#" {
			result := MEMES{
				rank,
				link,
			}
			results = append(results, result)
			rank += 1
		}
	}
	return results, err
}

func Combination(url string) []MEMES {
	resp, err := ResponseURL(url)
	if err != nil {
		log.Println("error to take response of the page ", err)
		log.Println(err.Error())
		return nil
	}
	time.Sleep(time.Second * 5)
	body, err := MEMEPAGE(resp)
	if err != nil {
		log.Println("error to take page response maybe data is missing ", err)
		log.Println(err.Error())
		return body
	}
	return body
}

func main() {
	url := "https://www.reddit.com/r/dankmemes/"
	body := Combination(url)
	fmt.Printf("%v\n", body)

}
