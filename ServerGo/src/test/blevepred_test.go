package test

import (
	"testing"

	"nikworkedprofile/GoApi/src/bleveSI"
)

func BenchmarkProcFileBleveS(b *testing.B) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcFileBreveSLOWLY("test777", "./view/22-06-2021")
}
