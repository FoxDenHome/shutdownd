package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {
	certFile := "server.pem"

	log.Print("Generating new certificate")
	// Generate new certificate
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Printf("Could not generate private key: %v", err)
		return
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"ShutdownD"},
			CommonName:   "server",
		},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365 * 10),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Printf("Failed to self-sign certificate: %s", err)
		return
	}
	privateBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		log.Printf("Unable to marshal private key: %v", err)
		return
	}
	out, err := os.Create(certFile)
	if err != nil {
		log.Printf("Unable to open cert.pem for writing: %v", err)
		return
	}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	pem.Encode(out, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privateBytes})
	_ = out.Close()
}
