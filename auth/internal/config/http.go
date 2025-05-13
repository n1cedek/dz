package env

import (
	"errors"
	"net"
	"os"
)

var _ HTTPConfig = (*httpConfig)(nil)

const (
	httpPort = "HTTP_PORT"
	httpHost = "HTTP_HOST"
)

type httpConfig struct {
	host string
	port string
}

func NewHttpConfig() (*httpConfig, error) {
	port := os.Getenv(httpPort)
	if len(port) == 0 {
		return nil, errors.New("failed to get http port")
	}
	host := os.Getenv(httpHost)
	if len(host) == 0 {
		return nil, errors.New("failed to get http port")
	}
	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

func (h *httpConfig) Address() string {
	return net.JoinHostPort(h.host, h.port)
}
