package bleveSI

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"

	"nikworkedprofile/GoApi/src/logenc"

	"github.com/blevesearch/bleve"
)

func ProcFileBreve(fileN string, file string) {
	var wg sync.WaitGroup
	var data logenc.LogList
	if len(file) <= 0 {
		return
	}

	dir := pathdata + "/blevestorage/"
	extension := ".bleve"
	metaname := dir + fileN + extension
	if !logenc.CheckFileSum(file, "") {
		return
	}

	index, err := bleve.Open(metaname)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(metaname, mapping)
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	// search for some text
	ch := make(chan string, 100)

	for i := runtime.NumCPU() + 1; i > 0; i-- {
		wg.Add(1)
		go func() {
			//wg.Add(1)
			defer wg.Done()

		brloop:
			for {
				select {
				case line, ok := <-ch:
					if !ok {
						break brloop
					}
					data = logenc.ProcLineDecodeXML(line)
					if len(data.XML_RECORD_ROOT) > 0 {
						index.Index(data.XML_RECORD_ROOT[0].XML_ULID, data)
					}
				}
			}

		}()

	}

	err = logenc.ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		fmt.Println("ReadLines: ", err)
		close(ch)
		return
	}
	close(ch)
	wg.Wait()
	index.Close()
	logenc.WriteFileSum(file, "")
}

func ProcBleveSearchv1(fileN string, word string) []string {
	dir := pathdata + "/blevestorage/"
	extension := ".bleve"
	filename := fileN
	metaname := dir + filename + extension
	index, _ := bleve.Open(metaname)
	query := bleve.NewMatchQuery(word)
	query.Fuzziness = 1
	mq := bleve.NewMatchPhraseQuery(word)
	rq := bleve.NewRegexpQuery(word)
	q := bleve.NewDisjunctionQuery(query, mq, rq)

	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Size = 1000000000000000000

	searchResult, _ := index.Search(searchRequest)
	searchRequest.Fields = []string{"XML_RECORD_ROOT"}

	docs := make([]string, 0)

	for _, val := range searchResult.Hits {
		id := val.ID
		docs = append(docs, id)
	}

	index.Close()
	return docs

}

func ProcFileBreveSLOWLY(fileName string, file string) {
	const pieces int = 10

	var wg sync.WaitGroup
	var lines []string

	if len(file) <= 0 {
		return
	}
	dir := pathdata + "/blevestorage/"
	extension := ".bleve"
	metaname := dir + fileName + extension
	if !logenc.CheckFileSum(file, "") {
		return
	}

	//metaname := "example.bleve"
	index, err := bleve.Open(metaname)

	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(metaname, mapping)
	}

	defer index.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	fileF, _ := os.Open(file)
	scanner := bufio.NewScanner(fileF)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 4*1024*1024)
	for scanner.Scan() {

		lines = append(lines, scanner.Text())
	}

	if len(lines) == 0 {
		return
	}

	if err != nil {
		fmt.Println("ReadLines: ", err)
		//close(ch)
		return
	}

	var datas [pieces][]logenc.LogList

	curNum := 0

	for _, line := range lines {

		datas[curNum] = append(datas[curNum], logenc.ProcLineDecodeXML(line))

		curNum++
		if curNum == pieces {
			curNum = 0
		}

	}
	for _, data := range datas {

		wg.Add(1)
		go func(dataPiece []logenc.LogList) {

			for _, data := range dataPiece {
				if len(data.XML_RECORD_ROOT) > 0 {
					index.Index(data.XML_RECORD_ROOT[0].XML_ULID, data)
				}
			}
			wg.Done()
		}(data)

	}

	wg.Wait()
	logenc.WriteFileSum(file, "")

}

//example Speed

func ProcBleveScorch(fileN string, file string) {
	if !logenc.CheckFileSum(file, "") {
		return
	}
	var wg sync.WaitGroup
	index, err := bleveIndex(fileN)
	if err != nil {
		fmt.Println(err)
		return
	}
	var data logenc.LogList
	ch := make(chan string, 100)
	for i := runtime.NumCPU() + 1; i > 0; i-- {
		wg.Add(1)
		go func() {
			//wg.Add(1)
			defer wg.Done()

		brloop:
			for {
				select {
				case line, ok := <-ch:
					if !ok {
						break brloop
					}
					data = logenc.ProcLineDecodeXML(line)
					if len(data.XML_RECORD_ROOT) > 0 {
						index.Index(data.XML_RECORD_ROOT[0].XML_ULID, data)
					}
				}
			}

		}()

	}

	err = logenc.ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		fmt.Println("ReadLines: ", err)
		close(ch)
		return
	}
	close(ch)
	wg.Wait()
	index.Close()
	logenc.WriteFileSum(file, "")

}
