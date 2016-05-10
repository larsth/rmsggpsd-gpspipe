package daemon

import (
	"bytes"
	"encoding/json"
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

func getChachedBinMsg(w http.ResponseWriter) (binMessage *binmsg.Message, ok bool) {
	var (
		err error
		msg string
	)

	if binMessage, err = cache.Get(); err != nil {
		msg = fmt.Sprintf(
			"Cannot get binary message from the cache, Error %s.",
			err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println(msg)

		return nil, false
	}
	return binMessage, true
}

//write Date* HTTP headers
func writeDateResponseHeaders(w http.ResponseWriter, binMessage *binmsg.Message) {
	var nowUTC time.Time = time.Now().UTC()

	w.Header().Set("Date", nowUTC.Format(rfc7231))
	w.Header().Set("Date-RFC-3339", nowUTC.Format(time.RFC3339))
	w.Header().Set("Date-RFC3339-Nano", nowUTC.Format(time.RFC3339Nano))
	w.Header().Set(
		"Date-Bin-Msg-RFC3339",
		binMessage.TimeStamp.Time.Format(time.RFC3339))
	w.Header().Set(
		"Date-Bin-Msg-RFC3339Nano",
		binMessage.TimeStamp.Time.Format(time.RFC3339Nano))
}

func parseForm(req *http.Request, w http.ResponseWriter) (ok bool) {
	var (
		err error
		msg string
	)
	if err = req.ParseForm(); err != nil {
		msg = fmt.Sprintf("Cannot parse HTTP GET query, Error %s.", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println(msg)
		return false
	}
	return true
}

func writeJSONHttpResponse(w http.ResponseWriter, binMessage *binmsg.Message) {
	var (
		p    []byte
		err  error
		msg  string
		pLen string
	)

	w.Header().Set("Content-Type", "application/json")
	if p, err = json.Marshal(binMessage); err != nil {
		msg = `Cannot parse the binary message into a JSON GET query, Error: `
		msg = fmt.Sprintf("%s '%s'.\n", msg, err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println(msg)
		return
	}

	pLen = strconv.Itoa(len(p))
	w.Header().Set("Content-Length", pLen)
	if len(p) > 0 {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(p)
	}

	return
}

func writeBinMsgHttpResponse(w http.ResponseWriter, binMessage *binmsg.Message) {
	var (
		p    []byte
		err  error
		msg  string
		pLen string
	)

	//Tell a User-Agent (browser) to not to display
	//the body, but save the body to a file.
	w.Header().Set("Content-Disposition", "attachment")

	if p, err = binMessage.MarshalBinary(); err != nil {
		msg = fmt.Sprintf(
			"Cannot marshal *binmsg.Message to a binary representation, Error %s.",
			err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println(msg)
		return
	}

	w.Header().Set("Content-Type", "application/x.rmsgdk.binmsg")
	w.Header().Set("Content-Transfer-Encoding", "binary")

	pLen = strconv.Itoa(len(p))
	w.Header().Set("Content-Length", pLen)
	if len(p) > 0 {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(p)
	}
	return
}

func writeCSVHttpResponse(w http.ResponseWriter, binMessage *binmsg.Message) {
	var (
		fixmode, alt, lat, lon, timestamp = binMessage.Strings()
		buf                               bytes.Buffer
		p                                 []byte
		pLen                              string
	)
	//Construct message
	buf.WriteString("fixmode")
	buf.WriteString(":")
	buf.WriteString(fixmode)
	buf.WriteString("\n")

	buf.WriteString("alt")
	buf.WriteString(":")
	buf.WriteString(alt)
	buf.WriteString("\n")

	buf.WriteString("lat")
	buf.WriteString(":")
	buf.WriteString(lat)
	buf.WriteString("\n")

	buf.WriteString("lon")
	buf.WriteString(":")
	buf.WriteString(lon)
	buf.WriteString("\n")

	buf.WriteString("timestamp_RFC3339")
	buf.WriteString(":")
	buf.WriteString(timestamp)
	buf.WriteString("\n")

	w.Header().Set("Content-Type", "text/csv")

	p = buf.Bytes()
	pLen = strconv.Itoa(len(p))
	w.Header().Set("Content-Length", pLen)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(p)
	return
}

func writeXWwwFormUrlencodedHttpResponse(
	w http.ResponseWriter,
	binMessage *binmsg.Message) {

	var (
		fixmode, alt, lat, lon, timestamp = binMessage.Strings()
		values                            url.Values
		p                                 []byte
		pLen                              string
	)

	values.Set("fixmode", fixmode)
	values.Set("alt", alt)
	values.Set("lat", lat)
	values.Set("lon", lon)
	values.Set("timestamp_rfc3339", timestamp)

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Cache-Control", "no-cache")

	p = []byte(values.Encode())
	pLen = strconv.Itoa(len(p))
	w.Header().Set("Content-Length", pLen)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(p)

	return
}
