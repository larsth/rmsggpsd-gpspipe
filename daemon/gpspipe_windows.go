//+build ignore
//Ignore, because build tags are broken in Go 1.6.3

// +build !linux !darwin !freebsd
package daemon

import "github.com/larsth/go-gpsdjson"

type GpsJsonConfig struct {
	ExecFileName   string            `json:"exec-filename"`
	ExecArgs       []string          `json:"exec-args"`
	TickerDuration gpsdjson.Duration `json:"ticker-duration,string"`
}

type GpsPipeCmd struct {
	GpsJsonConfig GpsJsonConfig `json:"gpspipe"`
}

func (_ *GpsPipeCmd) run() error {
	return ErrGoosWindows
}

func (_ *GpsPipeCmd) init() error {
	return nil
}
