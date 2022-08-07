package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	//"encoding/json"

	"html/template"
	"net/http"
	"path/filepath"

	"nikworkedprofile/GoApi/ServerGo/src/bleveSI"
	"nikworkedprofile/GoApi/ServerGo/src/logenc"
	logs "nikworkedprofile/GoApi/ServerGo/src/logs_app"
	"nikworkedprofile/GoApi/ServerGo/src/web/util"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	pathdata = "/var/local/logi2"
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	search        string
	datestartend  string
	savefiles     []string
	stringF       bool
	SearchMap     map[string]logenc.LogList
	date_layout   = "01/02/2006"
	startUnixTime int64
	endUnixTime   int64
	pointH        string
	filename      string
	PrevNetConn   net.Conn
)

type MyStruct struct {
	DirN string
	File string
}

// RootHandler - http handler for handling / path
func RootHandler(w http.ResponseWriter, _ *http.Request) {

	files := []string{
		"web/templates/index.tmpl",
		"web/templates/footer.tmpl",
		//"./ui/html/footer.partial.tmpl",
		"web/templates/header.tmpl",
		"web/templates/wscontent.tmpl",
		"web/templates/card.tmpl",
	}
	t := template.New("index").Delims("<<", ">>")

	t, err := t.Parse("footer")
	if err != nil {
		logs.FatalLogger.Println("Parse footer" + err.Error())
		log.Fatal("Problem with template \"footer\"")
	}
	t, err = t.Parse("header")
	if err != nil {
		logs.FatalLogger.Println("Parse header" + err.Error())
		log.Fatal("Problem with template \"header\"")
	}
	t, err = t.Parse("loading")
	if err != nil {
		logs.FatalLogger.Println("Parse loading" + err.Error())
		log.Fatal("Problem with template \"header\"")
	}
	t, err = t.Parse("wscontent")
	if err != nil {
		logs.FatalLogger.Println("Parse wscontentent" + err.Error())
		log.Fatal("Problem with template \"wscontent\"")
	}
	t, err = t.Parse("card")
	if err != nil {
		logs.FatalLogger.Println("Parse card" + err.Error())
		log.Fatal("Problem with template \"card\"")
	}
	t, err = t.ParseFiles(files...)
	t = template.Must(t, err)
	if err != nil {
		logs.FatalLogger.Println("Parse templates" + err.Error())
		panic(err)
	}
	var fileList = make(map[string]interface{})
	fileList["FileList"] = util.Conf.Dir

	t.Execute(w, fileList)
}

// WSHandler - Websocket handler
func WSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WSHandler .................................................")
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		logs.FatalLogger.Println("Open Websoket connection" + err.Error())
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	logs.InfoLogger.Println("Open Websoket connection")
	filenameB, _ := base64.StdEncoding.DecodeString(mux.Vars(r)["b64file"])

	filename = string(filenameB)
	if filenameB == nil {
		return
	}

	if filename == "undefined" {
		ViewDir(conn, search)
	}

	if savefiles == nil {
		Indexing(filename)
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
		Indexing(filename)
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

	//
	/* go func() { */
	/* 	PrevNetConn = conn.UnderlyingConn() */
	/* 	for { */
	/* 		if PrevNetConn != conn.UnderlyingConn() { */
	/* 			break //выход из цикла при создании другого канала соединения */
	/* 		} */
	/* 		if logenc.CheckFileSum(filename, "", "") { */
	/* 			Indexing(filename) */
	/* 		} */

	/* 	} */
	/* }() */

	// If the file is found, only then start tailing the file.
	// This is to prevent arbitrary file access. Otherwise send a 403 status
	// This should take care of stacking of filenames as it would first
	// be searched as a string in the index, if not found then rejected.

	if ok {
		logs.InfoLogger.Println("TailFile:" + filename + "Search:" + search)
		util.TailFile(conn, filename, search, SearchMap)

	}

	context.Clear(r)
}

var (
	countSearch = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logi2_search_reguest_count",
		Help: "How many search request",
	})
)

func SearchHandler(_ http.ResponseWriter, r *http.Request) {
	countSearch.Inc()
	search = r.URL.Query().Get("search_string")
	logs.InfoLogger.Println("Search handler: ", search)
	fmt.Println("SEARCHHANDLER:", search)
	context.Clear(r)
}

func DataHandler(_ http.ResponseWriter, r *http.Request) {
	datestartend = r.URL.Query().Get("daterange")
	fmt.Println("DATAHANDLER:", datestartend)
	//SEARCHHANDLER: 01/01/2021 - 01/15/2021
	if len(datestartend) != 0 {
		datastart := string(datestartend[0:10])
		// daystart := string(datestartend[0:2])
		//monthstart := string(datestartend[3:5])
		//yearstart := string(datestartend[6:10])

		//dayend := string(datestartend[13:15])
		//monthend := string(datestartend[16:18])
		//yearend := string(datestartend[19:23])
		//dataend := string(datestartend[13:23])

		timeendUnix, _ := time.Parse(date_layout, "01/15/2021")
		timestartUnix, _ := time.Parse(date_layout, datastart)

		//fmt.Println("Common", dataend, "Unix", timeendUnix.Unix())
		//fmt.Println("Common", datastart, "Unix", timestartUnix.Unix())

		//fmt.Println("Parse d:m:y", daystart, ":", monthstart, ":", yearstart)
		//fmt.Println("Parse d:m:y", dayend, ":", monthend, ":", yearend)
		startUnixTime = timestartUnix.Unix()
		endUnixTime = timeendUnix.Unix()
	}
	context.Clear(r)
}

func PointHandler(_ http.ResponseWriter, r *http.Request) {
	pointH = r.URL.Query().Get("pointer")
	fmt.Println("POINTHANDLER:", pointH)
	context.Clear(r)
}

//NOT fileUtils !!!
func Indexing(fileaddr string) {
	//var SearchMap map[string]string
	if fileaddr == "undefined" {
		return
	} else {
		fileN := filepath.Base(fileaddr)
		//logs.InfoLogger.Println("Indexing by bleve" + fileN)

		//fmt.Println(fileaddr)
		//logenc.Replication(fileaddr)
		//go func() {
		//conn.WriteMessage(websocket.TextMessage, []byte("Indexing file, please wait"))
		bleveSI.ProcBleve(fileN, fileaddr)
		//conn.WriteMessage(websocket.TextMessage, []byte("Indexing complated"))
		//}()
		SearchMap = logenc.ProcMapFile(fileaddr)
	}
}

//View List of Dir
//NOT fileUtils !!!
func ViewDir(conn *websocket.Conn, search string) {
	//Delete file with all
	//:TODO
	//Пагинация из первых ulid см как для отдельных файлов
	var fileList = make(map[string][]string)
	files, _ := ioutil.ReadDir(pathdata + "/repdata")
	//"/home/nik/projects/Course/tmcs-log-agent-storage/"
	//"./view"
	countFiles := (len(files))
	conn.WriteMessage(websocket.TextMessage, []byte("Indexing file, please wait"))
	logs.InfoLogger.Println("View Dir" + " Search:" + search)

	fileList["FileList"] = util.Conf.Dir
	fmt.Println("start")
	/*for i := 0; i < countFiles; i++ {
		fileName := fileList["FileList"][i]
		Indexing(fileName)
		fmt.Println(fileName)
	}
	*/
	go util.UlidPaginationDir(conn, countFiles, fileList)
	//CountPage := "<countpage>" + strconv.Itoa(10) + "</countpage>"
	//conn.WriteMessage(websocket.TextMessage, []byte(CountPage))
	util.TailDir2(conn, search, SearchMap, startUnixTime, endUnixTime, countFiles, fileList)
	conn.WriteMessage(websocket.TextMessage, []byte("Indexing complated"))
	startUnixTime = 0
	endUnixTime = 0

}
