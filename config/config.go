package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Configuration 项目配置
type Configuration struct {
	// gtp apikey
	ApiKey string `json:"api_key"`
	// 自动通过好友
	AutoPass bool `json:"auto_pass"`
	// 远程文件地址
	RemoteUrl string `json:"remote_url"`
	// 告警群 需要将该群保存到通讯录
	AlarmGroupName string `json:"group_name"`
	// 告警好友备注
	AlarmUserName string `json:"user_remark_name"`
	// 敏感词
	FilterName string `json:"filter_name"`
	// 图片链接
	QunUrl string `json:"remote_qun_url"`
	//进群提示语
	JoinGroupTip string `json:"join_group_tip"`

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
	if "1" == botNum && Config.ApiKey != configBody.ApiKey1 {
		log.Println("bot"+ botNum +"：apikey有变更! " + Config.ApiKey + " ===> " + configBody.ApiKey1)
		Config.ApiKey = configBody.ApiKey1
	}
	if "2" == botNum && Config.ApiKey != configBody.ApiKey2 {
		log.Println("bot"+ botNum +"：apikey有变更! " + Config.ApiKey + " ===> " + configBody.ApiKey2)
		Config.ApiKey = configBody.ApiKey2
	}
	if "3" == botNum && Config.ApiKey != configBody.ApiKey3 {
		log.Println("bot"+ botNum +"：apikey有变更! " + Config.ApiKey + " ===> " + configBody.ApiKey3)
		Config.ApiKey = configBody.ApiKey3
	}
	if "4" == botNum && Config.ApiKey != configBody.ApiKey4 {
		log.Println("bot"+ botNum +"：apikey有变更! " + Config.ApiKey + " ===> " + configBody.ApiKey4)
		Config.ApiKey = configBody.ApiKey4
	}
	if "5" == botNum && Config.ApiKey != configBody.ApiKey5 {
		log.Println("bot"+ botNum +"：apikey有变更! " + Config.ApiKey + " ===> " + configBody.ApiKey5)
		Config.ApiKey = configBody.ApiKey5
	}
	if "6" == botNum && Config.ApiKey != configBody.ApiKey6 {
		log.Println("bot"+ botNum +"：apikey有变更! " + Config.ApiKey + " ===> " + configBody.ApiKey6)
		Config.ApiKey = configBody.ApiKey6
	}
	if "7" == botNum && Config.ApiKey != configBody.ApiKey7 {
		log.Println("bot"+ botNum +"：apikey有变更! " + Config.ApiKey + " ===> " + configBody.ApiKey7)
		Config.ApiKey = configBody.ApiKey7
	}
	if "8" == botNum && Config.ApiKey != configBody.ApiKey8 {
		log.Println("bot"+ botNum +"：apikey有变更! " + Config.ApiKey + " ===> " + configBody.ApiKey8)
		Config.ApiKey = configBody.ApiKey8
	}

	if Config.AutoPass != configBody.AutoPass {
		log.Println("bot"+ botNum +"：autoPass有变更! " + strconv.FormatBool(Config.AutoPass) + " ===> " + strconv.FormatBool(configBody.AutoPass))
		Config.AutoPass = configBody.AutoPass
	}

	if Config.FilterName != configBody.FilterWords {
		log.Println("bot"+ botNum +"：filterWords有变更! " + Config.FilterName + " ===> " + configBody.FilterWords)
		Config.FilterName = configBody.FilterWords
	}

	if ""!= configBody.GroupName && Config.AlarmGroupName != configBody.GroupName {
		log.Println("bot"+ botNum +"：alarmGroupName有变更! " + Config.AlarmGroupName + " ===> " + configBody.GroupName)
		Config.AlarmGroupName = configBody.GroupName
	}

	if ""!= configBody.UserRemarkName && Config.AlarmUserName != configBody.UserRemarkName {
		log.Println("bot"+ botNum +"：alarmUserName有变更! " + Config.AlarmUserName + " ===> " + configBody.UserRemarkName)
		Config.AlarmUserName = configBody.UserRemarkName
	}

	if ""!= configBody.QunUrl && Config.QunUrl != configBody.QunUrl {
		log.Println("bot"+ botNum +"：QunUrl有变更! " + Config.QunUrl + " ===> " + configBody.QunUrl)
		Config.QunUrl = configBody.QunUrl
	}

	if ""!= configBody.JoinGroupTip && Config.JoinGroupTip != configBody.JoinGroupTip {
		log.Println("bot"+ botNum +"：JoinGroupTip有变更! " + Config.JoinGroupTip + " ===> " + configBody.JoinGroupTip)
		Config.JoinGroupTip = configBody.JoinGroupTip
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

		// 如果环境变量有配置，读取环境变量
		ApiKey := os.Getenv("ApiKey")
		AutoPass := os.Getenv("AutoPass")
		if ApiKey != "" {
			Config.ApiKey = ApiKey
		}
		if AutoPass == "true" {
			Config.AutoPass = true
		}
	})
	return Config
}

// RemoteConfigResponseBody 远程config.json参数体
type RemoteConfigResponseBody struct {
	ApiKey1 string `json:"api_key-1"`
	ApiKey2 string `json:"api_key-2"`
	ApiKey3 string `json:"api_key-3"`
	ApiKey4 string `json:"api_key-4"`
	ApiKey5 string `json:"api_key-5"`
	ApiKey6 string `json:"api_key-6"`
	ApiKey7 string `json:"api_key-7"`
	ApiKey8 string `json:"api_key-8"`
	AutoPass bool `json:"auto_pass"`
	GroupName string `json:"group_name"`
	FilterWords string `json:"filter_words"`
	UserRemarkName string `json:"user_remark_name"`
	//进群提示语
	JoinGroupTip string `json:"join_group_tip"`
	// 图片链接
	QunUrl string `json:"remote_qun_url"`
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
	for i ,v := range os.Args {
		if 1 == i {
			return v
		}
	}
	return "1"
}
