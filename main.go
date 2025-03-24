package main

import (
	closer "chaosAgent/conn/close"
	chaoshttp "chaosAgent/pkg/http"
	"chaosAgent/pkg/log"
	"chaosAgent/pkg/options"
	"chaosAgent/pkg/tools"
	"chaosAgent/transport"
	"chaosAgent/web/api"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

var pidFile = "/var/run/chaos.pid"

func main() {
	options.NewOptions()
	log.InitLog(&options.Opts.LogConfig)

	options.Opts.SetOthersByFlags()

	// new transport newConn
	clientInstance, err := chaoshttp.NewHttpClient(options.Opts.TransportConfig)
	if err != nil {
		logrus.Errorf("create transport client instance failed, err: %s", err.Error())
		handlerErr(err)
	}
	transportClient := transport.NewTransportClient(clientInstance)
	api1 := api.NewAPI()
	err = api1.Register(transportClient)

	if err != nil {
		logrus.Errorf("register api failed, err: %s", err.Error())
		handlerErr(err)
	}

	// listen server
	go func() {
		defer tools.PanicPrintStack()
		err := http.ListenAndServe(":"+options.Opts.Port, nil)
		if err != nil {
			logrus.Warningln("Start http server failed")
			handlerErr(err)
		}
	}()

	handlerSuccess()

	closeClient := closer.NewClientCloseHandler(transportClient)
	tools.Hold(closeClient)
}

func handlerSuccess() {
	pid := os.Getpid()
	err := writePid(pid)
	if err != nil {
		logrus.Panic("write pid: ", pidFile, " failed. ", err)
	}
}

func handlerErr(err error) {
	if err == nil {
		return
	}
	logrus.Warningf("start agent failed because of %v", err)
	writePid(-1)
	logrus.Errorf("chaos agent will exit")
	os.Exit(1)
}

func writePid(pid int) error {
	file, err := os.OpenFile(pidFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(strconv.Itoa(pid))
	return err
}
