package main

import (
	"GOProject/Spider"
	"fmt"
	"log"
	"time"
)

func main() {
	var productCount int
	var singleImageCount int

	productNameList := Spider.LoadProduct("./tag.txt")

	fmt.Print("请输入爬取的商品数量(30的倍数): ")
	_, err := fmt.Scanln(&productCount)
	if err != nil {
		log.Printf("商品数量输入错误: %s\n", err.Error())
		return
	}

	fmt.Print("请输入单个商品图片数量(20的倍数): ")
	_, err = fmt.Scanln(&singleImageCount)
	if err != nil {
		log.Printf("图片数量输入错误: %s\n", err.Error())
		return
	}

	for l := range productNameList {
		productName := string(productNameList[l])

		productIdList := Spider.GetProductList(productName, productCount)
		for i := range productIdList {
			imageUrlList := Spider.GetCommentImageList(productIdList[i], singleImageCount)
			for j := range imageUrlList {
				path := "./data/" + productName + "/"
				Spider.Mkdir(path)
				Spider.SaveImage(path, imageUrlList[j])
				time.Sleep(time.Millisecond * 300)
			}
		}
		fmt.Printf("商品%s爬取完成\n", productName)
	}

	fmt.Printf("商品全部爬取完成, 回车关闭窗口!\n")
	_, _ = fmt.Scanln()
}
