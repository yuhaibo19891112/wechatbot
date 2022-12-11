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

}

var Config *Configuration
var once sync.Once

func Timer() {
	SetToLocal()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		SetToLocal()
	}
}


func SetToLocal() {
	LoadConfig()

	configBody := RemoteConfigHttp(Config.RemoteUrl)

	botNum := fetchSetUpParam()
	
	if configBody == nil || botNum == "" {
		return
	}

	if "1" == botNum {
		Config.ApiKey = configBody.ApiKey1
	}
	if "2" == botNum {
		Config.ApiKey = configBody.ApiKey2
	}
	if "3" == botNum {
		Config.ApiKey = configBody.ApiKey3
	}
	if "4" == botNum {
		Config.ApiKey = configBody.ApiKey4
	}
	if "5" == botNum {
		Config.ApiKey = configBody.ApiKey5
	}
	if "6" == botNum {
		Config.ApiKey = configBody.ApiKey6
	}
	if "7" == botNum {
		Config.ApiKey = configBody.ApiKey7
	}
	if "8" == botNum {
		Config.ApiKey = configBody.ApiKey8
	}

	Config.AutoPass = configBody.AutoPass
	Config.FilterName = configBody.FilterWords
	if ""!= configBody.GroupName {
		Config.AlarmGroupName = configBody.GroupName
	}

	if ""!= configBody.UserRemarkName {
		Config.AlarmUserName = configBody.UserRemarkName
	}
}

// 获取启动参数
func fetchSetUpParam() string {
	for i ,v := range os.Args {
		if 1 == i {
			return v
		}
	}
	return "1"
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


// RemoteConfigResponseBody 请求体
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
}

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
	log.Println(string(respBody))
	err = json.Unmarshal(respBody, remoteConfigBody)

	if err != nil {
		return nil
	}

	return remoteConfigBody
}
