package tcp

import (
	"net"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type TcpListener struct {
	Name              string
	Protocol          string
	Port              string
	Logger            zerolog.Logger
	ConnectionHandler func(net.Conn)
}

func StartListener(opts TcpListener) {
	listener, err := net.Listen(opts.Protocol, ":"+opts.Port)

	if err != nil {
		opts.Logger.Err(err).Msg("failed to start net listener")
		return
	}

	opts.Logger.Info().Msgf("%s running on port %s", opts.Name, opts.Port)

	defer func() {
		err := listener.Close()
		log.Err(err).Msg("failed to close listener")
	}()

	for {
		connection, err := listener.Accept()
		if err != nil {
			opts.Logger.Err(err).Msg("failed to accept connection")
			return
		}

		go opts.ConnectionHandler(connection)
	}
}
