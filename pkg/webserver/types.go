package webserver

import (
    "net/http"
)

type WebServerInt interface {
    Index(w http.ResponseWriter, r *http.Request)
    Filter(w http.ResponseWriter, r *http.Request)
    Prioritize(w http.ResponseWriter, r *http.Request)
    Start()
    Stop()
}

// Web Server parameters
type WebServerParameters struct {
    Port           int    // webhook server port
    CertFile       string // path to the x509 certificate for https
    KeyFile        string // path to the x509 private key matching `CertFile`
}

type webServer struct {
    server *http.Server
}

