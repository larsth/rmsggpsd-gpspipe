package daemon

import (
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/larsth/go-gpsdjson"
	"github.com/larsth/rmsggpsd-gpspipe/errors"
	"github.com/larsth/writeerror"
)

type GpsJsonConfig struct {
	ExecFileName   string            `json:"exec-filename"`
	ExecArgs       []string          `json:"exec-args"`
	TickerDuration gpsdjson.Duration `json:"ticker-duration,string"`
}

type GpsPipeCmd struct {
	GpsJsonConfig GpsJsonConfig `json:"gpspipe"`
	Cmd           *exec.Cmd     `json:"-"`
	StdOutPipe    io.ReadCloser `json:"-"`
	isRunning     bool          `json:"-"`
	mutex         sync.Mutex    `json:"-"`
}

func (cmd *GpsPipeCmd) run() error {
	cmd.mutex.Lock()
	defer cmd.mutex.Unlock()

	if cmd.isRunning == false {
		go gpsPipe(cmd)
		cmd.isRunning = true
		return nil
	} else {
		return errors.Annotatef(ErrIsRunning, "%s: %s",
			"(ErrIsRunning): Cannot (*daemon.GpsPipeCmd).Run()",
			"Is already running")
	}
}

//gpsPipe is a go routine that "speaks" with the external gpspipe program.
func gpsPipe(cmd *GpsPipeCmd) {
	const cutDuration = time.Duration(time.Millisecond * 200)
	var (
		ticker       *time.Ticker
		err          error
		annotatedErr error
		d            time.Duration
	)

	if err = cmd.Cmd.Start(); err != nil {
		annotatedErr = errors.Annotatef(err, "%s: %s",
			"FATAL ERROR, gpspipe go routine",
			"Cannot start external command: 'gpspipe'")

		writeerror.AndExit(annotatedErr, 3)
		return //kill this go routine
	}

	if cmd.GpsJsonConfig.TickerDuration.D < cutDuration {
		d = cutDuration
	} else {
		d = cmd.GpsJsonConfig.TickerDuration.D
	}
	ticker = time.NewTicker(d)

	//loop for infinity, or until an error occurs ...
	for {
		select {
		case _ = <-ticker.C:
			//			//A 'cmd.Config.TickerDuration.Duration' duration of time had elapsed ...
			//			//Read from the gpspipe external executable's standard output (STDOUT):
			//			if gpsdErr, err = gpsPipePutMessage(cmd); err != nil {
			//				annotatedErr = errors.Annotatef(err, "%s: %s",
			//					"go routine: daemon.gpsPipe(cmd *GpsPipeCmd)",
			//					"daemon.gpsPipePutMessage(*GpsPipeCmd) FATAL error")
			//				defer writeerror.AndExit(annotatedErr, 4)
			//				return //kill this go routine
			//			}
			//			if len(gpsdErr) > 0 {
			//				log.Printf("gpsd error: %s", gpsdErr)
			//			}
		}
	}
}
