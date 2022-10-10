package http

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/litestack-hq/lgst-tracker/helpers/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	t.Run("Handle ping without errors", func(t *testing.T) {
		logger := zerolog.New(os.Stdout)
		m := HttpModule{
			Logger: logger,
			Conf:   &config.Config{},
		}

		body := []byte("PING!")

		r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))
		w := httptest.NewRecorder()

		m.PingHandler(w, r)

		res := w.Result()
		responseBytes, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, 0, len(responseBytes))
	})
}
