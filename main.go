package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
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
	// time.Sleep(time.Second * 60)
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

func DownloadMEME(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {

		// panic(err)
		return body, err
	}
	return body, err
}

func CreateFolder() {
	_, err := os.Stat("memes")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("memes", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}

func GenerateFile(filename string) (*os.File, error) {
	create, err := os.Create("./memes/" + filename + ".jpg")
	if err != nil {
		return nil, err
	}
	return create, nil
}

func AppendData(file *os.File, data []byte) (string, error) {
	size, err := file.WriteString(string(data))
	if err != nil {
		return "", err
	}
	return string(size), nil
}

func Combination(url string) []MEMES {
	resp, err := ResponseURL(url)
	if err != nil {
		log.Println("error to take response of the page ", err)
		log.Println(err.Error())
		return nil
	}
	// time.Sleep(time.Second * 5)
	body, err := MEMEPAGE(resp)
	if err != nil {
		log.Println("error to take page response maybe data is missing ", err)
		log.Println(err.Error())
		return body
	}
	return body
}

func SaveIntoSystem(url string) {
	body := Combination(url)
	// fmt.Printf("%v\n", body)
	CreateFolder()
	for _, resp := range body {
		data, err := ResponseURL(resp.ImageUrl)
		if err != nil {
			log.Println("error to take response of the page ", err)
			log.Println(err.Error())
		}
		body, err := DownloadMEME(data)
		if err != nil {
			log.Println("Error to get Image response ", body)
			log.Println(err)
			log.Println(err.Error())
		}
		filename := strings.Replace(resp.ImageUrl, "/", "", -1)
		file, err := GenerateFile(filename)
		if err != nil {
			log.Printf("Error to create File %v\n", filename)
			log.Println(err)
			log.Println(err.Error())
		}
		_, err = AppendData(file, body)
		if err != nil {
			log.Println("Error to Append data into File ", filename)
			log.Println(err)
			log.Println(err.Error())
		}
		fmt.Printf("%v\n", filename)

	}
	fmt.Println("Memes Downloaded Successfully")

}

func main() {
	url := "https://www.reddit.com/r/dankmemes/"
	SaveIntoSystem(url)
}
