package main

import (
	"dbcore/core"
	"dbcore/pbfiles"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	core.InitConfig() // 初始化系统配置
	core.InitDB()     // 初始化DB配置

	server := grpc.NewServer()
	pbfiles.RegisterDBServiceServer(server, &core.DBService{})
	// 监听8080
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
