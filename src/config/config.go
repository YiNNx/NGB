package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"ngb/util"
	"os"
)

var (
	// C 全局配置文件，在Init调用前为nil
	C *Config
)

// Config 配置
type Config struct {
	App        app        `yaml:"app"`
	Postgresql postgresql `yaml:"postgresql"`
	Jwt        jwt        `yaml:"jwt"`
	Log        log        `yaml:"log"`
}

type app struct {
	Addr string `yaml:"addr"`
}

type postgresql struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type jwt struct {
	Secret string `yaml:"secret"`
}

type log struct {
	Path string `yaml:"path"`
	File string `yaml:"file"`
}

func init() {
	configFile := "default.yml"

	if v, ok := os.LookupEnv("ENV"); ok {
		configFile = v + ".yml"
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("./env/config/%s", configFile))

	if err != nil {
		panic(err)
		return
	}

	config := &Config{}

	err = yaml.Unmarshal(data, config)

	if err != nil {
		fmt.Println("Unmarshal config error!")
		panic(err)
		return
	}

	C = config

	util.Logger.Info("Config " + configFile + " loaded.")
	fmt.Println("Config " + configFile + " loaded.")

}
