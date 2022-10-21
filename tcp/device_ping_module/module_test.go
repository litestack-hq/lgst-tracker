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

var testPort int
var testLogger zerolog.Logger
var testModule tcp.Module

func TestMain(m *testing.M) {
	var err error
	testLogger = zerolog.New(os.Stdout)
	testPort, err = utils.GetFreePort()
	if err != nil {
		testLogger.Panic().Err(err).Msg("get free port failed")
	}

	testModule = New(tcp.ModuleOpts{
		Name:        "Test device ping module",
		Protocol:    tcp.PROTOCOL_TCP,
		DefaultPort: fmt.Sprintf("%d", testPort),
		Logger:      testLogger,
	})

	go testModule.Run()

	m.Run()

	err = testModule.Close()
	if err != nil {
		testLogger.Panic().Err(err).Msg("failed to close TCP listener")
	}
}

func TestTcpConnection(t *testing.T) {
	url := fmt.Sprintf("0.0.0.0:%d", testPort)
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
