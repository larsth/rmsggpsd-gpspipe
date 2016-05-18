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
		err      error
		fixmodes string
		times    string
		alt      float32
		lat      float32
		lon      float32
		fixmode  gpsfix.FixMode
		t        time.Time
		urlValue string
	)

	if alt, urlValue, err = parseFloat32FromUrlValue(v, "alt"); err != nil {
		return nil, errors.Annotate(err,
			"Cannot parse altitude string info floating point 32 bit value")
	}
	if lat, urlValue, err = parseFloat32FromUrlValue(v, "lat"); err != nil {
		return nil, errors.Annotatef(err, "%s :'%s'",
			`Cannot parse latitide string`,
			urlValue,
			`info a floating point 32 bit value`)
	}
	if lon, urlValue, err = parseFloat32FromUrlValue(v, "lon"); err != nil {
		return nil, errors.Annotate(err,
			"Cannot parse longitude string info floating point 32 bit value")
	}
	if fixmodes = v.Get("fixmode"); len(fixmodes) == 0 {
		return nil, errors.New(
			`The value from the "fixmode" key does not exist, or is empty`)
	}
	if fixmode, err = gpsfix.Parse(fixmodes); err != nil {
		return nil, errors.Annotate(err, "Cannot parse fixmode string")
	}
	if times = v.Get("time"); len(times) == 0 {
		return nil, errors.New(
			`The value from the "time" key does not exist, or is empty`)
	}
	if t, err = time.Parse(time.RFC3339, times); err != nil {
		return nil, errors.Annotatef(err, "%s %s",
			`Could not parse RFC3339 ("2006-01-02T15:04:05Z07:00")`,
			`encoded time string`)
	}

	return mkBinMsg(alt, lat, lon, fixmode, t), nil
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
		return errors.Annotatef(err, "%s", "TBD")
	}
	defer r.Body.Close()

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return errors.Annotatef(err, "%s", "TBD")
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
			`Cannot make a binmsg.Message with keys from a url values`)
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
