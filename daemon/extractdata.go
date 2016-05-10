package daemon

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/larsth/go-gpsdfilter"
	"github.com/larsth/go-gpsdjson"
	"github.com/larsth/go-gpsdparser"
	"github.com/larsth/go-gpsdscanner"
	"github.com/larsth/go-gpsfix"
	"github.com/larsth/go-rmsggpsbinmsg"
	"github.com/larsth/rmsggpsd/errors"
)

type extractData struct {
	mutex     sync.Mutex
	ioReader  io.Reader
	scanner   *gpsdscanner.Scanner
	filter    *gpsdfilter.Filter
	parseArgs gpsdparser.ParseArgs
}

type extractedData struct {
	Time             time.Time
	GpsdErrorMessage string
	Message          *binmsg.Message
}

func newExtractData(r io.Reader) (*extractData, error) {
	var err error

	ed := new(extractData)
	ed.ioReader = r
	if ed.scanner, err = gpsdscanner.New(r); err != nil {
		return nil, errors.Annotate(err,
			"gpsdscanner.New(io.Reader) error")
	}
	ed.filter = gpsdfilter.New()
	if err = addFilterRules(ed.filter); err != nil {
		return nil, errors.Annotate(err,
			"daemon.addFilterRules(*go-gpsdfilter.Filter) error")
	}
	ed.parseArgs.Filter = ed.filter

	return ed, nil
}

func (ed *extractData) ExtractData() (*extractedData, error) {
	var (
		err          error
		annotatedErr error
		i            interface{}
		TPV          *gpsdjson.TPV
		ERROR        *gpsdjson.ERROR
		eData        *extractedData
	)

	ed.mutex.Lock()
	defer ed.mutex.Unlock()

	if _, ed.parseArgs.Data, err = ed.scanner.Scan(); err != nil {
		ed.parseArgs.Data = nil
		annotatedErr = errors.Annotate(err,
			"scanning for a line (a \\n or \\r\\n terminated line) error")
		return nil, annotatedErr
	}
	if _, err = ed.filter.Filter(ed.parseArgs.Data); err != nil {
		ed.parseArgs.Data = nil
		annotatedErr = errors.Annotate(err,
			"Cannot filter gpsd JSON document")
		return nil, annotatedErr
	}
	if i, err = gpsdparser.Parse(&ed.parseArgs); err != nil {
		ed.parseArgs.Data = nil
		annotatedErr = errors.Annotate(err,
			"Cannot parse the filtered gpsd JSON document")
		return nil, annotatedErr
	}

	switch i.(type) {
	case *gpsdjson.TPV:
		TPV = i.(*gpsdjson.TPV)
		if eData, err = extractGpsdJsonTPV(TPV); err != nil {
			annotatedErr = errors.Annotate(err,
				"Cannot extract a class=TPV gpsd JSON document")
			return nil, annotatedErr
		}
		return eData, nil

	case *gpsdjson.ERROR:
		ERROR = i.(*gpsdjson.ERROR)

		eData = extractGpsdJsonERROR(ERROR)
		return eData, nil
	}
	//Below: this is a implicit switch i.(type) default: case ...

	//A gpsd JSON document class from a filter rule is unhandled
	// - this is a programming error.

	//The program can do nothing about a programming error, and ignoring
	//the problem is not acceptable, so:
	log.SetFlags(log.Ldate | log.Llongfile | log.Ltime | log.LUTC | log.Lmicroseconds)
	annotatedErr = errors.Errorf("%s: %s - %s. %s. (Type: %s)",
		"func (ed *extractData) ExtractData() (*ExtractedData, error)",
		"switch i.(type), case implicit default",
		"Unhandled gpsdjson JSON document type",
		"This is a programming error",
		fmt.Sprintf("%#v", i))
	log.Println(annotatedErr.Error())

	return nil, annotatedErr
}

func extractGpsdJsonTPV(tpv *gpsdjson.TPV) (*extractedData, error) {
	var (
		err          error
		annotatedErr error
		eData        = new(extractedData)
	)
	eData.Message = new(binmsg.Message)
	eData.Time = time.Now().UTC()
	eData.Message.Gps.FixMode = tpv.Fix

	if eData.Message.Gps.FixMode == gpsfix.FixNone ||
		eData.Message.Gps.FixMode == gpsfix.FixNotSeen {
		eData.Message.TimeStamp.Time = eData.Time
		eData.Message.Gps.Latitude = 0.0
		eData.Message.Gps.Longitude = 0.0
		eData.Message.Gps.Altitude = 0.0
	} else {
		if eData.Time, err = time.Parse(time.RFC3339, tpv.Time); err != nil {
			annotatedErr = errors.Annotate(err,
				"Error parsing to a time.Time from the 'tpv.Time' string")
			return nil, annotatedErr
		}
		eData.Message.TimeStamp.Time = eData.Time
		eData.Message.Gps.SetLat(tpv.Lat)
		eData.Message.Gps.SetLon(tpv.Lon)
		eData.Message.Gps.SetAlt(tpv.Alt)
	}

	return eData, nil
}

func extractGpsdJsonERROR(ERROR *gpsdjson.ERROR) *extractedData {
	var eData = new(extractedData)

	eData.Time = time.Now().UTC()
	eData.Message = nil
	eData.GpsdErrorMessage = ERROR.Message

	return eData
}
