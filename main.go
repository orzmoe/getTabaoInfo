package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	html, err := getHtml()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(html)
}
func getHtml() (string, error) {
	_url := "https://detail.m.tmall.com/item.htm?id=636880836537"
	client := &http.Client{}
	request, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	resp, err := client.Do(request) //发送请求
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return "", err
	}
	return string(content), nil
}
