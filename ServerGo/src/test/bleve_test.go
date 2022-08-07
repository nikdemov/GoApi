package test

import (
	"testing"

	"nikworkedprofile/GoApi/src/bleveSI"
)

func BenchmarkProcFileBleve(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleve("test4", "./view/22-06-2021")

}

/*func BenchmarkProcFileSSS(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleveScorch("test123", "./view/22-06-2021")
	//generator.Remove("/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded193","test")

}
*/
func BenchmarkProcFileBatch(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleve("test87", "./view/22-06-2021")
	//generator.Remove("/home/nik/projects/Course/logi2/genrlogs./gen_logs_coded193","test")

}

/* func TestProcBleveS(t *testing.T) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleve("test7", "./view/22-06-2021")

}
*/
