package device_ping_module

import (
	"bufio"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	SAMPLE_DEVICE_HEARTBEAT_RESPONSE = `*HQ,7301037602,V4V1,20221019154330#`
	SAMPLE_DEVICE_HEARTBEAT          = `*HQ,7301037602,V1,142852,V,0635.7945,N,00320.4042,E,000.00,246,191022,FFFFFBFF,621,30,0,0,6#`
)

func TestHandleHeartBeat(t *testing.T) {
	server, client := net.Pipe()
	module := devicePingModule{
		Logger: testLogger,
	}

	go module.handleHeartBeat(SAMPLE_DEVICE_HEARTBEAT, server)

	response, err := bufio.NewReader(client).ReadString('#')
	if err != nil {
		t.Error(err)
	}

	assert.Regexp(t, `\*HQ,\d{10},V4V1,\d{14}#`, response)
}
