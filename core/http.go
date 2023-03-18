package core

import (
	"fmt"
	"log"
	"net/http"
)

func InitHttp() {
	go func() {
		http.HandleFunc("/reload", func(writer http.ResponseWriter, request *http.Request) {
			InitConfig()
		})
		log.Println(fmt.Sprintf("启动内部服务,端口是:%d", SysConfig.AppConfig.HttpPort))
		http.ListenAndServe(fmt.Sprintf(":%d", SysConfig.AppConfig.HttpPort), nil)
	}()
}
