package http

import (
	"github.com/litestack-hq/lgst-tracker/helpers/config"
	"github.com/rs/zerolog"
)

type HttpModule struct {
	Conf   *config.Config
	Logger zerolog.Logger
}

func NewHttpModule(opts HandlerOpts) *HttpModule {
	h := HttpModule(opts)
	return &h
}
