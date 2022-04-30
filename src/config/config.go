package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

	// 读取配置文件
	data, err := ioutil.ReadFile(fmt.Sprintf("./env/config/default.yml"))

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
	fmt.Println("Config loaded.")

}

