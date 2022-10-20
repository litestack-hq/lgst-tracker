package device_ping_module

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/litestack-hq/lgst-common/helpers/utils"
	"github.com/litestack-hq/lgst-tracker/tcp"
	"github.com/rs/zerolog"
)

// Test that New works
// Test that run works
// Test that connects are accepted
// test that the heartbeat handler is called

var port int
var logger zerolog.Logger
var module tcp.Module

func TestMain(m *testing.M) {
	var err error
	logger = zerolog.New(os.Stdout)
	port, err = utils.GetFreePort()
	if err != nil {
		logger.Panic().Err(err).Msg("get free port failed")
	}

	module = New(tcp.ModuleOpts{
		Name:        "Test device ping module",
		Protocol:    tcp.PROTOCOL_TCP,
		DefaultPort: fmt.Sprintf("%d", port),
		Logger:      logger,
	})

	go module.Run()

	m.Run()

	err = module.Close()
	if err != nil {
		logger.Panic().Err(err).Msg("failed to close TCP listener")
	}
}

func TestTcpConnection(t *testing.T) {
	url := fmt.Sprintf("0.0.0.0:%d", port)
	addr, err := net.ResolveTCPAddr(tcp.PROTOCOL_TCP, url)
	if err != nil {
		t.Error(err)
	}

	conn, err := net.DialTCP(tcp.PROTOCOL_TCP, nil, addr)
	if err != nil {
		t.Error(err)
	}

	err = conn.Close()
	if err != nil {
		t.Error(err)
	}
}
