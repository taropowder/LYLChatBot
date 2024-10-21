package conf

import (
	"github.com/eatmoreapple/openwechat"
	log "github.com/sirupsen/logrus"
)

type DatabaseConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
}

type ApiAuthConfig struct {
	Token  string `yaml:"token"`
	Source string `yaml:"source"`
}

type AbstractHandlerConfig struct {
	Cookie string `yaml:"cookie"`
}

type Config struct {
	LogLevel     log.Level `yaml:"log_level"`
	GlobalUsers  []string  `yaml:"global_users"`
	GlobalGroups []string  `yaml:"global_groups"`
	WaitTime     int       `yaml:"wait_time"`
	Handlers     struct {
		Abstract AbstractHandlerConfig `yaml:"abstract"`
	} `yaml:"handlers"`
	Gpt struct {
		NewBing NewBingConfig `yaml:"new_bing"`
		Qwen    QwenConfig    `yaml:"qwen"`
		Gemini  GeminiConfig  `yaml:"gemini"`
	} `yaml:"gpt"`
	Redis       RedisConfig     `yaml:"redis"`
	Database    DatabaseConfig  `yaml:"database"`
	APIAuth     []ApiAuthConfig `yaml:"api_auth"`
	APIAddress  string          `yaml:"api_address"`
	Proxy       string          `yaml:"proxy"`
	HowNewsSite string          `yaml:"how_news_site"`
	HowNewsApi  string          `yaml:"how_news_api"`
	Google      struct {
		ApiKey string `yaml:"api_key"`
		CX     string `yaml:"cx"`
	} `yaml:"google"`
}

type NewBingConfig struct {
	URL        string            `yaml:"url"`
	Proxy      string            `yaml:"proxy"`
	RetryTimes int               `yaml:"retry_times"`
	Timeout    int               `yaml:"time_out"`
	Headers    map[string]string `yaml:"headers"`
}

type QwenConfig struct {
	ApiKey     string `yaml:"api_key"`
	Model      string `yaml:"model"`
	RetryTimes int    `yaml:"retry_times"`
	Timeout    int    `yaml:"time_out"`
}

type GeminiConfig struct {
	ApiKey string `yaml:"api_key"`
	Proxy  string `yaml:"proxy"`
}

var ConfigureInstance = Config{}
var BotInstance = &openwechat.Bot{}
var ConfigFilePath = ""
var DebugMode = false

func NewDefaultConfig() Config {
	return Config{
		LogLevel: log.InfoLevel,
		WaitTime: 60,
	}
}
