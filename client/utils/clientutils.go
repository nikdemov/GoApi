package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Db       DatabaseConfig `mapstructure:"database"`
	DataBase DatabaseConfig `toml:"database"`
}
type DatabaseConfig struct {
	Host  []string `mapstructure:"hostname"`
	Hostt []string `toml:"hostname"`
	Port  string
}

const pathdata = "/var/local/logi2"

func CheckConfig() ([]string, string) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(pathdata + "/")
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("couldn't load config:", err)
		//os.Exit(1)
	}
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}
	Ip := c.Db.Host
	Port := c.Db.Port
	Ip = removeDuplicateStr(Ip)

	return Ip, Port
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
