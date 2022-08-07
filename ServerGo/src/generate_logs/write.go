package generate_logs

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	"time"

	"nikworkedprofile/GoApi/ServerGo/src/logenc"
	logs "nikworkedprofile/GoApi/ServerGo/src/logs_app"
	"nikworkedprofile/GoApi/ServerGo/src/web/util"
)

var (
	ctxCF, cancelCF = context.WithCancel(context.Background())
	letterRunes     = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func ProcGenWriteF(ctx context.Context) {
	//Example()

	filesFrom := string(util.GetOutboundIP()[len(util.GetOutboundIP())-3:])
	logenc.CreateDir(pathdata + "/repdata/")
	file, err := os.OpenFile(pathdata+"/repdata/gen_files_write"+filesFrom, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logs.FatalLogger.Println("Create gen file" + err.Error())
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop Gen")
			return
		default:
			LINE := StructFile("101")

			rand.Seed(time.Now().UnixNano())

			InfoLogger := log.New(file, "", 0)

			infof := func(info string) {
				InfoLogger.Output(2, logenc.EncodeLine(info))
			}

			infof(LINE)

			//time.Sleep(time.Nanosecond * 1000000)
			time.Sleep(2000 * time.Millisecond)
			fmt.Println("Message add :D" + RandStringRunes(3))

		}

	}

}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
