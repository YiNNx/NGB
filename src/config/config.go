package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

func init() {
	configFile := "default.yml"

	// 如果有设置 ENV ，则使用ENV中的环境
	if v, ok := os.LookupEnv("ENV"); ok {
		configFile = v + ".yml"
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(fmt.Sprintf("./env/config/%s", configFile))

	if err != nil {
		//Logger.Println("Read config error!")
		//Logger.Panic(err)
		panic(err)
		return
	}

	config := &Config{}

	err = yaml.Unmarshal(data, config)

	if err != nil {
		//Logger.Println("Unmarshal config error!")
		//Logger.Panic(err)
		fmt.Println("Unmarshal config error!")
		panic(err)
		return
	}

	C = config

	//Logger.Println("Config " + configFile + " loaded.")
	fmt.Println("Config " + configFile + " loaded.")

}
