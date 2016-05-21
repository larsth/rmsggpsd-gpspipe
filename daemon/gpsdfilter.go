package daemon

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/larsth/go-gpsdjson"
	"github.com/larsth/go-rmsggpsbinmsg"
	"github.com/larsth/rmsggpsd-gpspipe/errors"
)

type filter interface {
	ParseGpsdJson(p []byte) (*binmsg.Message, error)
}

type gpsdFilter struct {
	mutex       sync.RWMutex
	logger      *log.Logger
	gpsdJsonTpv *gpsdjson.TPV
}

type Class struct {
	Class string `json:"class"`
}

func newGpsdFilter(logger *log.Logger) *gpsdFilter {
	g := new(gpsdFilter)
	g.logger = logger
	return g
}

func (g *gpsdFilter) tpvGpsdJsonToBinMessage() (*binmsg.Message, error) {
	var (
		m   *binmsg.Message
		t   time.Time
		err error
	)

	if len(g.gpsdJsonTpv.Time) > 0 {
		t, err = time.Parse(time.RFC3339Nano, g.gpsdJsonTpv.Time)
		if err != nil {
			return nil, errors.Annotate(err,
				"Cannot parse gpsd TPV JSON document Time string.")
		}
	} else {
		t = time.Unix(0, 0)
	}

	m = binmsg.MkBinMsg(float32(g.gpsdJsonTpv.Alt),
		float32(g.gpsdJsonTpv.Lat),
		float32(g.gpsdJsonTpv.Lon),
		g.gpsdJsonTpv.Fix,
		t)

	return m, nil
}

func (g *gpsdFilter) ParseGpsdJson(p []byte) (*binmsg.Message, error) {
	var (
		class Class
		err   error
		m     *binmsg.Message
	)

	g.mutex.Lock()
	g.mutex.Unlock()

	//log the gpsd JSON document
	g.logger.Printf("#v", p)

	if err = json.Unmarshal(p, &class); err != nil {
		return nil, errors.Annotate(err,
			"Cannot parse gpsd JSON document."+
				" Finding \"class\" failed.")
	}

	if strings.Compare("TPV", class.Class) == 0 {
		if err = json.Unmarshal(p, g.gpsdJsonTpv); err != nil {
			return nil, errors.Annotate(err,
				"Cannot unmarshal a gpsd TPV JSON document.")
		}
		if m, err = g.tpvGpsdJsonToBinMessage(); err != nil {
			return nil, errors.Annotate(err,
				`Cannot create a *binmsg.Message `+
					`with data from a gpsd TPV JSON document.`)
		}
		return m, nil
	}

	return nil, nil
}
