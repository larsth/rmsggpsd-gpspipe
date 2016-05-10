package daemon

import (
	"log"

	"github.com/larsth/go-gpsdfilter"
	"github.com/larsth/rmsggpsd/errors"
)

var filter *gpsdfilter.Filter

var rmsggps1dFilterRules = []*gpsdfilter.Rule{
	&gpsdfilter.Rule{
		Class: "ERROR",
		DoLog: true,
		Type:  gpsdfilter.TypeParse,
	},
	&gpsdfilter.Rule{
		Class: "TPV",
		DoLog: true,
		Type:  gpsdfilter.TypeParse,
	},
	&gpsdfilter.Rule{
		Class: "DEVICE",
		DoLog: true,
		Type:  gpsdfilter.TypeLog,
	},
	&gpsdfilter.Rule{
		Class: "DEVICES",
		DoLog: true,
		Type:  gpsdfilter.TypeLog,
	},
	&gpsdfilter.Rule{
		Class: "POLL",
		DoLog: true,
		Type:  gpsdfilter.TypeLog,
	},
	&gpsdfilter.Rule{
		Class: "SKY",
		DoLog: true,
		Type:  gpsdfilter.TypeLog,
	},
	&gpsdfilter.Rule{
		Class: "VERSION",
		DoLog: true,
		Type:  gpsdfilter.TypeLog,
	},
	&gpsdfilter.Rule{
		Class: "WATCH",
		DoLog: true,
		Type:  gpsdfilter.TypeLog,
	},
	&gpsdfilter.Rule{
		Class: "ATT",
		DoLog: false,
		Type:  gpsdfilter.TypeIgnore,
	},
	&gpsdfilter.Rule{
		Class: "GST",
		DoLog: false,
		Type:  gpsdfilter.TypeIgnore,
	},
	&gpsdfilter.Rule{
		Class: "PPS",
		DoLog: false,
		Type:  gpsdfilter.TypeIgnore,
	},
	&gpsdfilter.Rule{
		Class: "TOFF",
		DoLog: false,
		Type:  gpsdfilter.TypeIgnore,
	},
}

func addFilterRules(f *gpsdfilter.Filter) error {
	var err, annotatedErr error

	if f == nil {
		return ErrNilFilter
	}

	for _, rule := range rmsggps1dFilterRules {
		if err = f.Add(rule); err != nil {
			annotatedErr = errors.Annotate(err,
				"(*gpsdfilter.Filter).AddRule error")
			return annotatedErr
		}
	}

	return nil
}

func initFilterRules() {
	filter = gpsdfilter.New()
	if err := addFilterRules(filter); err != nil {
		log.Fatalln(err.Error())
	}
}
