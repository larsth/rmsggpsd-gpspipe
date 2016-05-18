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

//write Date* HTTP headers
func writeDateResponseHeaders(w http.ResponseWriter, binMessage *binmsg.Message) {
	var nowUTC time.Time = time.Now().UTC()

	w.Header().Set("Date", nowUTC.Format(rfc7231))
	w.Header().Set("Date-RFC-3339", nowUTC.Format(time.RFC3339))
	w.Header().Set("Date-RFC3339-Nano", nowUTC.Format(time.RFC3339Nano))
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

func writeXWwwFormUrlencodedHttpResponse(w http.ResponseWriter) {
	var (
		m                               *binmsg.Message
		fixmode, alt, lat, lon, gpstime string
		bearingTime                     time.Time
		bearing                         float64
		bearingTimeString               string
		bearingString                   string
		values                          url.Values
		p                               []byte
		pLen                            string
	)

	m = thisGpsCache.Get()
	fixmode, alt, lat, lon, gpstime = m.Strings()
	bearing, bearingTime = bearingCache.Get()
	bearingString = strconv.FormatFloat(bearing, 'f', -1, 32)
	bearingTimeString = bearingTime.Format(time.RFC3339)

	values.Set("bearing", bearingString)
	values.Set("bearingtime", bearingTimeString)
	values.Set("gpsaltitude", alt)
	values.Set("gpsfixmode", fixmode)
	values.Set("gpslatitude", lat)
	values.Set("gpslongitude", lon)
	values.Set("gpstime", gpstime)

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Cache-Control", "no-cache")

	p = []byte(values.Encode())
	pLen = strconv.Itoa(len(p))
	w.Header().Set("Content-Length", pLen)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(p)

	return
}

func httpRequestHandler(w http.ResponseWriter, req *http.Request) {

	return
}
