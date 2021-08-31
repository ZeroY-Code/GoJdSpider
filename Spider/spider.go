package Spider

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func LoadProduct(path string) []string {
	_, err := os.Lstat(path)
	if err != nil {
		log.Printf("标签列表文件 %s 不存在", path)
		os.Exit(0)
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0)
	br := bufio.NewReader(file)

	var productNameList []string
	for {
		productName, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		productNameList = append(productNameList, string(productName))
	}

	return productNameList
}

func Mkdir(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Printf("目录 %s 创建失败, 原因: %s", path, err.Error())
		return
	}
}

func SaveImage(savePath string, imageUrl string) {
	s := strings.Split(imageUrl, "/")
	fileName := s[len(s)-1]
	filePath := savePath + "/" + fileName
	if _, err := os.Lstat(filePath); err == nil {
		return
	}

	file, _ := os.Create(filePath)

	client := http.Client{}
	request, _ := http.NewRequest("GET", imageUrl, nil)
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("图片%s下载失败，原因: %s", fileName, err.Error())
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	length, err := file.WriteAt(body, 0)
	if err != nil {
		fmt.Printf("图片%s保存失败，原因: %s", fileName, err.Error())
		return
	}
	nowTime := time.Now().Format("2006-01-02 03:04:05")
	fmt.Printf("[%s] 文件%s保存成功，大小%db\n", nowTime, fileName, length)
}
