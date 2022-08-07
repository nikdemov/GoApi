package logenc

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	logs "nikworkedprofile/GoApi/src/logs_app"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var callSearchCnt = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "logi2_rpc_request_seфrach_count",
	})

/* func init() {
	prometheus.MustRegister(callSearchCnt)
} */

type JsonRPC struct {
	workDir string
}

// Search API

type SearchArgs struct {
	Limmit int
	Find   string
}
type SearchResp struct {
	Lines []string
}

// Search RPC call
// example:
// curl -X POST -H "Content-Type: application/json" -d '{"id": 1, "method": "JsonRPC.Search", "params":[{"Limmit": 4,"Find":"s"}]}' http://localhost:8081/jrpc
func (t *JsonRPC) Search(args *SearchArgs, result *SearchResp) error {
	callSearchCnt.Add(1)
	resp := SearchResp{}
	fmt.Println("Search() in workDir", t.workDir)
	chRes := make(chan Data, 100)
	go func() {
		scan := &Scan{}
		scan.Find = t.workDir
		scan.Text = args.Find
		scan.ChRes = chRes
		scan.LimitResLines = args.Limmit
		scan.Search()
		close(scan.ChRes)
	}()

ext:
	for {
		select {
		case strm, ok := <-chRes:
			if !ok {
				break ext
			}
			resp.Lines = append(resp.Lines, strm.Line)
		}
	}

	//resp.Lines = []string{"1", "2"} // stub

	*result = resp
	return nil
}

// helper ReadWriteCloser interface for rpc handlers
type httpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *httpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *httpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *httpConn) Close() error                      { return nil }

// RunRPC run rpc routines
func RunRPC(workDir string) {
	fmt.Println("RunRPC()")

	uptime := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "logi2_rpc_uptime_seconds",
		})
	prometheus.MustRegister(uptime)

	go func() {
		for {
			uptime.Inc()
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	var listener net.Listener
	var err error
	listenErr := 0

	for ok := false; !ok; {
		listener, err = net.Listen("tcp", ":8081")
		if err != nil {
			if listenErr == 0 {
				fmt.Println(err)
			}
			listenErr++
			time.Sleep(time.Second * 3)
		}
		ok = (err == nil)

		if ok {
			defer listener.Close()
		}
	}

	fmt.Println("listen ok")

	rpcServer := rpc.NewServer()

	jrpc := &JsonRPC{workDir: workDir}
	rpcServer.Register(jrpc)

	serveJRPC := func(w http.ResponseWriter, r *http.Request) {
		serverCodec := jsonrpc.NewServerCodec(&httpConn{in: r.Body, out: w})
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		err := rpcServer.ServeRequest(serverCodec)
		if err != nil {
			logs.ErrorLogger.Println("Error while serving JSON request:" + err.Error())
			//log.Printf("Error while serving JSON request: %v", err)
			http.Error(w, "Error while serving JSON request, details have been logged.", 500)
			return
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/jrpc", serveJRPC).Headers("content-type", "application/json")
	//r.Handle("/metrics", promhttp.Handler())
	// cmd example for monitoring: expvarmon -ports="10010" -i 0.25s
	// or open in browser http://localhost:10010/debug/vars
	//r.Handle("/debug/vars", http.DefaultServeMux)

	if err := http.Serve(listener, r); err != nil {
		logs.ErrorLogger.Println("Http serve error" + err.Error())
	}
}
