//go:build dev
// +build dev

package test

import (
	"testing"

	"nikworkedprofile/GoApi/src/bleveSI"
)

func TestProcBleveScorch(t *testing.T) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleveScorch("test5", "./view/22-06-2021")
}

func testBleveSearch(t *testing.T) {
	type args struct {
		fileN string
		file  string
	}

	bleveSI.ProcBleveSearchv2("test4", "1")

}
