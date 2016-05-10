package daemon

import (
	"net/http"

	"github.com/larsth/go-rmsggpsbinmsg"
)

func binMsgWebRequestHandler(
	w http.ResponseWriter,
	req *http.Request,
	binMessage *binmsg.Message) {

	var ok bool

	w.Header().Set("Cache-Control", "no-cache")

	if ok = isGETHttpMethod(req, w); !ok {
		return
	}

	writeDateResponseHeaders(w, binMessage)

	if ok = parseForm(req, w); !ok {
		return
	}

	if _, ok = req.Form["json"]; ok {
		writeJSONHttpResponse(w, binMessage)
		return
	}

	if _, ok = req.Form["binmsg"]; ok {
		writeBinMsgHttpResponse(w, binMessage)
		return
	}

	if _, ok = req.Form["csv"]; ok {
		writeCSVHttpResponse(w, binMessage)
		return
	}

	//by default: return 'application/x-www-form-urlencoded' data
	writeXWwwFormUrlencodedHttpResponse(w, binMessage)
	return
}

func rootWebRequestHandler(w http.ResponseWriter, req *http.Request) {
	var (
		ok         bool
		binMessage *binmsg.Message
	)

	if binMessage, ok = getChachedBinMsg(w); !ok {
		return
	}

	binMsgWebRequestHandler(w, req, binMessage)
	return
}
