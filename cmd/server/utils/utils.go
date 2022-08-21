package utils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// 3. read server cert & key
func ServerCert() (*tls.Certificate, error) {
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	return &serverCert, nil
}

func ConfigTLS(serverCert tls.Certificate, certPool *x509.CertPool) *tls.Config {
	// 5. configuration of the certificate what we want to
	return &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
		ClientCAs:    certPool,
	}
}

func GetCertPool(caFile string) (*x509.CertPool, error) {
	// 1.read ca's cert, verify to client's certificate
	caPem, err := os.ReadFile(caFile) // "cert/ca-cert.pem"
	if err != nil {
		return nil, err
	}

	// 2. create cert pool and append ca's cert
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, errors.New("bad cert")
	}

	return certPool, nil
}
