package aliyun

import (
	"fmt"
	"log"
	"nolipix-img-api/config"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var SHClient *oss.Bucket
var initOnce sync.Once

func InitClient() {
	conf := config.GetConfig()
	fmt.Println(conf.Aliyun)
	client, err := oss.New(conf.Aliyun.OSS.Endpoint, conf.Aliyun.OSS.AccessKey, conf.Aliyun.OSS.AccessSecret)
	if err != nil {
		log.Fatal(err)
	}
	SHClient, err = client.Bucket(conf.Aliyun.OSS.Bucket)
	if err != nil {
		log.Fatal(err)
	}
}

func GetSHClient() *oss.Bucket {
	if SHClient == nil {
		initOnce.Do(InitClient)
	}
	return SHClient
}
