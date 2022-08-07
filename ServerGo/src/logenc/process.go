package logenc

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	logs "nikworkedprofile/GoApi/ServerGo/src/logs_app"

	"golang.org/x/net/html"
)

var (
	pathlogs = "/var/log/logi2"
	pathdata = "/var/local/logi2"
)
var (
	Logger *log.Logger
	mu     sync.Mutex
	ind    bool
	//fileSize int64

	//true untyped bool = true
)

/* var (
	sliceLoglist []LogList
) */

func ProcLine(line string) (csvF string) {

	if len(line) == 0 {

		return
	}
	lookFor := "<loglist>"
	xmlline := DecodeLine(line)
	contain := strings.Contains(xmlline, lookFor)
	if !contain {

		return xmlline
	}

	val, err := DecodeXML(xmlline)
	if err != nil {
		logs.WarningLogger.Println("Procline" + err.Error())

		return
	}

	csvline := EncodeCSV(val)
	//fmt.Print(csvline)
	return csvline
}

func ProcLineCSVstr(line string) (csvF string) {

	if len(line) == 0 {

		return
	}
	lookFor := "<loglist>"
	xmlline := DecodeLine(line)
	contain := strings.Contains(xmlline, lookFor)
	if !contain {

		return xmlline
	}

	//fmt.Print(csvline)
	return xmlline
}

/* func ProcLineCSVLoglost(line string) (val LogList) {

	if len(line) == 0 {

		return
	}
	lookFor := "<loglist>"
	xmlline := DecodeLine(line)
	contain := strings.Contains(xmlline, lookFor)
	if !contain {

		return
	}
	val, err := DecodeXML(xmlline)
	if err != nil {

		return
	}
	//fmt.Print(csvline)
	return val
} */

func procLineq(line string) (csvF string) {

	if len(line) == 0 {

		return
	}

	xmlline := DecodeLine(line)
	val, err := DecodeXML(xmlline)
	if err != nil {
		logs.WarningLogger.Println("procline" + err.Error())

		return
	}

	csvline := EncodeCSV(val)
	return csvline
}

func ProcFile(file string) {
	ch := make(chan string, 100)
	//log.Println("1")
	for i := runtime.NumCPU() + 1; i > 0; i-- {
		go func() {
			for {
				line := <-ch

				ProcLine(line)
			}

		}()
	}

	err := ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		logs.WarningLogger.Println("ProcFile Readlines" + err.Error())

		fmt.Println("ReadLines: ", err)
		close(ch)
		return
		//log.Fatalf("ReadLines: %s", err)
	}
	close(ch)
}

func ProcDir(dir string) {

	filepath.Walk(dir,
		func(path string, file os.FileInfo, err error) error {
			if err != nil {
				logs.ErrorLogger.Println("ProcDir" + err.Error())

				return err
			}
			if !file.IsDir() {

				ProcFile(path)
			}
			return nil
		})
}

//Write Decode logs (it is not working(may be))
func ProcWrite(dir string) {

	filepath.Walk(dir,
		func(path string, file os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !file.IsDir() {

				procFileWrite(path)
			}
			return nil
		})
}

func procFileWrite(file string) {
	CreateDir("./writedeclog")
	fileN := filepath.Base(file)

	filew, err1 := os.OpenFile(pathdata+"/writedeclog/"+fileN+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err1 != nil {
		log.Fatal("procFileWrite error", err1)
	}

	Logger = log.New(filew, "", 0)

	ch := make(chan string, 100)

	for i := runtime.NumCPU() + 1; i > 0; i-- {
		go func() {
			for {
				line := <-ch

				Logger.Println(procLineq(line))

			}

		}()
	}

	err := ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		fmt.Println("ReadLines: ", err)
		close(ch)
		return
	}
	close(ch)
}

///
/* func Promrun(port string) {
	//portStr := strconv.Itoa(port)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+port+"", nil))
}
*/
func ProcLineDecodeXML(line string) (val LogList) {

	if len(line) == 0 {

		return
	}
	lookFor := "<loglist>"
	xmlline := DecodeLine(line)
	contain := strings.Contains(xmlline, lookFor)
	if !contain {

		return
	}
	val, err := DecodeXML(xmlline)
	if err != nil {
		logs.ErrorLogger.Println("DecodeXML" + err.Error())

		return
	}
	return val
}

func ProcLineDecodeXMLUlid(line string) (ulid string) {

	if len(line) == 0 {

		return
	}
	lookFor := "<loglist>"
	xmlline := DecodeLine(line)
	contain := strings.Contains(xmlline, lookFor)
	if !contain {

		return
	}
	val, err := DecodeXML(xmlline)
	if err != nil {
		logs.ErrorLogger.Println("ProcLineDecodeXMLUlid" + err.Error())
		return
	}
	ulid = val.XML_RECORD_ROOT[0].XML_ULID
	return ulid
}

func ProcLineDecodeXMLType(line string) (typem string) {

	if len(line) == 0 {

		return
	}
	lookFor := "<loglist>"
	xmlline := DecodeLine(line)
	contain := strings.Contains(xmlline, lookFor)
	if !contain {

		return
	}
	val, err := DecodeXML(xmlline)
	if err != nil {
		logs.ErrorLogger.Println("ProcLineDecodeXMLType" + err.Error())
		return
	}
	typem = val.XML_RECORD_ROOT[0].XML_TYPE
	return typem
}

func ProcMapFile(file string) map[string]LogList {
	//func ProcMapFile(file string) {
	if len(file) <= 0 {
		log.Println("ProcMapFile file = 0")

		return nil
	}
	ch := make(chan string, 100)
	SearchMap := make(map[string]LogList)
	var wg sync.WaitGroup
	var data LogList
	go func() {
		for {
			line, ok := <-ch
			if !ok {
				break
			}
			wg.Add(1)
			go func(line string) {
				//wg.Add(1)
				defer wg.Done()
				if len(line) != 0 {
					data = ProcLineDecodeXML(line)
					if len(data.XML_RECORD_ROOT) > 0 {
						mu.Lock()
						SearchMap[data.XML_RECORD_ROOT[0].XML_ULID] = data
						mu.Unlock()
					}
				}
				//data = ProcLineDecodeXML(line)
				//datas = ProcLineCSVLoglost(line)\
				//mu.Lock()
				/* if len(data.XML_RECORD_ROOT) > 0 {
					mu.Lock()
					SearchMap[data.XML_RECORD_ROOT[0].XML_ULID] = data
					mu.Unlock()
				} */
				//mu.Unlock()
				//defer wg.Done()
			}(line)

		}
	}()
	err := ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		logs.ErrorLogger.Println("ProcMapFile:" + err.Error())

		fmt.Println("ReadLines: ", err)
		close(ch)
		return SearchMap
	}

	close(ch)
	wg.Wait()
	return SearchMap
}

//slowely not used
func ProcMapFileREZERV(file string) {
	if len(file) <= 0 {
		return
	}
	ch := make(chan string, 1000000)
	SearchMap := make(map[string]string)
	var data LogList
	var datas string
	err := ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		fmt.Println("ReadLines: ", err)
		close(ch)
		return
	}
	fmt.Println("run")
	for {

		line, ok := <-ch
		if !ok {
			break
		}
		data = ProcLineDecodeXML(line)
		datas = ProcLine(line)
		if len(data.XML_RECORD_ROOT) > 0 {
			SearchMap[data.XML_RECORD_ROOT[0].XML_ULID] = datas
		}

		if len(ch) == 0 {
			break
		}

	}
	close(ch)
}

func CheckFileSum(file string, typeS string) bool {
	ind = true
	checksum2 := FileMD5(file)
	fileN := filepath.Base(file)
	hashFileName := "md5" + typeS
	f, err := os.OpenFile(pathdata+"/"+hashFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checke(err)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), (checksum2 + " " + fileN)) {
			ind = false
			return ind
		} else {

			ind = true
			//WriteFileSum(file, typeS, path)
		}
		line++
	}
	scanner = nil
	return ind
}

func WriteFileSum(file string, typeS string) {

	checksum2 := FileMD5(file)

	fileN := filepath.Base(file)
	hashFileName := "md5" + typeS
	f, _ := os.OpenFile(pathdata+"/"+hashFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {

		if strings.Contains(scanner.Text(), (checksum2 + " " + fileN)) {

			ind = false
			return
		} else {
			ind = true
		}

		line++
	}
	scanner = nil
	//if ind == true && strings.Contains(scanner.Text(), (fileN)) {

	//	} else
	if ind {

		f.Write([]byte(checksum2 + " " + fileN + "\n"))
	}
	fi, _ := f.Stat()

	if fi.Size() == 0 {
		f.Write([]byte(checksum2 + " " + fileN + "\n"))
	}

}

func checke(e error) {
	if e != nil {
		panic(e)
	}
}

// FileMD5 создает md5-хеш из содержимого нашего файла.
func FileMD5(path string) string {
	h := md5.New()
	f, err := os.Open(path)
	//defer f.Close()
	if err != nil {
		logs.ErrorLogger.Println("Open Md file" + err.Error())
		f.Close()
		return "null"

	}
	defer f.Close()
	_, err = io.Copy(h, f)
	if err != nil {
		logs.FatalLogger.Println("Copy in md file" + err.Error())

		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

/*
func DeleteHTMLTeg(s string) (clean string) {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					clean = a.Val
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {

			f(c)
		}
	}
	f(doc)
	return clean
}
*/
//Collect all links from response body and return it as an array of strings
func GetLinks(body io.Reader) []string {
	var links []string
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			//: links list shoudn't contain duplicates
			return links
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}

				}
			}

		}
	}
}

func SearchT(dir string) {
	var text string
	var limit int

	var MassStr []Data

	fmt.Print("Enter limit: ")
	fmt.Scanln(&limit)
	fmt.Print("Enter text: ")
	fmt.Scanln(&text)

	chRes := make(chan Data, 100)
	go func() {
		scan := &Scan{}
		scan.Find = dir
		scan.Text = text
		scan.ChRes = chRes
		scan.LimitResLines = limit
		scan.Search()
		close(scan.ChRes)
	}()

ext:
	for i := 0; i < limit; i++ {

		data, ok := <-chRes
		if !ok {
			break ext
		}
		MassStr = append(MassStr, data)

	}
	sort.Slice(MassStr, func(i, j int) (less bool) {
		return MassStr[i].ID < MassStr[j].ID
	})
	fmt.Printf("%+v\n", MassStr)

}

//Convert number to word
var NumberToWord = map[int]string{
	1:  "one",
	2:  "two",
	3:  "three",
	4:  "four",
	5:  "five",
	6:  "six",
	7:  "seven",
	8:  "eight",
	9:  "nine",
	10: "ten",
	11: "eleven",
	12: "twelve",
	13: "thirteen",
	14: "fourteen",
	15: "fifteen",
	16: "sixteen",
	17: "seventeen",
	18: "eighteen",
	19: "nineteen",
	20: "twenty",
	30: "thirty",
	40: "forty",
	50: "fifty",
	60: "sixty",
	70: "seventy",
	80: "eighty",
	90: "ninety",
}

func convert1to99(n int) (w string) {
	if n < 20 {
		w = NumberToWord[n]
		return
	}

	r := n % 10
	if r == 0 {
		w = NumberToWord[n]
	} else {
		w = NumberToWord[n-r] + "-" + NumberToWord[r]
	}
	return
}

func convert100to999(n int) (w string) {
	q := n / 100
	r := n % 100
	w = NumberToWord[q] + " " + "hundred"
	if r == 0 {
		return
	} else {
		w = w + " and " + convert1to99(r)
	}
	return
}

func Convert1to1000(n int) (w string) {
	if n > 1000 || n < 1 {
		panic("func Convert1to1000: n > 1000 or n < 1")
	}

	if n < 100 {
		w = convert1to99(n)
		return
	}
	if n == 1000 {
		w = "one thousand"
		return
	}
	w = convert100to999(n)
	return
}
