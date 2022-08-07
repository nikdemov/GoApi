package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"nikworkedprofile/GoApi/client/utils"

	"github.com/SCU-SJL/menuscreen"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		println("From server:", string(buf[0:n])) //From server
	}
}

func main() {

	utils.CallClear()
	fmt.Println("Now you can use only VFS(stable) or WEB(stable) and stop this service")
	c, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	go reader(c)
	for {
		//reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")

		//text, _ := reader.ReadString('\n') //Send server
		//ui terminal
		idx, text := menuClientMain()
		_, err := c.Write([]byte(text)) //Send server
		if err != nil {
			log.Fatal("write error:", err)
			break
		}
		switch idx {
		case 0:
			for {
				text = utils.WebMenu("0")

				_, err := c.Write([]byte(text)) //Send server
				if err != nil {
					log.Fatal("write error:", err)
					break
				}
				if text == "stop" {
					break
				}
				time.Sleep(1e9)
			}
		case 4:
			return
		case 7:
			utils.Edit()
		case 8:
			//test.Edit()
		}

		time.Sleep(1e9)
	}
}

func menuClientMain() (int, string) {
	menu, err := menuscreen.NewMenuScreen()
	if err != nil {
		panic(err)
	}
	defer menu.Fini()
	menu.SetTitle("ControlPanel").
		SetLine(0, "WEB").
		//SetLine(1, "VFS").
		SetLine(1, "STOPWEB").
		//SetLine(3, "STOPVFS").
		SetLine(2, "STOP CLIENT").
		SetLine(3, "STOP SERVER").
		SetLine(4, "GENERATE LOG FILES").
		SetLine(5, "GENERATE LOG FILE").
		SetLine(6, "STOP GEN LOG FILE").
		SetLine(7, "VIEW CONFIG FILE?").
		SetLine(8, "EDIT CONFIG FILE?").
		Start()
	idx, ln, ok := menu.ChosenLine()
	if !ok {
		fmt.Println("you did not chose any items.")
		return idx, ln

	}
	fmt.Printf("you've chosen %d line, content is: %s\n", idx, ln)
	return idx, ln
}
