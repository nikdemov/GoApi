package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	logs "nikworkedprofile/GoApi/src/logs_app"
	"nikworkedprofile/GoApi/src/server"

	//logs "github.com/nikworkedprofile/GoApi/src/logs_app"

	"github.com/kardianos/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const serviceName = "Logi2 version 0.1"
const serviceDescription = "Service for monitoring"

var (
	uptime_server = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "logi2_uptime_server_seconds",
		Help: "Time server run",
	})
)
var ctx, _ = context.WithCancel(context.Background())
var empty []string

type program struct{}

// playType indicates how to play a gauge.
func (p program) Start(s service.Service) error {
	logs.InfoLogger.Println("Service Start")
	fmt.Println(s.String() + " started")
	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	logs.InfoLogger.Println("Service Stop")
	fmt.Println(s.String() + " stopped")
	return nil
}

func (p program) run() {
	startTime := time.Now()

	go func() {
		time.Sleep(time.Second * 2)
		for {
			uptime := time.Since(startTime)
			uptime_server.Set(float64(int64(uptime / time.Second)))
		}
	}()
	server.Server()

}

// WaitShutdown waits until is going to die
func WaitShutdown() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sigc
	log.Printf("signal received [%v] canceling everything\n", s)
}

func main() {
	logs.InitLog()
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		logs.ErrorLogger.Println("Cannot create the service: " + err.Error())
	}
	err = s.Run()
	if err != nil {
		logs.ErrorLogger.Println("Cannot start the service:" + err.Error())
	}
	defer s.Run()
	WaitShutdown()

}
