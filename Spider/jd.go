package Spider

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36 Edg/92.0.902.78"
)

func GetProductList(keyword string, productCount int) []string {
	url := "https://search.jd.com/s_new.php"
	client := &http.Client{
		Timeout: time.Second * 3,
	}

	var productIdList []string

	for i := 0; i < productCount/30; i++ {
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Set("referer", "https://search.jd.com/Search")
		request.Header.Set("user-agent", userAgent)

		query := request.URL.Query()
		query.Add("keyword", keyword)
		query.Add("page", strconv.Itoa(i))
		//query.Add("psort", "4")
		request.URL.RawQuery = query.Encode()

		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("商品列表请求出错: %s", err.Error())
			return nil
		}

		body, _ := ioutil.ReadAll(response.Body)
		dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		dom.Find("li.gl-item").Each(func(i int, selection *goquery.Selection) {
			productId, _ := selection.Attr("data-sku")
			productIdList = append(productIdList, productId)
		})
		time.Sleep(time.Millisecond * 300)
	}

	return productIdList
}

func GetCommentImageList(productId string, singleImageCount int) []string {
	url := "https://club.jd.com/discussion/getProductPageImageCommentList.action"
	client := http.Client{
		Timeout: time.Second * 3,
	}

	var imageUrlList []string
	for page := 1; page < singleImageCount/20+1; page++ {
		request, _ := http.NewRequest("GET", url, nil)

		request.Header.Set("User-Agent", userAgent)

		query := request.URL.Query()
		query.Add("productId", productId)
		query.Add("isShadowSku", "0")
		query.Add("page", strconv.Itoa(page))
		query.Add("pageSize", "20")
		request.URL.RawQuery = query.Encode()

		response, err := client.Do(request)
		if err != nil {
			log.Printf("评论图片请求出错: %s", err.Error())
			return nil
		}
		body, _ := ioutil.ReadAll(response.Body)

		var responseJson struct {
			ImageComment struct {
				ImageCount int `json:"imgCount"`
				ImageList  []struct {
					ImageUrl string `json:"imageUrl"`
				} `json:"imgList"`
			} `json:"imgComments"`
		}
		err = json.Unmarshal(body, &responseJson)
		if err != nil {
			log.Printf("json解析出错: %s", err.Error())
			return nil
		}
		responseImageList := responseJson.ImageComment.ImageList
		for i := range responseImageList {
			imageUrlList = append(imageUrlList, "https:"+responseImageList[i].ImageUrl)
		}
	}

	return imageUrlList
}
