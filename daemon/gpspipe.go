package daemon

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/larsth/go-rmsggpsbinmsg"
	"github.com/larsth/linescanner"
	"github.com/larsth/rmsggpsd-gpspipe/errors"
	"github.com/larsth/writeerror"
)

type GpsPipe struct {
	mutex       sync.Mutex    `json:"-"`
	hadRunOnce  bool          `json:"-"`
	cmd         *exec.Cmd     `json:"-"`
	stdOutPipe  io.ReadCloser `json:"-"`
	gpsdFilter  filter        `json:"-"`
	lineScanner *linescanner.LineScanner
	logger      *log.Logger
	Config      GpsPipeConfig `json:"gpspipe"`
}

func (g *GpsPipe) init() error {
	var (
		err error
	)

	if len(g.Config.ExecArgs) == 0 {
		return errors.Annotate(err, `ERROR: `+
			`Zero arguments used to run the external gpspipe program.`)
	}
	g.cmd = exec.Command(g.Config.ExecFileName, g.Config.ExecArgs...)
	if g.stdOutPipe, err = g.cmd.StdoutPipe(); err != nil {
		return errors.Annotate(err, "Cannot get the standard OUT pipe"+
			"used to read from the the external gpspipe program.")
	}

	g.gpsdFilter = newGpsdFilter(g.logger)
	if g.lineScanner, err = linescanner.New(g.stdOutPipe); err != nil {
		return errors.Annotate(err, "Error while initializing the linescanner")
	}

	return nil
}

func (g *GpsPipe) Run(logger *log.Logger) error {
	var (
		err error
	)
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if !g.hadRunOnce {
		if logger == nil {
			return errors.New("*log.Logger is nil")
		}
		g.logger = logger

		if err = g.init(); err != nil {
			return errors.Annotate(err, "Cannot init *daemon.GpsPipeCmd")
		}
		go gpspipeGoRoutine(g)
		g.hadRunOnce = true
		return nil
	}
	return errors.Annotatef(ErrIsRunning, "%s: %s",
		"(ErrIsRunning): Cannot (*daemon.GpsPipeCmd).Run()",
		"Is already running.")
}

func goroutineFATAL(err error, msg string) {
	annotatedErr := errors.Annotate(err, msg)
	s := errors.Details(annotatedErr)
	fmt.Fprintln(os.Stderr, s)
	writeerror.AndExit(annotatedErr, 3)
	//does not return ...
}

func pipeInit(g *GpsPipe, d *time.Duration) error {
	const minDuration = time.Duration(time.Millisecond * 100)
	var err error

	if err = g.cmd.Start(); err != nil {
		return errors.Annotate(err, "Cannot start "+
			"the external gpspipe program.")
	}

	if g.Config.TickerDuration.D < minDuration {
		*d = minDuration
	} else {
		*d = g.Config.TickerDuration.D
	}
	return nil
}

func (g *GpsPipe) readJson() (p []byte, err error) {
	const maxLoops = 128

	for g.lineScanner.Scan() {
		if g.lineScanner.ReadCount() > maxLoops {
			return nil, errors.Errorf(
				"line scanner, maximum of loops exceeded: %d loops",
				maxLoops)
		}
	}
	if g.lineScanner.Err() != nil {
		return nil, errors.Annotate(g.lineScanner.Err(), "linescanner error")
	}
	return g.lineScanner.Bytes(), nil
}

func (g *GpsPipe) pipeAction(ticker *time.Ticker) error {
	var (
		p   []byte
		err error
		m   *binmsg.Message
	)

	select {
	case _ = <-ticker.C:
		if p, err = g.readJson(); err != nil {
			return errors.Annotate(err, `Cannot read a gpsd JSON document`)
		}
		if m, err = g.gpsdFilter.ParseGpsdJson(p); err != nil {
			return errors.Annotate(err, `Cannot parse a gpsd JSON document`)
		}
		if m != nil {
			//save this GPS coordinate to the cache
			thisGpsCache.Put(m)
		}
	}
	return nil
}

//pipe is a go routine that read gpsd JSON documents via
//the external gpspipe program.
func gpspipeGoRoutine(g *GpsPipe) {
	var (
		err    error
		ticker *time.Ticker
		d      time.Duration
	)

	if err = pipeInit(g, &d); err != nil {
		goroutineFATAL(err, "FATAL ERROR, gpspipe go routine: Cannot init")
		//pipeFATAL does not return, but make it clear that this go routine is killed
		return
	}

	ticker = time.NewTicker(d)

	//loop for infinity, or until an error occurs ...
	for {
		if err = g.pipeAction(ticker); err != nil {
			msg := "FATAL ERROR, gpspipe go routine, pipe action error"
			goroutineFATAL(err, msg)
			//pipeFATAL does not return, but make it clear that this go routine is killed
			return
		}
	}
}
