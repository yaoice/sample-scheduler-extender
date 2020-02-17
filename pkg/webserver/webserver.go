package webserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/yaoice/sample-scheduler-extender/pkg/scheduler"
	"io"
	"k8s.io/klog"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
    "net/http"
	"sync"
)

var (
	once sync.Once
	ws   *webServer
	err  error
)

func NewWebServer(webHook WebServerParameters) (WebServerInt, error) {
	once.Do(func() {
		ws, err = newWebHookServer(webHook)
	})
	return ws, err
}

func newWebHookServer(webHook WebServerParameters) (*webServer, error) {
	// load tls cert/key file
	/*    tlsCertKey, err := tls.LoadX509KeyPair(webHook.CertFile, webHook.KeyFile)
	      if err != nil {
	          return nil, err
	      }
	*/
	ws := &webServer{
		server: &http.Server{
			Addr: fmt.Sprintf(":%v", webHook.Port),
			//            TLSConfig: &tls.Config{Certificates: []tls.Certificate{tlsCertKey}},
		},
	}

	// add routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", ws.Index)
	mux.HandleFunc("/filter", ws.Filter)
	mux.HandleFunc("/prioritize", ws.Prioritize)
	ws.server.Handler = mux
	return ws, nil
}

func (ws *webServer) Start() {
	if err := ws.server.ListenAndServe(); err != nil {
		klog.Errorf("Failed to listen and serve webhook server: %v", err)
	}
}

func (ws *webServer) Stop() {
	klog.Infof("Got OS shutdown signal, shutting down wenhook server gracefully...")
	ws.server.Shutdown(context.Background())
}

func (ws *webServer) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to sample-scheduler-extender!\n")
}

func (ws *webServer) Filter(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	var extenderArgs schedulerapi.ExtenderArgs
	var extenderFilterResult *schedulerapi.ExtenderFilterResult
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		extenderFilterResult = &schedulerapi.ExtenderFilterResult{
			Error: err.Error(),
		}
	} else {
		extenderFilterResult = scheduler.Filter(extenderArgs)
	}

	if response, err := json.Marshal(extenderFilterResult); err != nil {
		klog.Fatalln(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (ws *webServer) Prioritize(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	var extenderArgs schedulerapi.ExtenderArgs
	var hostPriorityList *schedulerapi.HostPriorityList
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		klog.Errorln(err)
		hostPriorityList = &schedulerapi.HostPriorityList{}
	} else {
		hostPriorityList = scheduler.Prioritize(extenderArgs)
	}

	if response, err := json.Marshal(hostPriorityList); err != nil {
		klog.Fatalln(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// func Bind(w http.ResponseWriter, r *http.Request) {}
