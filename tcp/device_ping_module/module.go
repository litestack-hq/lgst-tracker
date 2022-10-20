package device_ping_module

import (
	"bufio"
	"net"
	"strings"

	"github.com/litestack-hq/lgst-tracker/tcp"
	"github.com/rs/zerolog"
)

type devicePingModule struct {
	Name        string
	Protocol    string
	DefaultPort string
	Logger      zerolog.Logger
	Listener    net.Listener
	// ConnectionClosed bool
}

func New(opts tcp.ModuleOpts) tcp.Module {
	return &devicePingModule{
		Name:        opts.Name,
		Protocol:    opts.Protocol,
		DefaultPort: opts.DefaultPort,
		Logger:      opts.Logger,
	}
}

func (m *devicePingModule) connectionHandler(conn net.Conn) {
	defer conn.Close()

	for {
		data, err := bufio.NewReader(conn).ReadString('#')
		if heartBeatRegex.MatchString(data) {
			m.handleHeartBeat(data, conn)
			return
		}

		if err != nil {
			if err.Error() != "EOF" {
				m.Logger.Err(err).Msg("failed to read data")
			}
			return
		}

		m.Logger.Info().Str("data", data).Msg("device data")
	}
}

func (m *devicePingModule) Run() {
	var err error
	m.Listener, err = net.Listen(m.Protocol, ":"+m.DefaultPort)

	if err != nil {
		msg := err.Error()
		if !strings.Contains(msg, "use of closed network connection") {

			m.Logger.Err(err).Msg("failed to start net listener")
		}
		return
	}

	defer m.Listener.Close()

	m.Logger.Info().Msgf("%s running on port %s", m.Name, m.DefaultPort)

	for {
		connection, err := m.Listener.Accept()
		if err != nil {
			m.Logger.Err(err).Msg("failed to accept connection")
			return
		}

		go m.connectionHandler(connection)
	}
}

func (m *devicePingModule) Close() error {
	return m.Listener.Close()
}
