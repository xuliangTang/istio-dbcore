package core

import (
	"dbcore/helpers"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

type DBConfig struct {
	DSN           string `yaml:"dsn"`
	MaxOpenConn   int    `yaml:"maxOpenConn"`
	MinIdleConn   int    `yaml:"maxIdleConn"`
	MaxLifeSecond int    `yaml:"maxLifeTime"` // 秒
}

type AppConfig struct {
	RpcPort  int `yaml:"rpcPort"`  //grpc端口
	HttpPort int `yaml:"httpPort"` //http端口
}

type SysConfigStruct struct {
	AppConfig *AppConfig `yaml:"appConfig"`
	DBConfig  *DBConfig  `yaml:"dbConfig"`
	APIList   []*API     `yaml:"apis"`
}

func (this *SysConfigStruct) FindAPI(name string) *API {
	for _, api := range this.APIList {
		if api.Name == name {
			return api
		}
	}
	return nil
}

var SysConfig *SysConfigStruct // 全局配置对象
const SysConfigPath = "./app.yml"

func InitConfig() { // 初始化配置
	SysConfig = &SysConfigStruct{}
	f, err := os.Open(SysConfigPath)
	helpers.Error(err, "找不到系统配置文件")
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, SysConfig)
	helpers.Error(err, "配置文件不正确")
}
