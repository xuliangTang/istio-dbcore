package cmd

import (
	"dbcore/core"
	"dbcore/pbfiles"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "dbcore",
	Short: "通用DB服务",
	Long:  `通用DB服务`,
	Run: func(cmd *cobra.Command, args []string) {
		core.InitConfig() // 初始化系统配置
		core.InitDB()     // 初始化DB配置
		core.InitHttp()

		server := grpc.NewServer()
		pbfiles.RegisterDBServiceServer(server, &core.DBService{})

		// 监听
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", core.SysConfig.AppConfig.RpcPort))
		log.Println(fmt.Sprintf("DBCore服务开始启动，监听端口:%d", core.SysConfig.AppConfig.RpcPort))
		if err != nil {
			log.Fatal(err)
		}
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(reloadCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version of dbcore",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.1.0")
	},
}

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "reload sys config (app.yaml)",
	Run: func(cmd *cobra.Command, args []string) {
		core.InitConfig()
		rsp, err := http.Get(fmt.Sprintf("http://localhost:%d/reload", core.SysConfig.AppConfig.HttpPort))
		if err != nil {
			log.Println(err)
		}
		defer rsp.Body.Close()
		if rsp.StatusCode == http.StatusOK {
			log.Println("重载配置文件成功")
		} else {
			log.Println("重载配置文件失败")
		}
	},
}
