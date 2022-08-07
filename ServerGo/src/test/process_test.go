package test

import (
	"testing"

	"nikworkedprofile/GoApi/ServerGo/src/logenc"
)

func BenchmarkProcMapFile(b *testing.B) {
	type args struct {
		file string
	}

	logenc.ProcMapFile("./view/22-06-2021")
	//t.StartTimer()
}

func BenchmarkProcMapFilePP(b *testing.B) {
	type args struct {
		file string
	}

	logenc.ProcMapFileREZERV("./view/22-06-2021")

}
