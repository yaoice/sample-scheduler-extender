package main

import (
	"flag"
	"github.com/yaoice/sample-scheduler-extender/pkg/webserver"
	"k8s.io/klog"
	"os"
	"os/signal"
	"syscall"
)

var webServer webserver.WebServerParameters

func main() {
	// parse parameters
	flag.Parse()

	// init webhook api
	ws, err := webserver.NewWebServer(webServer)
	if err != nil {
		panic(err)
	}

	// start webhook server in new routine
	go ws.Start()
	klog.Infoln("Server started")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	ws.Stop()
}

func init() {
	// read parameters
	flag.IntVar(&webServer.Port, "port", 8888, "The port of webhook server to listen.")
	flag.StringVar(&webServer.CertFile, "tlsCertPath", "/etc/webhook-demo/certs/cert.pem", "The path of tls cert")
	flag.StringVar(&webServer.KeyFile, "tlsKeyPath", "/etc/webhook-demo/certs/key.pem", "The path of tls key")
}


