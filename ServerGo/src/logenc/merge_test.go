package logenc

import (
	"fmt"
	"testing"
	"time"
)

func TestMergeLines(t *testing.T) {
	//var wg sync.WaitGroup

	ch1 := make(chan LogList, 10)
	ch2 := make(chan LogList, 10)

	now := time.Now().UTC()
	go func() {
		gen := func() LogList {
			var listlog LogList
			listlog.XML_RECORD_ROOT = make([]Log, 1)
			listlog.XML_RECORD_ROOT[0].GenTestULID(now)
			now = now.Add(time.Millisecond)
			return listlog
		}

		ch1 <- gen()
		ch1 <- gen()
		ch2 <- gen()

		dup := gen()

		ch2 <- dup
		ch1 <- dup

		ch1 <- gen()

		close(ch1)
		close(ch2)
	}()
	want := 5
	ch3 := MergeLines(ch1, ch2)
	got := 0
	for val := range ch3 {

		if len(val.XML_RECORD_ROOT) != 0 {
			got++
			fmt.Println(val.XML_RECORD_ROOT[0].XML_ULID)
		}
	}
	if got != 5 {
		t.Errorf("MergeLines() = %v, want %v", got, want)
	}

}

func TestReplication(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			Replication("/home/nik/projects/Course/logi2/service/genrlogs./gen_logs_coded193")
		})
	}
}

func ATestMerge(t *testing.T) {

	Merge("./testmerge/", "/home/nik/projects/Course/logi2/service/genrlogs./gen_logs_coded193")

}
