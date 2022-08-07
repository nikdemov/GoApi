package util

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	logs "nikworkedprofile/GoApi/src/logs_app"
)

// Config - struct to hold the config
type Config struct {
	Dir1  []string
	Dir   []string
	Cron  string
	Cron1 string
}

// Conf global config
var Conf Config

// ParseConfig - function to manage config
func ParseConfig(dir []string, cron string) error {
	// Parse cron
	// Rules for cron :
	// the string should be of type [^0](\d*)(h|d) and the integer should be positive
	// If this exact format is not presented, it will fail.

	timeUnit := cron[len(cron)-1]
	if timeUnit != 'h' && timeUnit != 'd' {
		return fmt.Errorf("invalid time unit in cron arg: %s", cron)
	}

	timeValue, err := strconv.ParseInt(cron[:len(cron)-1], 10, 32)
	if err != nil {
		logs.ErrorLogger.Println("invalid time value" + err.Error())
		return fmt.Errorf("invalid time value in cron arg: %s", cron)
	}
	if timeValue < 0 {
		return fmt.Errorf("invalid time value in cron arg: %s", cron)
	}

	if (timeUnit == 'h' && timeValue >= 10000) || (timeUnit == 'd' && timeValue >= 365) {
		fmt.Fprintf(os.Stderr, "Whoah Dude !, That's a long time you put there...")
	}

	// First Index
	IndexFiles(dir)
	tmp := make([]interface{}, len(dir))
	for idx, x := range dir {
		tmp[idx] = x
		log.Println(x)
	}

	// Setting up cron job to keep indexing the files
	if timeValue > 0 {
		repeat := time.Duration(timeValue) * time.Hour

		if timeUnit == 'd' {
			repeat = repeat * 24
		}
		//fmt.Println(repeat)
		//	go MakeAndStartCron(repeat, func(v ...interface{}) error { !!!
		go MakeAndStartCron(repeat, func(v ...interface{}) error {
			tmp := make([]string, len(v))
			for idx, val := range v {
				tmp[idx] = val.(string)
			}
			IndexFiles(tmp)
			return nil
		}, tmp...)
	}
	return nil
}
