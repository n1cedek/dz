package env

import (
	"errors"
	"google.golang.org/grpc/credentials"
	"os"
)

const (
	tlsFile = "TLS_CERT_FILE"
	tlsKey  = "TLS_KEY_FILE"
)

var _ TLSConfig = (*tlsConfig)(nil)

type tlsConfig struct {
	certFile string
	keyFile  string
}

func NewTLSConfig() (*tlsConfig, error) {
	certF := os.Getenv(tlsFile)
	if len(certF) == 0 {
		return nil, errors.New("failed to get cert file")
	}
	keyF := os.Getenv(tlsKey)
	if len(keyF) == 0 {
		return nil, errors.New("failed to get key file")
	}
	return &tlsConfig{
		certFile: certF,
		keyFile:  keyF,
	}, nil
}

func (t *tlsConfig) GetTLSConfig() (credentials.TransportCredentials, error) {
	cred, err := credentials.NewServerTLSFromFile(t.certFile, t.keyFile)
	if err != nil {
		return nil, err
	}
	return cred, nil
}
