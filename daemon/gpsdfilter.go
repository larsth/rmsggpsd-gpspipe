package daemon

import (
	"log"
	"sync"

	"github.com/larsth/go-gpsdjson"
	"github.com/larsth/go-rmsggpsbinmsg"
)

type gpsdFilter struct {
	mutex         sync.RWMutex
	m             *binmsg.Message
	logger        *log.Logger
	gpsdJsonTpv   *gpsdjson.TPV
	gpsdJsonError *gpsdjson.ERROR
}

func newGpsdFilter(logger *log.Logger) *gpsdFilter {
	g := new(GpsdFilter)
	g.logger = logger
	g.m = new(binmsg.Message)
	return g
}

func (g *gpsdFilter) ParseGpsdJsonDocument(p []byte) error {
	return nil
}
