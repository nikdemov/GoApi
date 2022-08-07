package test

import (
	"testing"

	"nikworkedprofile/GoApi/src/logenc"
)

func BenchmarkCheckFileSum(b *testing.B) {
	type args struct {
		file string
	}

	logenc.CheckFileSum("./view/22-06-2021", "")
	//logenc.CheckFileSum("/home/nik/projects/Course/logi2/logtest/gen_logs1")

}
