package daemon

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/larsth/go-rmsggpsbinmsg"
	_ "github.com/larsth/rmsggpsd-gpspipe/cache"
	"github.com/larsth/rmsggpsd-gpspipe/errors"
)

var (
	httpd  http.Server
	router mux.Router
)

var HttpdLogger *log.Logger

func startGpsppipe(gpspipeLogger *log.Logger, c *JsonConfig) error {
	var (
		err error
	)
	if c.GpsPipe != nil {
		if err = c.GpsPipe.Run(gpspipeLogger); err != nil {
			return errors.Annotate(err,
				`Cannot start external gpspipe command.`)
		}
	}
	go gpspipeGoRoutine(c.GpsPipe)
	return nil
}

// Configure, and start the web server ...
func startHttpd(addr string) error {
	var (
		err error
	)

	(&router).HandleFunc("/", httpRequestHandler)

	(&httpd).Addr = addr
	(&httpd).Handler = &router
	(&httpd).ReadTimeout = 30 * time.Second
	(&httpd).WriteTimeout = 30 * time.Second
	/* MaxHeaderBytes: 1 << 12 = 4096 bytes, which is usually one OS page */
	(&httpd).MaxHeaderBytes = 1 << 12
	(&httpd).SetKeepAlivesEnabled(true)

	if err = (&httpd).ListenAndServe(); err != nil {
		return errors.Annotate(err,
			`Cannot start the web server.`)
	}
	return nil
}

func Start(gpspipeLogger, otherGpslogger *log.Logger, c *JsonConfig) error {
	var (
		err error
		m   *binmsg.Message
	)

	if err = readJsonConfigDocument(c); err != nil {
		return errors.Annotate(err,
			`Error reading the JSON configuration file.`)
	}

	if c.GpsPipe == nil && c.ThisGps == nil {
		return errors.Annotatef(err, "%s, &s: Nothing to do!",
			`Both the "gpspipe"`,
			`and the "this-gps" JSON objects does not exists`)
	}

	if c.ThisGps != nil {
		m = binmsg.MkBinMsg(c.ThisGps.Alt,
			c.ThisGps.Lat,
			c.ThisGps.Lon,
			c.ThisGps.FixMode,
			time.Now().UTC())
		thisGpsCache.Put(m)
	}

	//start the 'other GPS' HTTP client
	go httpClient(c, otherGpslogger)

	if err = startGpsppipe(gpspipeLogger, c); err != nil {
		return errors.Trace(err)
	}

	//Last thing to do:
	if err = startHttpd(c.Httpd.AddrString); err != nil {
		return errors.Trace(err)
	}

	return nil
}
