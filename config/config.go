package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

const RESPONSE_SUCCESS = 1
const RESPONSE_LOGOUT = 4

const BASE_URL_IMAGE = "https://image.han-dress.cn/"

var IMAGE_DIR string

var RemoteDeepfaceApis = []string{
	"http://127.0.0.1:9000/analyze",
	"http://127.0.0.1:9001/analyze",
	"http://127.0.0.1:9002/analyze",
}

var InvalidDeepfaceApis []string

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	API      APIConfig      `yaml:"api"`
	Redis    RedisConfig    `yaml:"redis"`
	Session  SessionConfig  `yaml:"session"`
	Cookie   CookieConfig   `yaml:"cookie"`
}

type ServerConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	SSL        bool   `yaml:"ssl"`
	PidFile    string `yaml:"pid_file"`
	ImageDir   string `yaml:"image_dir"`
	SelfieRoot string `yaml:"selfie_root"`
	VideoDir   string `yaml:"video_dir"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

type APIConfig struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	EnableHTTPS bool   `yaml:"enable_https"`
}

type RedisConfig struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Password   string `yaml:"password"`
	Db         int    `yaml:"db"`
	Expiration uint   `yaml:"expiration"`
}

type SessionConfig struct {
	Name           string `yaml:"name"`
	GcMaxLifeTime  uint   `yaml:"gc_max_life_time"`
	PhpName        string `yaml:"php_name"`
	RedisKeyPrefix string `yaml:"redis_key_prefix"`
}

type CookieConfig struct {
	MaxAge   int    `yaml:"max_age"`
	Path     string `yaml:"path"`
	Domain   string `yaml:"domain"`
	Secure   bool   `yaml:"secure"`
	HttpOnly bool   `yaml:"http_only"`
}

var Conf string

func Init(config string) {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	Conf = config
	if Conf == "" {
		Conf = workDir + "/config/config.yaml"
	}
}
func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(Conf)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	IMAGE_DIR = cfg.Server.ImageDir
	return &cfg, nil
}
