package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/parnurzeal/gorequest"
	"gorm.io/driver/mysql"
)

type ConsulAddr struct {
	ConsulAddr string `toml:"consul_addr"`
	ConfigPath string `toml:"config_path"`
}

type Config struct {
	App struct {
		Name       string `default:"app" toml:"name"`
		Port       uint16 `default:"8080" toml:"port"`
		Mode       string `default:"production" toml:"mode"` // allowed values: development, testing, production
		GatewayURL string `default:"" mapstructure:"gateway_url" toml:"gateway_url"`
	} `toml:"app"`
	MySQL  *mysql.Config `toml:"mysql"`
	Aliyun struct {
		Region          string `default:"" toml:"region"`
		AccessKeyID     string `default:"" mapstructure:"access_key_id" toml:"access_key_id"`
		AccessKeySecret string `default:"" mapstructure:"access_key_secret" toml:"access_key_secret"`
		OSS             struct {
			Bucket       string `default:"" toml:"bucket"`
			Endpoint     string `default:"" toml:"endpoint"`
			CustomDomain string `default:"" mapstructure:"custom_domain" toml:"custom_domain"`
			AccessKey    string `default:"" mapstructure:"access_key" toml:"access_key"`
			AccessSecret string `default:"" mapstructure:"access_secret" toml:"access_secret"`
		} `toml:"oss"`
	} `toml:"aliyun"`
}

var cfg Config

func GetConfig() *Config {
	return &cfg
}

func (c *Config) IsProductionMode() bool {
	return c.App.Mode == "production"
}

func (c *Config) IsDevelopmentMode() bool {
	return c.App.Mode == "development"
}

func (c *Config) IsTestingMode() bool {
	return c.App.Mode == "testing"
}

func InitConfig(config *ConsulAddr) {

	if err := loadConfigFromRemote(config.ConsulAddr+config.ConfigPath, &cfg); err != nil {
		log.Fatal(err)
	}

	if cfg.App.GatewayURL == "" {
		cfg.App.GatewayURL = fmt.Sprintf("http://localhost:%d", cfg.App.Port)
	} else {
		cfg.App.GatewayURL = strings.TrimRight(cfg.App.GatewayURL, "/")
	}
}

type ConsulConfig struct {
	LockIndex   int    `json:"LockIndex"`
	Key         string `json:"Key"`
	Flags       int    `json:"Flags"`
	Value       string `json:"Value"`
	CreateIndex int    `json:"CreateIndex"`
	ModifyIndex int    `json:"ModifyIndex"`
}

func loadConfigFromRemote(consulApi string, dst interface{}) error {
	_, body, errs := gorequest.New().Get(consulApi).EndBytes()
	if errs != nil {
		return errs[0]
	}

	consulCfg := []ConsulConfig{}
	if err := json.Unmarshal(body, &consulCfg); err != nil {
		return err
	}

	cs, err := base64.StdEncoding.DecodeString(consulCfg[0].Value)
	if err != nil {
		return err
	}

	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return errors.New("dst must be ptr")
	}
	switch reflect.TypeOf(dst).Elem().Kind().String() {
	case "struct":
		if _, err := toml.Decode(string(cs), dst); err != nil {
			return err
		}
	case "string":
		d := dst.(*string)
		*d = string(cs)
	default:
		return errors.New("unknow dst type")
	}

	return nil
}
