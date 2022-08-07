package logenc

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	logs "nikworkedprofile/GoApi/src/logs_app"
)

/* var (
	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "logi2_counter",
			Help: "search lines counter",
		})

	gaugeRT = prometheus.NewGauge(
		prometheus.GaugeOpts{

			Name: "logi_gauge_time_read_lines",
			Help: "Time read lines msec",
		})

	counterES = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "logi_counter_equal_strings",
			Help: "search lines counter equal string",
		})
) */
/*
func init() {
	prometheus.MustRegister(counter)
	prometheus.MustRegister(gaugeRT)
	prometheus.MustRegister(counterES)
}
*/
type Scan struct {
	Find          string
	Text          string
	ChRes         chan Data
	LimitResLines int // 0 - unlimited
}

type Data struct {
	ID   int32
	Line string
}

func (t *Scan) procLineSearch(lineS string) (csvF string, xmlL string) {

	if len(lineS) == 0 {
		return
	}

	xmlline := DecodeLine(lineS)
	val, err := DecodeXML(xmlline)
	if err != nil {
		logs.ErrorLogger.Println("procLineSearch" + err.Error())

		return
	}
	csvline := EncodeCSV(val)

	return csvline, val.XML_RECORD_ROOT[0].XML_MESSAGE
}

func (t *Scan) procFileSearch(file string, wg *sync.WaitGroup, countL *int32, numId *int32) {

	ch := make(chan Data, 100) // буфер для ускорения записи (нет блокировки чтением)
	var csvSt Data
	var mes string
	t.LimitResLines = t.LimitResLines + 1000
	for i := runtime.NumCPU() + 1; i > 0; i-- {

		go func() {
			wg.Add(1)
			defer wg.Done()
			//End:
		ext:
			for {
				select {
				case data, ok := <-ch:
					if !ok {
						break ext
					}

					if t.LimitResLines <= int(atomic.LoadInt32(countL)) && t.LimitResLines != 0 {
						break
					}

					csvSt.Line, mes = t.procLineSearch(data.Line)
					csvSt.ID = data.ID
					s := string(mes)

					if strings.Contains(s, t.Text) {
						//counterES.Add(1)
						t.ChRes <- csvSt
						//i++
						atomic.AddInt32(countL, 1)

					}
					//	counter.Add(1)

				}

			}

		}()
	}

	//start := time.Now()

	err := ReadLines(file, func(line string) {
		data := Data{Line: line, ID: atomic.AddInt32(numId, 1)}
		ch <- data

	})
	if err != nil {
		logs.ErrorLogger.Println("ReadLines procFileSearch" + err.Error())

		fmt.Println("ReadLines: ", err)
		close(ch)
		return
	}

	close(ch)

	//	duration := time.Since(start)

	//	gaugeRT.Add(float64(duration.Milliseconds()))

}

func (t *Scan) Search() {
	var wg sync.WaitGroup
	var countL int32
	var numId int32
	filepath.Walk(t.Find,
		func(path string, file os.FileInfo, err error) error {

			if err != nil {
				logs.ErrorLogger.Println("Search" + err.Error())

				return err
			}
			if !file.IsDir() {

				t.procFileSearch(path, &wg, &countL, &numId)
			}

			return nil
		})

	wg.Wait()
}
