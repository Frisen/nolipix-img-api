package main

import (
	"flag"
	"fmt"
	"log"
	"nolipix-img-api/api"
	"nolipix-img-api/config"
	"nolipix-img-api/internal/aliyun"

	"github.com/BurntSushi/toml"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.toml", "config file path")
	flag.Parse()
	InitConfig(configPath)

	app := gin.New()
	app.MaxMultipartMemory = 20 << 20 //20M
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	app.Use(ginprom.PromMiddleware(nil))

	api.InitRouter(app)
	pprof.Register(app)
	if err := app.Run(fmt.Sprintf(":%d", config.GetConfig().App.Port)); err != nil {
		panic(err)
	}
}

func InitConfig(configPath string) {
	conf := new(config.ConsulAddr)
	_, err := toml.DecodeFile(configPath, conf)
	if err != nil {
		log.Fatal("decode config failed", err)
	}
	config.InitConfig(conf)
	aliyun.InitClient()
}
