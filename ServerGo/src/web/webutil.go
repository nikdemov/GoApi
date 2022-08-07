package web

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	logs "nikworkedprofile/GoApi/ServerGo/src/logs_app"
	"nikworkedprofile/GoApi/ServerGo/src/web/util"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

// Use - Stacking middlewares
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}

func CheckFiles(address string, port string, ctx context.Context) {
	//quit := make(chan bool)
	for range time.Tick(time.Second * 2) {

		if len(missadr) == 0 {
			missadr = append(missadr, "nope")
		}
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stop CheckFiles")
				return
			default:
				for i := range missadr {
					if missadr[i] != address {
						err := util.GetFiles(address, port)
						if err != nil {
							logs.ErrorLogger.Println("Getfiles" + err.Error())
							log.Println(err)
							fmt.Println(address)
							missadr = append(missadr, address)
						}
						err = util.ParseConfig(*dir, *cron) //INDEXING FILE

						if err != nil {
							logs.FatalLogger.Println("Parse config" + err.Error())
							log.Println("LOOP", err)
							panic(err)
						}

					}
				}

				fmt.Println(missadr)
				//wg.Add(1)
				go reconect(address)
				//wg.Wait()
				time.Sleep(time.Second * 10)
				continue

			}
		}
	}

}

func reconect(address string) {
	//defer wg.Done()

	for i := range missadr {
		if missadr[i] == address {
			copy(missadr[i:], missadr[i+1:]) // Shift a[i+1:] left one index.
			missadr[len(missadr)-1] = ""     // Erase last element (write zero value).
			missadr = missadr[:len(missadr)-1]
		}

	}

}

func EnterIp() {
	var data []byte
	for {

		fmt.Print("Enter IP:  ")
		fmt.Scanln(&limit)

		if limit == "stop" {
			//ipaddr = append(ipaddr, "localhost")
			limitSlice, _ := CheckConfig()
			ipaddr = append(ipaddr, limitSlice...)
			ipaddr = removeDuplicateStr(ipaddr)
			config := Config{DataBase: DatabaseConfig{Hostt: ipaddr, Port: "15000"}}
			data, _ = toml.Marshal(&config)
			break
		} else if util.CheckIPAddress(limit) {
			ipaddr = append(ipaddr, limit)
			limitSlice, _ := CheckConfig()
			ipaddr = append(ipaddr, limitSlice...)
			ipaddr = removeDuplicateStr(ipaddr)
			config := Config{DataBase: DatabaseConfig{Hostt: ipaddr, Port: "15000"}}

			data, _ = toml.Marshal(&config)
		}
	}

	err3 := ioutil.WriteFile(pathdata+"/config.toml", data, 0666)

	if err3 != nil {

		log.Fatal(err3)
	}
	fmt.Println("Written")

}

func CheckConfig() ([]string, string) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(pathdata + "/")
	if err := v.ReadInConfig(); err != nil {
		logs.ErrorLogger.Println("Read config" + err.Error())
		fmt.Println("couldn't load config:", err)
		//os.Exit(1)
	}
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		logs.ErrorLogger.Println("Unmarchal config" + err.Error())
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

func EnterIpReady(ipmas []string) {

	var data []byte
	ipaddr = ipmas
	limitSlice, _ := CheckConfig()
	ipaddr = append(ipaddr, limitSlice...)
	ipaddr = removeDuplicateStr(ipaddr)
	config := Config{DataBase: DatabaseConfig{Hostt: ipaddr, Port: "15000"}}
	data, _ = toml.Marshal(&config)

	err3 := ioutil.WriteFile(pathdata+"/config.toml", data, 0666)

	if err3 != nil {

		log.Fatal(err3)
	}
	fmt.Println("Written")

}

var startTime time.Time

func uptime() time.Duration {
	return time.Since(startTime)
}
