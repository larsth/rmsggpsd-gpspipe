package daemon

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/larsth/go-gpsfix"
	"github.com/larsth/go-rmsggpsbinmsg"
	"github.com/larsth/rmsggpsd-gpspipe/errors"
)

var ErrNoSuchUrlKey = errors.New("No such url.Values key")

func parseFloat32FromUrlValue(v url.Values, key string) (float32, string, error) {
	var (
		NaN float32 = float32(math.NaN())
		s   string
		err error
		f   float64
	)
	if s = v.Get(key); len(s) == 0 {
		return NaN, s, errors.Annotatef(ErrNoSuchUrlKey,
			"url.Values[%s] not found", key)
	}
	if f, err = strconv.ParseFloat(s, 32); err != nil {
		return NaN, s, errors.Annotatef(err,
			"cannot parse url value: '%s' from key: '%s' %s",
			s, key, `to a floating point 32-bit value`)
	}
	return float32(f), "", nil
}

func makeBinMsgFromUrlValues(v url.Values) (*binmsg.Message, error) {
	var (
		err        error
		alt        float32
		lat        float32
		lon        float32
		gpsfixmode string
		fixmode    gpsfix.FixMode
		gpstime    string
		t          time.Time
		value      string
	)

	alt, value, err = parseFloat32FromUrlValue(v, "gpsaltitude")
	if err != nil {
		return nil, errors.Annotatef(err, "%s %s %s",
			`Cannot parse altitude string`,
			value,
			`into a floating point 32 bit value`)
	}
	lat, value, err = parseFloat32FromUrlValue(v, "gpslatitude")
	if err != nil {
		return nil, errors.Annotatef(err, "%s %s %s",
			`Cannot parse latitide string`,
			value,
			`into a floating point 32 bit value`)
	}
	lon, value, err = parseFloat32FromUrlValue(v, "gpslongitude")
	if err != nil {
		return nil, errors.Annotatef(err, "%s %s %s",
			`Cannot parse longitude string`,
			value,
			`into a floating point 32 bit value`)
	}
	gpsfixmode = v.Get("gpsfixmode")
	if len(gpsfixmode) == 0 {
		return nil, errors.New(
			`The value from the "fixmode" key does not exist, or is empty`)
	}
	fixmode, err = gpsfix.Parse(gpsfixmode)
	if err != nil {
		return nil, errors.Annotate(err, "Cannot parse fixmode string")
	}
	gpstime = v.Get("gpstime")
	if len(gpstime) == 0 {
		return nil, errors.New(
			`The value from the "time" key does not exist, or is empty`)
	}
	t, err = time.Parse(time.RFC3339, gpstime)
	if err != nil {
		return nil, errors.Annotatef(err, "%s %s",
			`Could not parse RFC3339 ("2006-01-02T15:04:05Z07:00")`,
			`encoded time string`)
	}

	return binmsg.MkBinMsg(alt, lat, lon, fixmode, t), nil
}

func handleHttpClientResponse(addr string) error {
	var (
		r      *http.Response
		err    error
		body   []byte
		sbody  string
		values url.Values
		m      *binmsg.Message
	)

	if r, err = http.DefaultClient.Get(addr); err != nil {
		return errors.Annotatef(err, "%s: %s",
			`Cannot HTTP GET from web server`, addr)
	}
	defer r.Body.Close()

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return errors.Annotatef(err, "%s: %s",
			`Cannot read the HTTP body (NOT HTML!) fetched from web server`,
			addr)
	}
	sbody = string(body)

	if r.StatusCode != http.StatusOK {
		return errors.Errorf("%s: %s.\n\tStatus Code: %d\n\tHTTP body: %s\n",
			`Rmsggpsd HTTP client`,
			`Recieved something else than status code "200 OK"`,
			r.StatusCode,
			sbody)
	}
	if values, err = url.ParseQuery(sbody); err != nil {
		return errors.Errorf("%s: %s: %s.\n\tHTTP body: %s\n",
			`Rmsggpsd HTTP client: `,
			`Could not parse Content-Type (MIME encoding: `,
			`'application/x-www-form-urlencoded')`,
			sbody)
	}
	if m, err = makeBinMsgFromUrlValues(values); err != nil {
		return errors.Annotate(err,
			`Cannot make a binmsg.Message with keys from "+ 
            "the associative array of url values`)
	}
	otherGpsCache.Put(m)
	return nil
}

func httpClient(c *JsonConfig, logger *log.Logger) {
	const minDuration = time.Millisecond * 5
	var (
		addr           string
		tickerDuration time.Duration
		ticker         *time.Ticker
	)

	// c, and logger must be non nil pointers
	if logger == nil {
		s := "Other gps http client: logger of type *log.logger is nil."
		fmt.Fprintln(os.Stderr, s)
		os.Exit(3)
		return
	}
	if c == nil {
		s := "Other gps http client. *daemon.JsonConfig is nil"
		logger.Println(s)
		fmt.Fprintln(os.Stderr, s)
		os.Exit(4)
		return
	}

	addr = c.OtherGps.AddrString

	tickerDuration = c.OtherGps.TickerDuration.D
	if tickerDuration.Nanoseconds() < minDuration.Nanoseconds() {
		tickerDuration = minDuration
	}
	ticker = time.NewTicker(tickerDuration)

	for {
		select {
		case _ = <-ticker.C:
			if err := handleHttpClientResponse(addr); err != nil {
				details := errors.Details(err)
				logger.Println(details)
			}
		}
	}
}
