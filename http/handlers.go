package http

import (
	"io"
	"net/http"

	"github.com/go-chi/render"
)

func (m *HttpModule) PingHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		m.Logger.Err(err).Msg("failed to read ping data")
		render.Status(r, http.StatusBadRequest)
		return
	}

	m.Logger.Info().Str("data", string(body)).Msg("ping received")

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "Ok")
}
