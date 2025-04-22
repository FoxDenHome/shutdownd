package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func GenerateSelfSignedCert(certFile string, hostname string) error {
	var err error
	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			return fmt.Errorf("could not get hostname: %v", err)
		}
	}
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return fmt.Errorf("could not generate private key: %v", err)
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
		return fmt.Errorf("failed to self-sign certificate: %s", err)
	}
	privateBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return fmt.Errorf("unable to marshal private key: %v", err)
	}
	out, err := os.Create(certFile)
	if err != nil {
		return fmt.Errorf("unable to open cert.pem for writing: %v", err)
	}
	defer func() {
		_ = out.Close()
	}()

	err = pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if err != nil {
		return fmt.Errorf("failed to PEM encode CERTIFICATE: %v", err)
	}
	err = pem.Encode(out, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privateBytes})
	if err != nil {
		return fmt.Errorf("failed to PEM encode EC PRIVATE KEY: %v", err)
	}

	return nil
}
