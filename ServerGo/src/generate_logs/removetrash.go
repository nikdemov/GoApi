package generate_logs

import (
	"io/ioutil"
	"strings"

	"nikworkedprofile/GoApi/ServerGo/src/logenc"
)

//Remove file in dir
func RemoveByConfig() {
	//Remove(pathdata+"/genrlogs./", "gen_logs_coded")
	Remove(pathdata+"/repdata/", "gen_logs_coded")
}

func Remove(dirpath string, lineS string) {
	//var count int = 0
	files, _ := ioutil.ReadDir(dirpath)

	for _, file := range files {
		//go R(count)
		fileN := file.Name()
		contain := strings.Contains(fileN, lineS)
		if contain {
			logenc.DeleteOldsFiles(dirpath+fileN, "")
		}

	}

}
