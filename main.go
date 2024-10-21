package main

import (
	"GinBlog/controller"
	"GinBlog/dao/mysql"
	"GinBlog/dao/redis"
	"GinBlog/logger"
	"GinBlog/pkg/snowflake"
	"GinBlog/router"
	"GinBlog/setting"
	"fmt"
	"os"
)

// @title GinBlog项目接口文档
// @version 1.0
// @description Go web博客项目

// @contact.name knoci
// @contact.email knoci@foxmail.com

// @host 127.0.0.1:8808
// @BasePath /api/v1

func main() {
	// 1.读取配置
	if len(os.Args) < 2 {
		fmt.Println("need config file.eg: GinBlog config.yaml")
		return
	}
	config_set := os.Args[1]
	if lenth := len(config_set); config_set[lenth-4:] == ".exe" {
		config_set = config_set[0 : lenth-4]
	}
	err := setting.Init(config_set)
	if err != nil {
		fmt.Printf("load setting failed: %v", err)
		return
	}
	// 2.加载日志
	err = logger.Init(setting.Conf.LogConfig, setting.Conf.Mode)
	if err != nil {
		fmt.Printf("load logger failed: %v", err)
		return
	}
	// 3.配置mysql
	err = mysql.Init(setting.Conf.MySQLConfig)
	if err != nil {
		fmt.Printf("init mysql failed: %v", err)
		return
	}
	defer mysql.Close()
	// 4.配置redis
	err = redis.Init(setting.Conf.RedisConfig)
	if err != nil {
		fmt.Println("init redis failed: %v", err)
		return
	}
	defer redis.Close()
	// 5.获取路由并运行服务
	if err := snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed: %v", err)
		return
	}

	if err := controller.InitTrans("zh"); err != nil {
		fmt.Println("init trans failed: %v", err)
		return
	}

	r := router.InitRouter(setting.Conf.Mode)
	err = r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed: %v", err)
		return
	}
}
