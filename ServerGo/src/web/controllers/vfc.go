package controllers

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	logs "nikworkedprofile/GoApi/ServerGo/src/logs_app"

	"github.com/gorilla/mux"
	"github.com/shurcooL/httpfs/union"
	"github.com/spf13/afero"
)

// RunHTTP run http api
func VFS(port string, ctx context.Context) (err error) {
	logs.InfoLogger.Println("Start VFS")

	fmt.Println("Start VFS")
	addr := ":" + port
	//dir := "/home/maxxant/Documents/log"
	//dir := "./tmcs-log-agent-storage/"
	dir := pathdata + "/repdata/"
	//dir := "/home/nik/projects/Course/tmcs-log-agent-storage/"

	var listener net.Listener
	//var err error
	listenErr := 0

	// wait for listening started
	for ok := false; !ok; {

		listener, err = net.Listen("tcp", addr)
		if err != nil {
			if listenErr == 0 {
				fmt.Println(err)
			}
			listenErr++
			time.Sleep(time.Second * 3)
		}
		ok = (err == nil)

		if ok {
			defer listener.Close()
		}
		//bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	//fmt.Println("listen ok: ", addr)

	fsbase := afero.NewBasePathFs(afero.NewOsFs(), dir)
	fsInput := afero.NewReadOnlyFs(fsbase)
	fsRoot := union.New(map[string]http.FileSystem{
		"/data": afero.NewHttpFs(fsInput),
	})

	router := mux.NewRouter()

	fileserver := http.FileServer(fsRoot)
	router.PathPrefix("/vfs/").Handler(http.StripPrefix("/vfs/", fileserver))
	//fmt.Println("running VFS" + " port: " + addr)
	//fmt.Println("Run new terminal for use service")

	srv := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		if err := srv.Serve(listener); err != nil {
			logs.ErrorLogger.Println("Http server error" + err.Error())
			fmt.Println("Http serve error", err)
		}
	}()
	<-ctx.Done()
	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		logs.FatalLogger.Println("VFC Shutdown Failed" + err.Error())
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}
