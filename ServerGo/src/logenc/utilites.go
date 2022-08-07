package logenc

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	logs "nikworkedprofile/GoApi/src/logs_app"

	"github.com/oklog/ulid/v2"
)

//XML_Structure
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

func (me *Log) GenTestULID(tt time.Time) {
	now := time.Now().UnixNano()
	entropy := rand.New(rand.NewSource(now))
	timestamp := ulid.Timestamp(tt)
	me.XML_ULID = ulid.MustNew(timestamp, entropy).String()
}

var (
	count = 0
)

const (
	XOR_KEY = 59
	//shortForm = "2006.01.02-15.04.05"
)

//Read lines
func ReadLines(path string, fn func(line string)) error {
	file, err := os.Open(path)
	if err != nil {
		logs.ErrorLogger.Println("Readlines :" + err.Error())

		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 4*1024*1024)

	//c := 0
	for scanner.Scan() {
		//c++
		//println(c)
		fn(scanner.Text())
	}
	return scanner.Err()
}

/* func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
} */

func DecodeLine(line string) string {
	data, err := base64.StdEncoding.DecodeString(line)

	if err != nil {
		logs.ErrorLogger.Println("Decode line" + err.Error())

		fmt.Println("error:", err)
		return ""
	}

	if len(data) <= 0 {
		return ""
	}

	k := 0
	for {
		//XOR with lines
		data[k] ^= XOR_KEY
		k++
		if k >= len(data) {
			break
		}
	}
	//print("start1")
	//print(string(data))
	//print("end1")
	return string(data)
}

func EncodeLine(line string) string {
	//data := base64.StdEncoding.Strict().EncodeToString([]byte(line))
	//result := []byte(data)
	//print("EncodeLine")
	if len(line) <= 0 {
		return ""
	}
	result := []byte(line)
	k := 0
	for {
		//XOR with lines
		result[k] ^= XOR_KEY
		k++
		if k >= len(result) {
			break
		}
	}
	//print(line)
	data := base64.StdEncoding.Strict().EncodeToString(result)

	return data
}

func DecodeXML(line string) (LogList, error) {
	//print("DecodeXML")

	var v = LogList{}

	err := xml.Unmarshal([]byte(line), &v)

	return v, err
}

func EncodeXML(tmp LogList) (v string) {
	//print("DecodeXML")

	//var v = LogList{}

	k, _ := xml.Marshal(tmp)
	v = string(k)
	return v
}

func datestr2time(str string) time.Time {
	// format example: 08092021224536920  from xml
	//"02012006150405.000"

	const shortForm = "02012006150405.000"

	str2 := string(str[0:14]) + "." + string(str[14:17])
	t, _ := time.Parse(shortForm, str2)
	return t
}

func EncodeCSV(val LogList) string {
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)
	for _, logstr := range val.XML_RECORD_ROOT {
		//TIME
		//print("EncodeCSV")
		t := datestr2time(logstr.XML_TIME)
		typeM := "INFO"
		switch logstr.XML_TYPE {
		case "1":
			typeM = "DEBUG"
		case "2":
			typeM = "WARNING"
		case "3":
			typeM = "ERROR"
		case "4":
			typeM = "FATAL"
		}
		//TYPE
		/* typeM := "INFO"
		if logstr.XML_TYPE == "1" {
			typeM = "DEBUG"
		} else if logstr.XML_TYPE == "2" {
			typeM = "WARNING"
		} else if logstr.XML_TYPE == "3" {
			typeM = "ERROR"
		} else if logstr.XML_TYPE == "4" {
			typeM = "FATAL"
		} */
		//id := fmt.Sprint(count)
		err := writer.Write([]string{typeM, logstr.XML_APPNAME, logstr.XML_APPPATH, logstr.XML_APPPID, logstr.XML_THREAD, t.Format(time.RubyDate), logstr.XML_ULID, logstr.XML_MESSAGE, logstr.XML_DETAILS, logstr.DT_FORMAT})
		count++
		if err != nil {
			logs.FatalLogger.Println("EncodeCSV" + err.Error())

			log.Fatalln("error writing record to csv:", err)
		}
	}

	writer.Flush()
	return buf.String()
}

/*
func DecodeString1(val LogList) string {
	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)
	for _, logstr := range val.XML_RECORD_ROOT {
		//TIME
		//print("EncodeCSV")
		t := datestr2time(logstr.XML_TIME)
		//TYPE
		typeM := "INFO"
		if logstr.XML_TYPE == "1" {
			typeM = "DEBUG"
		} else if logstr.XML_TYPE == "2" {
			typeM = "WARNING"
		} else if logstr.XML_TYPE == "3" {
			typeM = "ERROR"
		} else if logstr.XML_TYPE == "4" {
			typeM = "FATAL"
		}
		//id := fmt.Sprint(count)
		err := writer.Write([]string{typeM, logstr.XML_APPNAME, logstr.XML_APPPATH, logstr.XML_APPPID, logstr.XML_THREAD, t.Format(time.RubyDate), logstr.XML_ULID, logstr.XML_MESSAGE, logstr.XML_DETAILS, logstr.DT_FORMAT})
		count++
		if err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	writer.Flush()
	return buf.String()
} */

//rune ='symbol'
func Remove(s string, symbol rune) string {
	return strings.Map(
		func(r rune) rune {
			if r != symbol {
				return r
			}
			return -1
		},
		s,
	)
}

//ПОлучить расширение файла
func GetExtensionFromFile(filename string) string {
	return path.Ext(filename)
}
