package sockets

import (
	"bufio"
	"net"
)

func DevicePingHandler(connection net.Conn, opts HandlerOpts) {
	/*
		TODO
		Validate message before parsing
		- Heartbeat
		- GPS ping
		- NBR
		- ...etc
	*/

	defer connection.Close()

	for {
		data, err := bufio.NewReader(connection).ReadString('#')
		if err != nil {
			if err.Error() != "EOF" {
				opts.Logger.Err(err).Msg("failed to read data")
			}
			return
		}

		opts.Logger.Info().Str("data", data).Msg("device data")
	}
}
