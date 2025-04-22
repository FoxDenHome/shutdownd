package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path"
	"time"
)

type Logger interface {
	Info(eventID uint32, msg string) error
	Error(eventID uint32, msg string) error
	Close() error
}

func (h *shutdownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Write([]byte("POST only"))
		return
	}

	switch r.URL.Path {
	case "/shutdown":
		h.logger.Info(1, "Shutdown initiated")
		err := h.doShutdown()
		if err != nil {
			h.logger.Error(1, fmt.Sprintf("Shutdown start error: %v", err))
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	case "/abort":
		h.logger.Info(1, "Shutdown aborted")
		err := h.doShutdownAbort()
		if err != nil {
			h.logger.Error(1, fmt.Sprintf("Shutdown abort error: %v", err))
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	default:
		w.WriteHeader(404)
		w.Write([]byte("Path not mapped"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func (h *shutdownHandler) execute() (ssec bool, errno uint32) {
	defer h.logger.Close()

	configDir, err := getConfigDir(h.logger)
	if err != nil {
		h.logger.Error(1, fmt.Sprintf("Could not locate config.json: %v", err))
		return
	}

	certFile := path.Join(configDir, "cert.pem")
	if !fileExists(certFile) {
		h.logger.Info(1, "Generating new certificate")
		// Generate new certificate
		hostname, err := os.Hostname()
		if err != nil {
			h.logger.Error(1, fmt.Sprintf("Could not get hostname: %v", err))
			return
		}
		priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			h.logger.Error(1, fmt.Sprintf("Could not generate private key: %v", err))
			return
		}
		template := x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				Organization: []string{"ShutdownD"},
				CommonName:   hostname,
			},
			NotBefore: time.Now().Add(-time.Hour),
			NotAfter:  time.Now().Add(time.Hour * 24 * 365 * 10),

			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
		}
		certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
		if err != nil {
			h.logger.Error(1, fmt.Sprintf("Failed to self-sign certificate: %s", err))
			return
		}
		privateBytes, err := x509.MarshalECPrivateKey(priv)
		if err != nil {
			h.logger.Error(1, fmt.Sprintf("Unable to marshal private key: %v", err))
			return
		}
		out, err := os.Create(certFile)
		if err != nil {
			h.logger.Error(1, fmt.Sprintf("Unable to open cert.pem for writing: %v", err))
			return
		}
		pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
		pem.Encode(out, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privateBytes})
		_ = out.Close()
		h.logger.Info(1, "Successfully generated new certificate")
	}

	caCert, err := os.ReadFile(path.Join(configDir, "server.pem"))
	if err != nil {
		h.logger.Error(1, fmt.Sprintf("Could not read server.pem: %v", err))
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

	h.logger.Info(1, "ShutdownD listening")

	err = server.ListenAndServeTLS(certFile, certFile)
	if err != nil {
		h.logger.Error(1, fmt.Sprintf("Could not listen on HTTP: %v", err))
		return
	}
	return
}
