package logs

import (
	"fmt"
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	FatalLogger   *log.Logger
	pathlogs      = "/var/log/logi2"
	pathdata      = "/var/local/logi2"
)

func InitLog() {
	CreateDir(pathlogs)
	CreateDir(pathdata)
	CreateDir(pathdata + "/replace")
	CreateDir(pathdata + "/blevestorage")
	fileConf, err := os.OpenFile(pathdata+"/config.toml", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File created successfully")
	defer fileConf.Close()

	filemd5, err := os.Create("/var/local/logi2/md5")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File created successfully")
	defer filemd5.Close()
	file, err := os.OpenFile(pathlogs+"/logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		//log.Fatal(err)
		return
	}

	//if err := os.Truncate(pathlogs+"/logs.log", 0); err != nil {
	//	log.Printf("Failed to truncate: %v", err)
	//}

	//logs, _ := ioutil.TempFile("", "logs.log")

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	FatalLogger = log.New(file, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
	defer file.Close()

}
func CreateDir(dirpath string) {
	//fileN := filepath.Base(path)
	//Create a folder/directory at a full qualified path
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		//ErrorLogger.Println("Create Dir" + err.Error())
		log.Println("CreateDir:", err)
	}
}
