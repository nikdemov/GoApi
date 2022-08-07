package web

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hpcloud/tail"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"nikworkedprofile/GoApi/ServerGo/src/generate_logs"
	"nikworkedprofile/GoApi/ServerGo/src/logenc"
	logs "nikworkedprofile/GoApi/ServerGo/src/logs_app"
	"nikworkedprofile/GoApi/ServerGo/src/web/controllers"
	"nikworkedprofile/GoApi/ServerGo/src/web/util"

	"github.com/alecthomas/kingpin"
	"github.com/gorilla/mux"
	"github.com/shurcooL/httpfs/union"
	"github.com/spf13/afero"
)

var (
	//dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default("./view").ExistingFilesOrDirs()
	//dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default("/home/nik/projects/Course/logi2/repdata/").ExistingFilesOrDirs()
	//dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default("/home/nik/projects/Course/tmcs-log-agent-storage/").ExistingFilesOrDirs()
	dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default(pathdata + "/repdata").ExistingFilesOrDirs()
	//port = kingpin.Flag("port", "Port number to host the server").Short('p').Default("15000").Int()
	//port            *int
	cron            = kingpin.Flag("cron", "configure cron for re-indexing files, Supported durations:[h -> hours, d -> days]").Short('t').Default("0h").String()
	missadr         []string
	limit           string
	ipaddr          []string
	wg              sync.WaitGroup
	status          bool = false
	ctxCF, cancelCF      = context.WithCancel(context.Background())
	pathdata             = "/var/local/logi2"
	search          string
	datestartend    string
	savefiles       []string
	stringF         bool
	quit            (chan bool)
)

/* var (
	uptime_server_web = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "logi2_uptime_server_web_seconds",
		Help: "How many time server run",
	})
) */

type DatabaseConfig struct {
	Host  []string `mapstructure:"hostname"`
	Hostt []string `toml:"hostname"`
	Port  string
}

type Config struct {
	Db       DatabaseConfig `mapstructure:"database"`
	DataBase DatabaseConfig `toml:"database"`
}

func ProcWeb(dir1 string, slice []string, ctx context.Context) (err error) {
	//startTime := time.Now()

	status = false
	if dir1 == "-x" {
		status = true
	}

	generate_logs.Remove(pathdata+"/testsave/", "gen_logs_coded")
	generate_logs.Remove(pathdata+"/testsave/", "md5")

	logenc.CreateDir(pathdata + "/repdata/")
	logenc.CreateDir(pathdata + "/testsave/")

	kingpin.Parse()

	go func() {
		for {

			err = util.ParseConfig(*dir, *cron) //INDEXING FILE

			if err != nil {
				logs.FatalLogger.Println("Indexing files on web" + err.Error())
				panic(err)

			}
			time.Sleep(time.Second * 10)
		}
	}()
	time.Sleep(time.Second * 5)
	for i := 0; i < len(util.Conf.Dir); i++ {
		filesList = append(filesList, files{util.Conf.Dir[i]})
	}
	infoip, err := ioutil.ReadFile(pathdata + "/config.toml")
	if err != nil {

		log.Fatal(err)
	}
	fmt.Println("ip", string(infoip))
	/* go func() {
		time.Sleep(time.Second * 10)
		for {
			uptime := time.Since(startTime)
			uptime_server_web.Set(float64(uptime) / float64(time.Second))
		}
	}() */
	go util.DiskInfo(pathdata + "/repdata")
	EnterIpReady(slice)
	/* if status {
		EnterIpReady(slice)
	} else {
		EnterIp()
	} */

	Ip, CPort := CheckConfig()
	for i := 0; i < len(Ip); i++ {

		go CheckFiles(Ip[i], CPort, ctxCF)
		time.Sleep(time.Second * 2)
	}
	/* <-ctx.Done()

	cancelCF()
	ctxCF, cancelCF = context.WithCancel(context.Background()) */

	//go CheckFiles("localhost", "10015", ctxCF)
	//vfs
	dir := pathdata + "/repdata/"
	time.Sleep(time.Second * 10)
	fsbase := afero.NewBasePathFs(afero.NewOsFs(), dir)
	fsInput := afero.NewReadOnlyFs(fsbase)
	fsRoot := union.New(map[string]http.FileSystem{
		"/data": afero.NewHttpFs(fsInput),
	})

	router := mux.NewRouter()
	fileserver := http.FileServer(fsRoot)

	router.HandleFunc("/ws/{b64file}", Use(controllers.WSHandler)).Methods("GET")
	router.HandleFunc("/", Use(controllers.RootHandler)).Methods("GET")
	router.HandleFunc("/searchproject", controllers.SearchHandler)
	router.HandleFunc("/datestartend", controllers.DataHandler)
	router.HandleFunc("/pointproject", controllers.PointHandler)
	router.PathPrefix("/vfs/").Handler(http.StripPrefix("/vfs/", fileserver))
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/files/{fetchPercentage}", Use(filesL)).Methods("GET")
	router.HandleFunc("/data/{b64file}", filedata).Methods("GET")

	//router.PathPrefix("/metrics/").Handler(promhttp)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/static")))
	//router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/static")))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.tmpl")
	})
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", 15000), //port
		Handler: router}
	go func() {
		if err = server.ListenAndServe(); err != nil {
			logs.ErrorLogger.Println("Listen service" + err.Error())
			log.Println("listen:", err)
		}

	}()
	<-ctx.Done()
	go func() {
		cancelCF()
		ctxCF, cancelCF = context.WithCancel(context.Background())
	}()

	log.Printf("server stopped")

	if err = server.Shutdown(context.Background()); err != nil {
		logs.FatalLogger.Println("Shutdown Web" + err.Error())
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}
func rootPage(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("This is root page"))
}

func filesL(w http.ResponseWriter, r *http.Request) {

	fetchPercentage, errInput := strconv.ParseFloat(mux.Vars(r)["fetchPercentage"], 64)

	fetchCount := 0

	if errInput != nil {
		fmt.Println("gg ", errInput.Error())
	} else {
		fetchCount = int(float64(len(filesList)) * fetchPercentage / 100)
		if fetchCount > len(filesList) {
			fetchCount = len(filesList)
		}
	}

	// write to response
	jsonList, err := json.Marshal(filesList[0:fetchCount])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("content-type", "application/json")
		w.Write(jsonList)
	}

}

func filedata(w http.ResponseWriter, r *http.Request) {
	filenameB, _ := base64.StdEncoding.DecodeString(mux.Vars(r)["b64file"])
	filename := string(filenameB)
	if filenameB == nil {
		return
	}

	/* if filename == "undefined" {
		ViewDir(conn, search)
	}
	*/
	if savefiles == nil {
		//Indexing(filename)
		savefiles = append(savefiles, filename)
	} else {
		for i := 0; i < len(savefiles); i++ {
			if filename != savefiles[i] {
				stringF = true
			} else {
				stringF = false
			}
		}

	}
	if stringF {
		//Indexing(filename)
		savefiles = append(savefiles, filename)

	}

	// sanitize the file if it is present in the index or not.

	filename = filepath.Clean(filename)
	ok := false
	for _, wFile := range util.Conf.Dir {
		if filename == wFile {
			ok = true
			break
		}
	}
	if ok {
		var listdata []string

		listdata = append(listdata, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><catalog>")
		fileN := filepath.Base(filename)
		t, err := tail.TailFile("/var/local/logi2/repdata/"+fileN,
			tail.Config{
				//ReOpen: true,
				Follow: false,
				Location: &tail.SeekInfo{
					//Offset: current,
					Whence: io.SeekStart, //!!!

				},
			})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error occurred in opening the file: ", err)
		}
		for line := range t.Lines {
			xmlsimple := logenc.ProcLineDecodeXML(line.Text)
			//util.EncodeXML(xmlsimple)
			listdata = append(listdata, logenc.EncodeXML(xmlsimple))
		}
		listdata = append(listdata, "</catalog>")
		result2 := strings.Join(listdata, " ")

		w.Header().Set("content-type", "application/xml")
		w.Write([]byte(result2))
	}

}

type files struct {
	Path string
}

var filesList = []files{}

func FileGet(fileName string, w http.ResponseWriter) {
	t, err := tail.TailFile(fileName,
		tail.Config{
			Follow: true,
			Location: &tail.SeekInfo{
				Whence: os.SEEK_END,
			},
		})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error occurred in opening the file: ", err)
	}
	for line := range t.Lines {
		xmlsimple := logenc.ProcLineDecodeXML(line.Text)
		w.Header().Set("content-type", "application/text")
		w.Write([]byte(logenc.EncodeXML(xmlsimple)))
	}
}
