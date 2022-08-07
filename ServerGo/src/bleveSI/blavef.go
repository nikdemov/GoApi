package bleveSI

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"nikworkedprofile/GoApi/ServerGo/src/logenc"
	logs "nikworkedprofile/GoApi/ServerGo/src/logs_app"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/oklog/ulid/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Logger   *log.Logger
	pathlogs = "/var/local/logi2"
	pathdata = "/var/log/logi2"
)

func bleveIndex(fileN string) (bleve.Index, error) {

	dir := pathdata + "/blevestorage/"
	extension := ".bleve"
	metaname := dir + fileN + extension
	index, err := bleve.Open(metaname)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		fmt.Println(err)

		index, err = bleve.NewUsing(metaname, mapping, scorch.Name, scorch.Name, nil)
		fmt.Println(err)

	}

	return index, err
}

var (
	count_indexing_mes = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logi2_indexing_strings_count",
		Help: "How many messages indexing by bleve",
	})
)

func ProcBleve(fileN string, file string) {
	var count int = 0
	if !logenc.CheckFileSum(file, "") {
		return
	}

	var wg sync.WaitGroup
	index, err := bleveIndex(fileN)
	if err != nil {
		logs.ErrorLogger.Println("Bleve Index" + err.Error())
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
			batch := index.NewBatch()
		brloop:
			for {

				select {
				case line, ok := <-ch:
					if !ok {
						break brloop
					}
					if count == 100 {
						err = index.Batch(batch)
						if err != nil {
							logs.WarningLogger.Println("Bleve Index Batch" + err.Error())
						}
						count = 0
						batch = index.NewBatch()
					}
					if len(line) != 0 {

						data = logenc.ProcLineDecodeXML(line)

						if len(data.XML_RECORD_ROOT) > 0 {
							err := batch.Index(data.XML_RECORD_ROOT[0].XML_ULID, data)
							if err != nil {
								logs.ErrorLogger.Println("Batch Bleve" + err.Error())
								close(ch)
								index.Close()
								return
							}
							count_indexing_mes.Inc()
							count++
						}
					}

				}
			}
			err = index.Batch(batch)
			if err != nil {
				logs.WarningLogger.Println("Bleve index batch" + err.Error())
			}

		}()
	}
	err = logenc.ReadLines(file, func(line string) {
		ch <- line
	})
	if err != nil {
		logs.ErrorLogger.Println("Http serve error" + err.Error())
		close(ch)
		index.Close()
		return
	}
	close(ch)
	wg.Wait()
	index.Close()
	logenc.WriteFileSum(file, "")
}

func ProcBleveSearchv2(fileN string, word string) []string {
	//var query *query.MatchQuery
	dir := pathdata + "/blevestorage/"
	extension := ".bleve"
	filename := fileN
	metaname := dir + filename + extension
	index, _ := bleve.OpenUsing(metaname, nil)

	query := bleve.NewMatchQuery(word)
	mq := bleve.NewMatchPhraseQuery(word)
	rq := bleve.NewRegexpQuery(word)
	q := bleve.NewDisjunctionQuery(query, mq, rq)
	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Size = 1000000000000000000
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		logs.WarningLogger.Println("Bleve Search" + err.Error())
		return []string{" "}
	}
	searchRequest.Fields = []string{"XML_RECORD_ROOT"}
	docs := make([]string, 0)
	for _, val := range searchResult.Hits {
		id := val.ID
		docs = append(docs, id)
	}
	//sort
	for i := len(docs); i > 0; i-- {
		for j := 1; j < i; j++ {
			j2, _ := ulid.Parse(docs[j-1])
			j1, _ := ulid.Parse(docs[j])
			if j2.Compare(j1) == 1 {
				intermediate := docs[j]
				docs[j] = docs[j-1]
				docs[j-1] = intermediate
			}

		}
	}

	index.Close()
	return docs

}
