package server

import (
	"log"
	"net"
	"strings"

	logs "nikworkedprofile/GoApi/src/logs_app"
	"nikworkedprofile/GoApi/src/web/util"
)

func enterIp(c net.Conn) []string {

	data := []byte("Input ip address for running service:\n Enter \"stop\" to run service") //Send Client
	_, err := c.Write(data)
	if err != nil {
		logs.FatalLogger.Println("Write writes data to the connection" + err.Error())

		log.Fatal("Write: ", err)
	}
	for {
		buf := make([]byte, 512)
		nr, _ := c.Read(buf)
		data := buf[0:nr]

		//data = []byte("В") //Send Client

		limit := string(data)
		limit = strings.TrimSpace(limit)
		if limit == "stop" {
			break
		} else if util.CheckIPAddress(limit) {
			data = []byte("Valid") //Send Client
			_, err := c.Write(data)
			if err != nil {
				logs.FatalLogger.Println("Write data on Enter Ip" + err.Error())
				log.Fatal("Write: ", err)
			}
			ipaddr = append(ipaddr, limit)
			ipaddr = removeDuplicateStr(ipaddr)

		} else {
			data = []byte("Invalid") //Send Client
			_, err := c.Write(data)
			if err != nil {
				logs.FatalLogger.Println("Invalid Write data on Enter Ip" + err.Error())
				log.Fatal("Write: ", err)
			}
		}
	}
	return ipaddr
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
