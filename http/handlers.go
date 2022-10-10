package http

import (
	"io"
	"net/http"
)

func (m *HttpModule) PingHandler(w http.ResponseWriter, r *http.Request) {
	/*
		Validate message before parsing
		- Heartbeat
		- GPS ping
		- NBR
		- ...etc
	*/

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		m.Logger.Err(err).Msg("failed to read ping data")

		return
	}

	m.Logger.Info().Str("data", string(body)).Msg("device ping")

	// TODO: Handle ping
}
