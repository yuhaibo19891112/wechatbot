package gtp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/869413421/wechatbot/config"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const BASEURL = "https://api.openai.com/v1/"

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     int    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {
	
	if filterWords(msg) {
		return "", errors.New("error words")
	}

    // 加一个中文符号将提问强制结束，不做"接话机器"
    // 测试过，即使标点重复或者错误，语义不受影响。
	msg = msg + "。"
	
	requestBody := ChatGPTRequestBody{
		Model:            "text-davinci-003",
		Prompt:           msg,
		MaxTokens:        1024,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)

	tr := &http.Transport{
		MaxIdleConns: 300,
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, time.Second*2) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			err = conn.SetDeadline(time.Now().Add(time.Second * 16)) //设置发送接受数据超时
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

	if err != nil {
		return "", err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL+"completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := config.Config.ApiKey
	apiKeyArray := strings.Split(apiKey, ",")
	randNum := rand.Intn(len(apiKeyArray))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ apiKeyArray[randNum])
	client := &http.Client{Transport: tr}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		return strconv.Itoa(response.StatusCode), errors.New(fmt.Sprintf("gtp api status code not equals 200,code is %d, msg %s", response.StatusCode, string(body)))
	}
	if err != nil {
		return "", err
	}

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}

	var reply string
	if len(gptResponseBody.Choices) > 0 {
		reply = gptResponseBody.Choices[0].Text
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}

func filterWords(msg string) bool  {
	if strings.Contains(msg, "台湾") ||
		strings.Contains(msg, "习大大") ||
		strings.Contains(msg, "习近平") ||
		strings.Contains(msg, "邓小平") ||
		strings.Contains(msg, "毛泽东") ||
		strings.Contains(msg, "冰毒") ||
		strings.Contains(msg, "政治") ||
		strings.Contains(msg, "推翻") ||
		strings.Contains(msg, "台灣") ||
		strings.Contains(msg, "習大大") {
		return true
	}
	remoteWords := config.Config.FilterName
	if remoteWords == "" || strings.TrimSpace(remoteWords) == ""{
		return false
	}
	words := strings.Split(remoteWords, ",")
	if len(words) == 0 {
		return false
	}
	for i := 0; i < len(words); i++ {
		if strings.TrimSpace(words[i]) != "" && strings.Contains(msg, words[i]) {
			return true
		}
	}
	return false
}

