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

type GpsPipeConfig struct {
	ExecFileName   string            `json:"exec-filename"`
	ExecArgs       []string          `json:"exec-args"`
	TickerDuration gpsdjson.Duration `json:"ticker-duration,string"`
}

type GpsPipe struct {
	mutex      sync.Mutex    `json:"-"`
	hadRunOnce bool          `json:"-"`
	Cmd        *exec.Cmd     `json:"-"`
	StdOutPipe io.ReadCloser `json:"-"`
	Config     GpsPipeConfig `json:"gpspipe"`
}

func (g *GpsPipe) init() error {
	var (
		err error
	)

	if len(g.Config.ExecArgs) == 0 {
		return errors.Annotate(err, "Zero arguments used to run the"+
			"the external gpspipe program.")
	}
	g.Cmd = exec.Command(g.Config.ExecFileName, g.Config.ExecArgs...)
	if g.StdOutPipe, err = g.Cmd.StdoutPipe(); err != nil {
		return errors.Annotate(err, "Cannot get the standard OUT pipe"+
			"used to read from the the external gpspipe program.")
	}
	if err = g.Cmd.Start(); err != nil {
		return errors.Annotate(err, "Cannot start "+
			"the external gpspipe program.")
	}
	return nil
}

func (g *GpsPipe) run() error {
	var (
		err error
	)
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if !g.hadRunOnce {
		if err = g.init(); err != nil {
			return errors.Annotate(err, "Cannot init *daemon.GpsPipeCmd")
		}
		go pipe(g)
		g.hadRunOnce = true
		return nil
	}
	return errors.Annotatef(ErrIsRunning, "%s: %s",
		"(ErrIsRunning): Cannot (*daemon.GpsPipeCmd).Run()",
		"Is already running.")
}

//pipe is a go routine that read gpsd JSON documents via
//the external gpspipe program.
func pipe(cmd *GpsPipe) {
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

	if cmd.Config.TickerDuration.D < cutDuration {
		d = cutDuration
	} else {
		d = cmd.Config.TickerDuration.D
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
