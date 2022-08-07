// Very basic socket server
// https://golangr.com/

package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	gen "nikworkedprofile/GoApi/ServerGo/src/generate_logs"
	logs "nikworkedprofile/GoApi/ServerGo/src/logs_app"
	"nikworkedprofile/GoApi/ServerGo/src/web"
)

var (
	ipaddr []string
	empty  []string

	mail string = "Succes"

	ctxVFS, cancelVFS = context.WithCancel(context.Background())
	ctxWEB, cancelWEB = context.WithCancel(context.Background())
	ctxGEN, cancelGEN = context.WithCancel(context.Background())
	sigc              = make(chan os.Signal, 1)
)

func echoServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			logs.FatalLogger.Println("Read buf serv" + err.Error())

			return
		}

		data := buf[0:nr]
		println("Server got:", string(data))
		s := strings.TrimSpace(string(data))
		/* if s == "VFS" {
			MesToClient(c, "Выбрана служба vfs\n")
			logs.InfoLogger.Println("Enter VFS service")
			go controllers.VFS("15000", ctxVFS)
		} */
		if s == "WEB" {
			MesToClient(c, "Выбрана служба web\n")
			logs.InfoLogger.Println("Enter web service")

			allip := enterIp(c)
			go web.ProcWeb("-x", allip, ctxWEB)
		}
		if s == "STOPWEB" {
			MesToClient(c, "Остановыка службы web\n")
			logs.InfoLogger.Println("Stop web service")

			go func() {
				cancelWEB()
				fmt.Println("stop WEB")
				ctxWEB, cancelWEB = context.WithCancel(context.Background())
			}()

		}

		/* if s == "STOPVFS" {
			MesToClient(c, "Остановыка службы vfs\n")
			go func() {
				cancelVFS()
				//fmt.Println("stop VFS")
				logs.InfoLogger.Println("Stop VFS service")
				ctxVFS, cancelVFS = context.WithCancel(context.Background())
			}()
			//cancel()
		} */
		if s == "STOP SERVER" {
			MesToClient(c, "Остановыка сервера\n")
			logs.ErrorLogger.Println("Stop server: ")
			//syscall.Kill()
			os.Exit(0)
		}
		//
		data = []byte(mail) //Send Client
		_, err = c.Write(data)
		if err != nil {
			logs.ErrorLogger.Println("Write writes data to the connection" + err.Error())
			log.Print("Write: ", err)
		}
		if s == "GENERATE LOG FILES" {
			MesToClient(c, "GENERATE LOG FILES\n")

			logs.InfoLogger.Println("Enter generate log files 10,20000000")
			gen.RemoveByConfig()
			go gen.ProcGenN(10, 200000)
		}
		if s == "GENERATE LOG FILE" {
			MesToClient(c, "GENERATE LOGS FILE\n")
			logs.InfoLogger.Println("Enter generate logs in file")
			go gen.ProcGenWriteF(ctxGEN)
		}
		if s == "STOP GEN LOG FILE" {
			MesToClient(c, "STOP GENERATE LOGS FILE\n")
			logs.InfoLogger.Println("Stop generate logs in file")
			go func() {
				cancelGEN()
				fmt.Println("stop Gen")
				ctxGEN, cancelGEN = context.WithCancel(context.Background())
			}()
		}
		if s == "VIEW CONFIG FILE" {
			//MesToClient(c, "VIEW CONFIG FILE\n")
			//logs.InfoLogger.Println("Enter VIEW CONFIG FILE")
			//content, err := ioutil.ReadFile("config.toml")
			//gen.ProcGenWriteF()
		}

	}
}

func Server() string {
	fmt.Println("Server start")
	go web.ProcWeb("-x", empty, ctxWEB)
	go func() {
		log.Println("App running, press CTRL + C to stop")
		select {}
	}()

	files, err := ioutil.ReadDir("/tmp/")
	if err != nil {
		logs.FatalLogger.Println("ReadDir" + err.Error())

		log.Fatal(err)
	}

	for _, f := range files {
		if f.Name() == "echo.sock" {
			os.Remove("/tmp/echo.sock")
			logs.FatalLogger.Println("Find echo.suck")
			log.Fatal("FIND echo.sock ")

		}
	}

	l, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		logs.FatalLogger.Println("Listen echo.sock" + err.Error())
		log.Fatal("listen error:", err)
	}
	//sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		// Wait for a SIGINT or SIGKILL:
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		// Stop listening (and unlink the socket if unix type):
		l.Close()
		// And we're done:
		os.Exit(0)
	}(sigc)
	for {
		fd, err := l.Accept()
		if err != nil {
			logs.FatalLogger.Println("Accept error echo.sock" + err.Error())
			log.Fatal("accept error:", err)
		}
		//shutdown.Listen()
		go echoServer(fd)

	}

}

func MesToClient(c net.Conn, message string) {
	data := []byte(message + "\n") //Send Client
	_, err := c.Write(data)
	if err != nil {
		logs.FatalLogger.Println("Message to client" + err.Error())

		log.Fatal("Write: ", err)
	}

}
