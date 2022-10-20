package tcp

import "github.com/rs/zerolog"

const (
	PROTOCOL_TCP = "tcp"
	PROTOCOL_UDP = "udp"
)

type ModuleOpts struct {
	Name        string
	Protocol    string
	DefaultPort string
	Logger      zerolog.Logger
}

type Module interface {
	Run()
	Close() error
}
