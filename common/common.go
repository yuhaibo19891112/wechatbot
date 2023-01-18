package common

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func LoadRemoteImg(finalUrl string, savePath string) (*os.File, error) {
	//读取url的信息，存入到文件
	resp, err := http.Get(finalUrl)
	if err != nil {
		log.Printf("get remote img error, %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	content, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("read img error, %v", readErr)
		return nil, err
	}
	ioutil.WriteFile(savePath, content, 0666)
	return os.Open(savePath)
}
