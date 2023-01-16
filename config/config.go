package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Configuration 项目配置
type Configuration struct {
	// 远程文件地址
	RemoteUrl string `json:"remote_url"`
	// 敏感词
	FilterName string `json:"filter_name"`
	// 图片链接
	QunUrl string `json:"remote_qun_url"`
	//进群提示语
	JoinGroupTip string `json:"join_group_tip"`
	// 接客语
	JiekeTip string `json:"jieke_tip"`
	// 系统用户
	SystemUser string `json:"system_user"`
	// 新闻早餐接口地址
	NewsApi string   `json:"news_api"`
}

// Config 公共参数
var Config *Configuration

// once 私有参数
var once sync.Once

// Timer 定时器
func Timer() {
	SetToLocal()

	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		SetToLocal()
	}
}

// SetToLocal 远程本地比较赋值
func SetToLocal() {

	LoadConfig()

	configBody := RemoteConfigHttp(Config.RemoteUrl)

	botNum := fetchSetUpParam()

	if configBody == nil || botNum == "" {
		return
	}

	// 机器人判断

	if Config.FilterName != configBody.FilterWords {
		log.Println("bot" + botNum + "：filterWords有变更! " + Config.FilterName + " ===> " + configBody.FilterWords)
		Config.FilterName = configBody.FilterWords
	}

	if "" != configBody.QunUrl && Config.QunUrl != configBody.QunUrl {
		log.Println("bot" + botNum + "：QunUrl有变更! " + Config.QunUrl + " ===> " + configBody.QunUrl)
		Config.QunUrl = configBody.QunUrl
	}

	if "" != configBody.JoinGroupTip && Config.JoinGroupTip != configBody.JoinGroupTip {
		log.Println("bot" + botNum + "：JoinGroupTip有变更! " + Config.JoinGroupTip + " ===> " + configBody.JoinGroupTip)
		Config.JoinGroupTip = configBody.JoinGroupTip
	}

	if "" != configBody.JiekeTip && Config.JiekeTip != configBody.JiekeTip {
		log.Println("bot" + botNum + "：JiekeTip有变更! " + Config.JiekeTip + " ===> " + configBody.JiekeTip)
		Config.JiekeTip = configBody.JiekeTip
	}

	if "" != configBody.NewsApi && Config.NewsApi != configBody.NewsApi {
		log.Println("bot" + botNum + "：NewsApi有变更! " + Config.NewsApi + " ===> " + configBody.NewsApi)
		Config.NewsApi = configBody.NewsApi
	}
}

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		Config = &Configuration{}
		f, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("open config err: %v", err)
			return
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(Config)
		if err != nil {
			log.Fatalf("decode config err: %v", err)
			return
		}
	})
	return Config
}

// RemoteConfigResponseBody 远程config.json参数体
type RemoteConfigResponseBody struct {
	GroupName      string `json:"group_name"`
	FilterWords    string `json:"filter_words"`
	UserRemarkName string `json:"user_remark_name"`
	//进群提示语
	JoinGroupTip string `json:"join_group_tip"`
	// 图片链接
	QunUrl string `json:"remote_qun_url"`
	// 接客语
	JiekeTip string `json:"jieke_tip"`
	// 新闻早餐接口地址
	NewsApi string   `json:"news_api"`
}

// RemoteConfigHttp 远程配置请求
func RemoteConfigHttp(remoteUrl string) *RemoteConfigResponseBody {

	resp, err := http.Get(remoteUrl)

	if err != nil {
		return nil
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil
	}

	if resp.StatusCode != 200 {
		return nil
	}

	remoteConfigBody := &RemoteConfigResponseBody{}
	err = json.Unmarshal(respBody, remoteConfigBody)

	if err != nil {
		return nil
	}

	return remoteConfigBody
}

// fetchSetUpParam 获取启动参数
func fetchSetUpParam() string {
	for i, v := range os.Args {
		if 1 == i {
			return v
		}
	}
	return "1"
}
