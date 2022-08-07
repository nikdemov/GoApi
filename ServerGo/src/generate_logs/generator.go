package generate_logs

import (
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"time"

	"nikworkedprofile/GoApi/src/logenc"
	logs "nikworkedprofile/GoApi/src/logs_app"
	"nikworkedprofile/GoApi/src/web/util"

	"github.com/Pallinder/go-randomdata"
	"github.com/oklog/ulid/v2"
)

var (
	pathlogs = "/var/log/logi2"
	pathdata = "/var/local/logi2"
)

type LogList struct {
	XMLName         xml.Name `xml:"loglist"` //dont touch XMLName
	XML_RECORD_ROOT []Log    `xml:"log"`
}
type Log struct {
	XML_APPNAME string `xml:"module_name,attr"`
	XML_APPPATH string `xml:"app_path,attr"`
	XML_APPPID  string `xml:"app_pid,attr"`
	XML_THREAD  string `xml:"thread_id,attr"`
	XML_TIME    string `xml:"time,attr"`
	XML_ULID    string `xml:"ulid,attr"`
	XML_TYPE    string `xml:"type,attr"`
	XML_MESSAGE string `xml:"message,attr"`
	XML_DETAILS string `xml:"ext_message,attr"`
	DT_FORMAT   string `xml:"ddMMyyyyhhmmsszzz,omitempty"`
}

var (
	Logger *log.Logger
	label  string = "00"
	//labeld     string
	countFile int = 0
	//countFiled int = 0
)

const (
	XOR_KEY = 59
)

func StructFile(count string) string {
	elem := "\""
	r := rand.New(rand.NewSource(99))
	XML_DETAILS := "Context:  -- void ::AbstractMonitor::,"
	now := time.Now().UnixNano()
	entropy := rand.New(rand.NewSource(now))
	timestamp := ulid.Timestamp(time.Now())
	XML_APPNAME := strconv.Itoa(r.Intn(10)) + " TEST"
	XML_APPPATH := "/" + strconv.Itoa(r.Intn(10)) + "/TEST/TEST"
	XML_APPPID := string(util.GetOutboundIP()[len(util.GetOutboundIP())-3:]) + "0" + count //strconv.Itoa(r.Intn(1000)) + "" // "7481,"
	XML_THREAD := strconv.Itoa(r.Intn(10)) + ""                                            //"88,"
	XML_MESSAGE := "Состояние '" + randomdata.IpV4Address() + "Cервер КС_UDP/Пинг'"
	XML_TYPE := strconv.Itoa(rand.Intn(4-1) + 1)
	address := randomdata.ProvinceForCountry("GB")
	//rand.Intn(max - min) + min
	//"29 05 2021 00 01 47 040"
	time1 := strconv.Itoa(rand.Intn(28-1)+1) + strconv.Itoa(rand.Intn(12-1)+1) + strconv.Itoa(rand.Intn(2022-2018)+2018) + "000147040"
	time_ulid := ulid.MustNew(timestamp, entropy)
	ulid1 := time_ulid.String()
	LINE := "<loglist><log module_name=" + elem + XML_APPNAME + elem +
		" app_path=" + elem + XML_APPPATH + elem +
		" app_pid=" + elem + XML_APPPID + elem +
		" thread_id=" + elem + XML_THREAD + elem +
		" time=" + elem + time1 + elem +
		" ulid=" + elem + ulid1 + elem +
		" type=" + elem + XML_TYPE + elem +
		" message=" + elem + XML_MESSAGE + elem +
		" ext_message=" + elem + XML_DETAILS + address + elem + "/></loglist>"

	//rand.Seed(time.Now().UnixNano())
	return LINE
}

func ProcGenN(count int, FileSize int64) {
	//Example()
	fmt.Println(" start  : ")

	filesFrom := string(util.GetOutboundIP()[len(util.GetOutboundIP())-3:])
	//	last3  := string(s[len(s)-3:])
	logenc.CreateDir(pathdata + "/repdata/")
	fmt.Println(" start Gen file: ", filesFrom)

	for {

		LINE := StructFile(strconv.Itoa(countFile))

		rand.Seed(time.Now().UnixNano())

		file, err := os.OpenFile(pathdata+"/repdata/gen_logs_coded"+label+filesFrom, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			logs.FatalLogger.Println("Create new gen file" + err.Error())
			//server.FatalLogger.Println("BlaveSearch" + err.Error())
			log.Fatal(err)
		}

		InfoLogger := log.New(file, "", 0)

		fi, _ := file.Stat()

		if fi.Size() >= FileSize {
			fmt.Println("Gen file: ", countFile)
			countFile++
			//logenc.WriteFileSum("./repdata/gen_logs_coded"+label+filesFrom, filesFrom, "./repdata/")
			label = strconv.Itoa(countFile)

		}

		infof := func(info string) {
			InfoLogger.Output(2, logenc.EncodeLine(info))
		}

		infof(LINE)

		//time.Sleep(time.Nanosecond * 1000000)

		if countFile >= count {
			return
		}

	}

}
