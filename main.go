package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	info, err := GetBabyInfo(os.Getenv("BABY_ID"), "0")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(info)
}

func getHtml(_url string, isMobile bool) (string, error) {

	client := &http.Client{}
	request, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return "", err
	}
	//Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1
	if isMobile {
		request.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")

	} else {
		request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36")

	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
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
func GetBabyInfo(id, babyType string) (BabyInfoResp, error) {
	babyInfoResp := BabyInfoResp{}
	if babyType == "1" || babyType == "0" {
		_url := "https://detail.m.tmall.hk/item.htm?id=" + id + "&toSite=main&sku_properties=29112:97926"
		html, err := getHtml(_url, true)
		if err != nil {
			return babyInfoResp, err
		}
		rgp := regexp.MustCompile(`og:image" content="(.+?)"`)
		images := rgp.FindStringSubmatch(html)
		if len(images) > 0 {
			rgp = regexp.MustCompile(`class="shop-name" title=".+?">(.+?)<`)
			shops := rgp.FindStringSubmatch(html)
			shopName, err := GbkToUtf8(shops[1])
			if err != nil {
				babyInfoResp.ShopName = shops[1]
			} else {
				babyInfoResp.ShopName = shopName
			}

			babyInfoResp.Image = images[1]
			babyInfoResp.Type = 1
			return babyInfoResp, nil
		}
	}
	if babyType == "2" || babyType == "0" {
		_url := "https://item.taobao.com/item.htm?ft=t&id=" + id
		html, err := getHtml(_url, false)
		if err != nil {
			return babyInfoResp, err
		}
		rgp := regexp.MustCompile(`pic.+?: '(.+?)',`)
		images := rgp.FindStringSubmatch(html)
		if len(images) > 0 {
			rgp = regexp.MustCompile(`shopName.+?: '(.+?)',`)
			shops := rgp.FindStringSubmatch(html)
			v, _ := zhToUnicode([]byte(shops[1]))

			babyInfoResp.ShopName = string(v)
			babyInfoResp.Image = images[1]
			babyInfoResp.Type = 2
		}
	}
	return babyInfoResp, nil
}

func zhToUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

type BabyInfoResp struct {
	ShopName string `json:"shop_name"`
	Image    string `json:"image"`
	Type     int    `json:"type"`
}

func GbkToUtf8(s string) (string, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return "", e
	}
	return string(d), nil
}
