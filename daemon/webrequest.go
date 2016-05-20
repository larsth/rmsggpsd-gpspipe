package daemon

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/larsth/go-rmsggpsbinmsg"
	"github.com/larsth/rmsggpsd-gpspipe/cache"
)

const rfc7231 = `Mon, 06 Jan 2006 15:04:05 GMT`

func isGETHttpMethod(req *http.Request, w http.ResponseWriter) (ok bool) {
	var msg string

	if strings.Compare("GET", req.Method) != 0 {
		msg = fmt.Sprintf("Method %s is not allowed", req.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		log.Println(msg)
		return false
	}
	return true
}

func parseForm(req *http.Request, w http.ResponseWriter) (ok bool) {
	var (
		err error
		msg string
	)
	if err = req.ParseForm(); err != nil {
		msg = fmt.Sprintf("Cannot parse HTTP GET query, Error %s.", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		HttpdLogger.Println(msg)
		return false
	}
	return true
}

func writeXWwwFormUrlencodedHttpResponse(
	w http.ResponseWriter,
	nowUTC time.Time,
	thisgpscache *cache.BinMsg,
	bearingcache *cache.Bearing) {
	var (
		m                               *binmsg.Message
		fixmode, alt, lat, lon, gpstime string
		tBearing                        string
		bearing                         string
		values                          url.Values
		p                               []byte
		pLen                            string
	)

	m = thisgpscache.Get()
	fixmode, alt, lat, lon, gpstime = m.Strings()
	_, _, bearing, tBearing = bearingcache.Get()

	values.Set("bearing", bearing)
	values.Set("bearingtime", tBearing)
	values.Set("gpsaltitude", alt)
	values.Set("gpsfixmode", fixmode)
	values.Set("gpslatitude", lat)
	values.Set("gpslongitude", lon)
	values.Set("gpstime", gpstime)

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Cache-Control", "no-cache")

	w.Header().Set("Date", nowUTC.Format(rfc7231))
	w.Header().Set("Date-RFC-3339", nowUTC.Format(time.RFC3339))
	w.Header().Set("Date-RFC3339-Nano", nowUTC.Format(time.RFC3339Nano))

	p = []byte(values.Encode())
	pLen = strconv.Itoa(len(p))
	w.Header().Set("Content-Length", pLen)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(p)

	return
}

func httpRequestHandler(w http.ResponseWriter, req *http.Request) {
	var nowUTC = time.Now().UTC()

	if !isGETHttpMethod(req, w) {
		return //response had already been written
	}
	if !parseForm(req, w) {
		return //response had already been written
	}
	//NOTE: thisGpsCache, and bearingCache are package visible variables
	writeXWwwFormUrlencodedHttpResponse(w, nowUTC, thisGpsCache, bearingCache)
}
