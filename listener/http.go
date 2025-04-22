package listener

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/FoxDenHome/shutdownd/util"
)

func (h *Listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		_, _ = w.Write([]byte("POST only"))
		return
	}

	switch r.URL.Path {
	case "/shutdown":
		_ = h.Logger.Info(1, "Shutdown initiated")
		err := h.doShutdown()
		if err != nil {
			_ = h.Logger.Error(1, fmt.Sprintf("Shutdown start error: %v", err))
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	case "/abort":
		_ = h.Logger.Info(1, "Shutdown aborted")
		err := h.doShutdownAbort()
		if err != nil {
			_ = h.Logger.Error(1, fmt.Sprintf("Shutdown abort error: %v", err))
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	default:
		w.WriteHeader(404)
		_, _ = w.Write([]byte("Path not mapped"))
		return
	}

	w.WriteHeader(200)
	_, _ = w.Write([]byte("OK"))
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func (h *Listener) execute() (ssec bool, errno uint32) {
	defer func() {
		_ = h.Logger.Close()
	}()

	configDir, err := util.GetConfigDir(h.Logger)
	if err != nil {
		_ = h.Logger.Error(1, fmt.Sprintf("Could not locate config.json: %v", err))
		return
	}

	certFile := path.Join(configDir, "cert.pem")
	if !fileExists(certFile) {
		_ = h.Logger.Info(1, "Generating new certificate")
		err = util.GenerateSelfSignedCert(certFile, "")
		if err != nil {
			_ = h.Logger.Error(1, fmt.Sprintf("Could not generate new certificate: %v", err))
			return
		}
		_ = h.Logger.Info(1, "Successfully generated new certificate")
	}

	caCert, err := os.ReadFile(path.Join(configDir, "server.pem"))
	if err != nil {
		_ = h.Logger.Error(1, fmt.Sprintf("Could not read server.pem: %v", err))
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	server := &http.Server{
		Addr:    ":666",
		Handler: h,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  caCertPool,
		},
	}

	h.onReady(server)

	_ = h.Logger.Info(1, "ShutdownD listening")

	err = server.ListenAndServeTLS(certFile, certFile)
	if err != nil {
		_ = h.Logger.Error(1, fmt.Sprintf("Could not listen on HTTP: %v", err))
		return
	}
	return
}
