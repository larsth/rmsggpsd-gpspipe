package daemon

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/larsth/go-rmsggpsbinmsg"
	"github.com/larsth/rmsg_gpsd_tcp/expvar"
	"github.com/larsth/rmsggpsd-gpspipe/errors"
)

var (
	httpd  http.Server
	router mux.Router
)

// Configure, and start the web server ...
func startHttpd(c *JsonConfig) error {
	var (
		err error
	)

	(&router).HandleFunc("/", rootWebRequestHandler)
	(&router).HandleFunc("/debug/expvar", expvar.ExpvarHandler)

	(&httpd).Addr = c.HttpAddr
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

func Start(c *JsonConfig) error {
	var (
		err        error
		binMessage *binmsg.Message
	)

	if err = readJsonConfigDocument(c); err != nil {
		return errors.Annotate(err,
			`Error reading the JSON configuration file.`)
	}

	if len(c.HttpAddr) == 0 {
		return errors.Annotate(err,
			"The HTTP address is empty in the JSON configuration document")
	}

	if c.GpsPipeCmd == nil && c.Gps == nil {
		return errors.Annotatef(err, "%s: %s!",
			`Both the \"gpspipe\", and the \"nogps" JSON objects does not exists`,
			`Nothing to do`)
	}

	if c.Gps != nil {
		binMessage = new(binmsg.Message)
		binMessage.TimeStamp.Time = time.Now().UTC()
		binMessage.Gps.Altitude = c.Gps.Alt
		binMessage.Gps.Latitude = c.Gps.Lat
		binMessage.Gps.Longitude = c.Gps.Lon
		binMessage.Gps.FixMode = c.Gps.FixMode
		cache.Put(binMessage)
	}

	if c.GpsPipeCmd != nil {
		if err = c.GpsPipeCmd.init(); err != nil {
			return errors.Annotate(err,
				`Cannot init external gpspipe command.`)
		}

		if err = c.GpsPipeCmd.run(); err != nil {
			return errors.Annotate(err,
				`Cannot start external gpspipe command.`)
		}
	}

	//Last thing to do:
	if err = startHttpd(c); err != nil {
		return errors.Annotate(err, "Cannot start the web server.")
	}

	return nil
}
