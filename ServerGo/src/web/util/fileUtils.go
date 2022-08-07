package util

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"nikworkedprofile/GoApi/src/bleveSI"
	"nikworkedprofile/GoApi/src/logenc"
	logs "nikworkedprofile/GoApi/src/logs_app"

	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var paginationUlids map[int]string

var FileUlids map[string]string

//TDOD
var pathdata = "/var/local/logi2"
var countViweMes = 100
var ctpath = "/var/log/logi2"
var PrevNetConn net.Conn
var (
	FileName []string
	FileList []string
	visited  map[string]bool

	// Global Map that stores all the files, used to skip duplicates while
	// subsequent indexing attempts in cron trigger
	indexMap           = make(map[string]bool)
	signature   bool   = false
	Fname       string = ""
	currentfile string
	page        int = 0
	hashSumFile string
)

type FileStruct struct {
	ID      int    `json:"id"`
	NAME    string `json:"filename"`
	HASHSUM string `json:"hashsum"`
}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type Map map[string]string

// TailFile - Accepts a websocket connection and a filename and tails the
// file and writes the changes into the connection. Recommended to run on
// a thread as this is blocking in nature

func TailFile(conn *websocket.Conn, fileName string, lookFor string, SearchMap map[string]logenc.LogList) {
	fileN := filepath.Base(fileName)

	if Fname != fileName {
		if Fname != "" {
			logenc.DeleteOldsFiles(pathdata+"/replace/"+filepath.Base(Fname), "")
		}
		Fname = fileName
		lookFor = ""

		for k := range paginationUlids {
			delete(paginationUlids, k)
		}
	}

	currentfile = fileN
	page = 0
	if lookFor == "" || lookFor == " " || lookFor == "Search" {

		hashSumFile = logenc.FileMD5(fileName)
		go func() {
			for {
				if hashSumFile != logenc.FileMD5(fileName) {
					hashSumFile = logenc.FileMD5(fileName)
					fmt.Println("hashSumFile", hashSumFile)
				} else if Fname != fileName {
					break
				}
			}
		}()
		go followCodeStatus(conn)
		UlidPaginationFile(conn, fileName)

		logenc.DeleteOldsFiles(pathdata+"/replace/"+fileN, "")
		conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
		TailingLogsInFileAll(0, fileName, conn, 0, page)
		logenc.DeleteOldsFiles(pathdata+"/replace/"+fileN, "")

		var countline int = 0
		var currentpage int = 0

		for {
			if (logenc.FileMD5(fileName) != hashSumFile) && currentfile == fileN {
				UlidPaginationFile(conn, fileName)
				logenc.DeleteOldsFiles(pathdata+"/replace/"+fileN, "")
				conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
				countline = TailingLogsInFileAll(0, fileName, conn, 0, page)
				logenc.DeleteOldsFiles(pathdata+"/replace/"+fileN, "")
				hashSumFile = logenc.FileMD5(fileName)
			} else if currentfile != fileN {

				break
			} else if countline >= 99 {
				//TransmitUlidPagination(conn, fileName)
				countline = 0
			} else if currentpage != page && currentfile == fileN {
				logs.InfoLogger.Println("Change page:" + string(rune(page)))

				logenc.DeleteOldsFiles(pathdata+"/replace/"+fileN, "")
				conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
				countline = TailingLogsInFileAll(0, fileName, conn, 0, page)
				logenc.DeleteOldsFiles(pathdata+"/replace/"+fileN, "")
				currentpage = page

			}

		}

		return

	}

	UlidC := bleveSI.ProcBleveSearchv2(fileN, lookFor)
	if len(UlidC) == 0 {
		println("Break")
		return
	} else {

		var countCheck int

		fmt.Println("countSearch", 0)
		for i := 0; i < len(UlidC); i++ {
			_, found := SearchMap[UlidC[i]]
			if found {
				countCheck++
			}
		}
		fmt.Println("...............countCheck", countCheck)
		CountPage := "<countpage>" + strconv.Itoa(countCheck) + "</countpage>"
		conn.WriteMessage(websocket.TextMessage, []byte(CountPage))
		countCheck = 0
		currentpage := 0
		fmt.Println("fileAdr", fileN, "lookFor", lookFor)
		fmt.Println("SearchMap: ", "UlidC: ", UlidC, "page: ", page, "conn", conn)
		tailLogsInFind(SearchMap, UlidC, page, conn)

		for {

			if Fname != fileName {
				break
			} else if currentpage != page && currentfile == fileN {
				tailLogsInFind(SearchMap, UlidC, page, conn)
				currentpage = page

			}
		}

		return

	}

}
func tailLogsInFind(SearchMap map[string]logenc.LogList, UlidC []string, page int, conn *websocket.Conn) {
	if page == 0 {
		page = 1
	}
	//conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
	fmt.Println(UlidC)
	fmt.Println(conn)
	for i := page*100 - 100; i < page*100; i++ {
		v, found := SearchMap[UlidC[i]]
		if found {
			fmt.Println(logenc.EncodeXML(v))
			conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(v)))
		}
	}

}

func TailingLogsInFileAll(countline int, fileName string, conn *websocket.Conn, current int64, page int) int {

	var statusPagination bool = false

	fileN := filepath.Base(fileName)
	original, err := os.Open(fileName)
	if err != nil {
		logs.ErrorLogger.Println("TailingLogsInFileAll " + err.Error())
		log.Println(err)
	}
	exec.Command("/bin/bash", "-c", "echo > "+pathdata+"/replace/"+fileN).Run()
	logenc.CopyFile(pathdata+"/replace/", fileN, original)
	//var countline int = 0
	taillog, err := tail.TailFile(pathdata+"/replace/"+fileN,
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
		return countline
	}
	//go taillog.StopAtEOF()
	//conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
	//go taillog.StopAtEOF()
	for line := range taillog.Lines {
		if page != 0 && line.Text != "" && line.Text != " " {
			//taillog.StopAtEOF()
			pagUlid := paginationUlids[page]
			csvsimpl := logenc.ProcLineDecodeXML(line.Text)
			currentUlid := csvsimpl.XML_RECORD_ROOT[0].XML_ULID
			if pagUlid == currentUlid {
				statusPagination = true
			}
			if statusPagination {
				countline++
				conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(csvsimpl)))
			}
		} else {
			csvsimpl := logenc.ProcLineDecodeXML(line.Text)
			countline++
			conn.WriteMessage(websocket.TextMessage, []byte(logenc.EncodeXML(csvsimpl)))
		}

		if countline == 100 {
			//taillog.Stop()
			return countline
		}

	}

	return countline

}

func followCodeStatus(conn *websocket.Conn) {
	//Reset(conn)
	if PrevNetConn != conn.UnderlyingConn() {
		PrevNetConn = conn.UnderlyingConn()
	}
	for {
		if PrevNetConn != conn.UnderlyingConn() {
			break //выход из цикла при создании другого канала соединения
		}
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			logs.ErrorLogger.Println("Read mes foolow Code Status " + err.Error())
			log.Println(err, "followCodeStatus")
			return
		}
		fmt.Println("msgType", msgType)
		page, err = strconv.Atoi(string(msg[:]))
		if err != nil {
			logs.ErrorLogger.Println("byte to int " + err.Error())
			currentfile = string(msg)
		}
		fmt.Println("Page", page)
	}

}

func UlidPaginationFile(conn *websocket.Conn, fileName string) {
	var CountPage string
	paginationUlids = make(map[int]string)
	var (
		strSlice []string

		countline int
		page      int    = 0
		firstUlid string = " "
	)
	taillog, err := tail.TailFile(fileName,
		tail.Config{
			Follow: false,
			Location: &tail.SeekInfo{
				Whence: io.SeekStart, //!!!

			},
		})
	if err != nil {
		logs.WarningLogger.Println("Occurred in opening the file Pagination file: " + err.Error())
		fmt.Fprintln(os.Stderr, "Error occurred in opening the file: ", err)
		return
	}
	for line := range taillog.Lines {
		strSlice = append(strSlice, logenc.ProcLineDecodeXMLUlid(line.Text))
		countline++
		if countline == 100 {
			page++
			countline = 0
			firstUlid = strSlice[0]
			paginationUlids[page] = firstUlid
			strSlice = nil

		}
		go taillog.StopAtEOF() //end tail and stop service
	}
	page++
	CountPage = "<countpage>" + strconv.Itoa(page) + "</countpage>"
	conn.WriteMessage(websocket.TextMessage, []byte(CountPage))
	firstUlid = strSlice[1]
	paginationUlids[page] = firstUlid

	countline = 0
	strSlice = nil

}

func UlidPaginationDir(conn *websocket.Conn, countFiles int, fileList map[string][]string) {
	//var CountPage string
	var CountPage string
	paginationUlids = make(map[int]string)
	FileUlids = make(map[string]string)
	var (
		strSlice []string

		countline int
		page      int = 1
		//fileName  string
		firstUlid string = " "
	)

	//fileList["FileList"] = util.Conf.Dir

	for i := 0; i < countFiles; i++ {
		fileName := fileList["FileList"][i]
		taillog, err := tail.TailFile(fileName,
			tail.Config{
				Follow: false,
				Location: &tail.SeekInfo{
					Whence: io.SeekStart, //!!!

				},
			})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error occurred in opening the file: ", err)
			logs.WarningLogger.Println("Occurred in opening the file Pagination Dir: " + err.Error())
			return
		}
		//go taillog.StopAtEOF()
		go taillog.StopAtEOF() //end tail and stop service
		for line := range taillog.Lines {
			strSlice = append(strSlice, logenc.ProcLineDecodeXMLUlid(line.Text))
			countline++
			if countline == 100 {
				page++
				countline = 0
				firstUlid = strSlice[1]
				paginationUlids[page] = firstUlid
				strSlice = nil
			}

		}
		if countline != 0 && countline < 100 && page == 0 {
			page++
			firstUlid = strSlice[1]
			paginationUlids[page] = firstUlid
		}
	}
	CountPage = "<countpage>" + strconv.Itoa(page) + "</countpage>"
	conn.WriteMessage(websocket.TextMessage, []byte(CountPage))
	//	firstUlid = strSlice[1]
	paginationUlids[page] = firstUlid
	countline = 0
	strSlice = nil
}

var (
	indexingfilet = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "logi2_indexing_file_count",
		Help: "How many files indexing to view on page",
	})
)

// IndexFiles - takes argument as a list of files and directories and returns
// a list of absolute file strings to be tailed
func IndexFiles(fileList []string) error {
	// Re-initialize the visited array
	visited = make(map[string]bool)

	// marking all entries possible stale
	// will be removed if not updated
	for k := range indexMap {
		indexMap[k] = false
	}

	for _, file := range fileList {
		dfs(file)
	}
	// Re-initialize the file list array
	FileList = make([]string, 0)

	// Iterate through the map that contains the filenames
	for k, v := range indexMap {
		if !v {
			delete(indexMap, k)
			continue
		}
		//fmt.Fprintln(os.Stderr, k)
		FileList = append(FileList, k)
	}
	//filepath.Base
	//filename
	for _, f := range FileList {
		fileN := filepath.Base(f)
		FileName = append(FileName, fileN)
	}
	Conf.Dir = FileList
	go func() {
		for _, file := range FileList {
			fileN := filepath.Base(file)
			//logs.InfoLogger.Println("Indexing by bleve" + fileN)
			bleveSI.ProcBleve(fileN, file)
		}
	}()
	Conf.Dir1 = FileName
	FileName = nil
	//metrics

	indexingfilet.Set(float64(len(Conf.Dir)))
	fmt.Fprintln(os.Stderr, "Indexing complete !, file index length: ", len(Conf.Dir))
	return nil
}

func dfs(file string) {
	// Mostly useful for first entry, as the paths may be like ../dir or ~/path/../dir
	// or some wierd *nixy style, Once the file is cleaned and made into an absolute
	// path, it should be safe as the next call is basepath(abspath) + "/" + name of
	// the file which should be accurate in all terms and absolute without any
	// funky conversions used in OS
	file = filepath.Clean(file)
	absPath, err := filepath.Abs(file)

	if err != nil {
		logs.FatalLogger.Println("Unable to get absolute path of the file: " + err.Error())
		fmt.Fprintf(os.Stderr, "Unable to get absolute path of the file %s; err: %s\n", file, err)
	}
	if _, ok := visited[file]; ok {
		// if the absolute path has been visited, return without processing
		return
	}
	visited[file] = true
	s, err := os.Stat(absPath)
	if err != nil {
		logs.WarningLogger.Println("Unable to stat file: " + file + ":" + err.Error())
		fmt.Fprintf(os.Stderr, "Unable to stat file %s; err: %s\n", file, err)
		return
	}
	// check if the file is a directory
	if s.IsDir() {
		basepath := filepath.Clean(file)
		filelist, _ := ioutil.ReadDir(absPath)
		for _, f := range filelist {
			dfs(basepath + string(os.PathSeparator) + f.Name())
			//dfs(basepath + string(os.PathSeparator) + f.Name())
		}
	} else if strings.ContainsAny(s.Mode().String(), "alTLDpSugct") {
		// skip these files
		// try including names PIPES
	} else {
		// only remaining file are ascii files that can be then differentiated
		// by the user as golang has only these many categorization
		// Note : this appends the absolute paths
		// Insert the absPath into the Map, avoids duplicates in successive cron runs
		indexMap[absPath] = true
	}
}

func TailDir2(conn *websocket.Conn, lookFor string, SearchMap map[string]logenc.LogList, startUnixTime int64, endUnixTime int64, countFiles int, fileList map[string][]string) {

	if (lookFor == "" || lookFor == " " || lookFor == "Search") && (startUnixTime == 0 || endUnixTime == 0) {
		go followCodeStatus(conn)
		var countline int = 0
		//var currentpage int = 0
		fileaddr := fileList["FileList"][0]
		fileN := filepath.Base(fileaddr)
		logenc.DeleteOldsFiles(pathdata+"/replace/"+fileN, "")
		conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
		countline = TailingLogsInFileAll(countline, fileaddr, conn, 0, page)
		logenc.DeleteOldsFiles(pathdata+"/replace/"+fileN, "")
		currentfile := 0
		currentpage := page
		PrevNetConn = conn.UnderlyingConn()
		for {
			if PrevNetConn != conn.UnderlyingConn() {
				break //выход из цикла при создании другого канала соединения
			}
			if countline < 99 && currentfile < countFiles {
				for i := 1; i < countFiles; i++ {
					currentfile++
					//Очистка поля
					//conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
					fileaddr := fileList["FileList"][i]
					//Indexing(fileaddr)
					fileN := filepath.Base(fileaddr)
					go logenc.Replication(fileaddr)
					bleveSI.ProcBleve(fileN, fileaddr)
					countline += TailingLogsInFileAll(countline, fileaddr, conn, 0, 0)
					logenc.DeleteOldsFiles(pathdata+"/replace/"+fileN, "")
					if countline > 99 {
						countline = 0
						break
					}
				}
			} else if currentpage != page {
				logs.InfoLogger.Println("Change page:" + string(rune(page)))

				for i := 0; i < countFiles; i++ {
					fileAdr := fileList["FileList"][i]
					fileN := filepath.Base(fileAdr)
					v, found := paginationUlids[page]
					if found == true {
						fmt.Println("Ulid", v)
						UlidC := bleveSI.ProcBleveSearchv2(fileN, v)
						fmt.Println("len Ulidc", len(UlidC))
						if len(UlidC) != 0 {
							countline = 0
							logenc.DeleteOldsFiles(pathdata+"/replace/"+fileN, "")
							conn.WriteMessage(websocket.TextMessage, []byte("<start></start>"))
							countline = TailingLogsInFileAll(countline, fileaddr, conn, 0, page)
							logenc.DeleteOldsFiles(pathdata+"/replace/"+fileN, "")
							currentpage = page

						}
					}
				}
			}
		}
	} else {
		for i := 0; i < countFiles; i++ {
			fileAdr := filepath.Base(fileList["FileList"][i])
			//fileN := filepath.Base(fileAdr)
			UlidC := bleveSI.ProcBleveSearchv2(fileAdr, lookFor)
			if len(UlidC) != 0 {
				fmt.Println("fileAdr", fileAdr, "lookFor", lookFor)
				fmt.Println("SearchMap: ", "UlidC: ", UlidC, "page: ", page, "conn")

				tailLogsInFind(SearchMap, UlidC, 0, conn)
			}

		}
		var countCheck int = 0
		for i := 0; i < countFiles; i++ {
			fileAdr := fileList["FileList"][i]
			fileN := filepath.Base(fileAdr)
			UlidC := bleveSI.ProcBleveSearchv2(fileN, lookFor)
			if len(UlidC) != 0 {
				//tailLogsInFind(SearchMap, UlidC, 1, conn)
				for i := 0; i < len(UlidC); i++ {
					_, found := SearchMap[UlidC[i]]
					if found {
						countCheck++
					}
				}
			}
		}
		//tailLogsInFind(SearchMap, fileList["FileList"][0], 1, conn)
		CountPage := "<countpage>" + strconv.Itoa(countCheck) + "</countpage>"
		conn.WriteMessage(websocket.TextMessage, []byte(CountPage))
		countCheck = 0

		currentpage := 0
		page = 0

		PrevNetConn = conn.UnderlyingConn()
		for {
			if PrevNetConn != conn.UnderlyingConn() {
				break //выход из цикла при создании другого канала соединения
			}
			if currentpage != page {
				for i := 0; i < countFiles; i++ {
					fileAdr := fileList["FileList"][i]
					fileN := filepath.Base(fileAdr)
					UlidC := bleveSI.ProcBleveSearchv2(fileN, lookFor)
					if len(UlidC) != 0 {
						fmt.Println("SearchMap: ", SearchMap, ": ", "page: ", page)

						tailLogsInFind(SearchMap, UlidC, page, conn)
						currentpage = page
					}

				}
			}
		}
	}
}

var (
	countFilesFromLink = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logi2_upload_files_vfs_count",
		Help: "How many files from vfs",
	})
)

func GetFiles(address string, port string) error {
	//var signature bool = false
	resp, err := http.Get("http://" + address + ":" + port + "/vfs/data/")
	if err != nil {
		logs.ErrorLogger.Println("Did not get files: " + err.Error())
		return err
		//log.Fatal(err)

	}
	countFilesFromLink.Inc()
	for _, v := range logenc.GetLinks(resp.Body) {

		fmt.Println(address)

		fullURLFile := "http://" + address + ":" + port + "/vfs/data/" + v

		fileURL, err := url.Parse(fullURLFile)
		if err != nil {
			logs.FatalLogger.Println("Parse: " + err.Error())
			log.Fatal("Parse", err)
		}
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName := segments[len(segments)-1]

		func() { // lambda for defer file.Close()
			file, err := os.OpenFile(pathdata+"/testsave/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				logs.FatalLogger.Println("Getfiles: " + err.Error())
				log.Fatal("Getfiles", err)
				//file.Close()
				//return
			}

			defer file.Close()

			client := http.Client{
				CheckRedirect: func(r *http.Request, _ []*http.Request) error {
					r.URL.Opaque = r.URL.Path
					return nil
				},
			}
			// Put content on file
			resp, err := client.Get(fullURLFile)
			if err != nil {
				logs.InfoLogger.Println("DeleteOldsFiles: " + err.Error())
				logenc.DeleteOldsFiles(pathdata+"/testsave/"+fileName, "")
				return
				//log.Fatal(err)
			}
			defer resp.Body.Close()
			contain := strings.Contains(fileName, "md5")
			if contain && logenc.CheckFileSum(pathdata+"/testsave/"+fileName, "rep") {
				signature = true

				fileS, err := os.OpenFile("/var/log/logi2/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					logs.InfoLogger.Println("DeleteOldsFiles: " + err.Error())
					log.Println("Open for copy", err)
				}
				defer fileS.Close()
				_, err = io.Copy(fileS, resp.Body)
				if err != nil {

					log.Println("Copy", err)
				}
				logenc.WriteFileSum(pathdata+"/testsave/"+fileName, "rep")
				log.Println("*1")
				logenc.DeleteOldsFiles(pathdata+"/testsave/"+fileName, "")

			} else if !contain {
				_, err = io.Copy(file, resp.Body)
				if err != nil {
					logs.InfoLogger.Println("Copy: " + err.Error())
					log.Println("Copy", err)
				}
			}

			if signature && !contain {
				last3 := fileName[len(fileName)-3:]
				if logenc.CheckFileSum(pathdata+"/testsave/"+fileName, last3) {
					log.Println("*2")
					logenc.DeleteOldsFiles(pathdata+"/repdata/"+fileName, "")
					logenc.Replication(pathdata + "/testsave/" + fileName)
					logenc.WriteFileSum(pathdata+"/testsave/"+fileName, "rep")
					fmt.Println("Merge", fileName)

					log.Println("*3")
					logenc.DeleteOldsFiles(pathdata+"/testsave/"+fileName, "")

				} else {
					logenc.Replication(pathdata + "/testsave/" + fileName)
					logenc.WriteFileSum(pathdata+"/testsave/"+fileName, "rep")
					fmt.Println("Merge", fileName)
					log.Println("*4")
					logenc.DeleteOldsFiles(pathdata+"/testsave/"+fileName, "")
				}

			} else if !signature && !contain {
				//time.Sleep(15 * time.Second)
				logenc.Replication(pathdata + "/testsave/" + fileName)
				logenc.WriteFileSum(pathdata+"/testsave/"+fileName, "rep")
				fmt.Println("Merge", fileName)
				log.Println("*5")
				logenc.DeleteOldsFiles(pathdata+"/testsave/"+fileName, "")
			}

		}()
	}
	return nil
}

//Disk Check
type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
	//limit big.Float = 0.8
)

var (
	count_disk_usage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "logi2_disk_usage_gigabytes",
		Help: "Disk usage gigabytes ",
	})
	count_disk_memory = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "logi2_disk_gigabytes",
		Help: "Disk available fo usage",
	})
	count_disk_memory_free = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "logi2_disk_free_gigabytes",
		Help: "Disk free",
	})
)

func FindOldestfile(dir string) {
	var name string
	var cutoff = time.Hour
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		logs.FatalLogger.Println("Find old files: " + err.Error())
		log.Fatal("FindOldestfile", err.Error())
	}
	now := time.Now()
	for _, info := range fileInfo {
		if diff := now.Sub(info.ModTime()); diff > cutoff {
			cutoff = now.Sub(info.ModTime())
			name = info.Name()

		}
	}
	logenc.DeleteOldsFiles(dir+name, "")
}

func DeleteFile90(dir string) {

	var cutoff = 24 * time.Hour * 90
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		logs.FatalLogger.Println("Delete file 90: " + err.Error())
		log.Fatal("DeleteFile90", err.Error())
	}
	now := time.Now()
	for _, info := range fileInfo {
		if diff := now.Sub(info.ModTime()); diff > cutoff {
			fmt.Printf("Deleting %s which is %s old\n", info.Name(), diff)
			logenc.DeleteOldsFiles(dir+info.Name(), "")

		}
	}

}

func CheckIPAddress(ip string) bool {
	/* if ip == "localhost" {
		fmt.Printf("IP Address: %s - Valid\n", ip)
		return true
	} else  */if net.ParseIP(ip) == nil {
		fmt.Printf("IP Address: %s - Invalid\n", ip)
		return false
	} else {
		fmt.Printf("IP Address: %s - Valid\n", ip)
		return true
	}

}
