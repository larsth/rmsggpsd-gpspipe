package daemon

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"runtime"
	"strings"
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
	Ed            *extractData  `json:"-"`
	isRunning     bool          `json:"-"`
	mutex         sync.Mutex    `json:"-"`
}

func (cmd *GpsPipeCmd) run() error {
	cmd.mutex.Lock()
	defer cmd.mutex.Unlock()

	if strings.Compare(runtime.GOOS, "windows") == 0 {
		return ErrGoosWindows
	}

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

func (cmd *GpsPipeCmd) init() error {
	var (
		err    error
		reader io.Reader
	)

	cmd.mutex.Lock()
	defer cmd.mutex.Unlock()

	if strings.Compare(runtime.GOOS, "windows") == 0 {
		return nil
	}

	if cmd.isRunning == false {
		cmd.Cmd = exec.Command(
			cmd.GpsJsonConfig.ExecFileName,
			cmd.GpsJsonConfig.ExecArgs...)
		if cmd.StdOutPipe, err = cmd.Cmd.StdoutPipe(); err != nil {
			return errors.Annotate(err,
				"Running cmd.Cmd.StdOutPipe() failed")
		}
		reader = io.Reader(cmd.StdOutPipe)
		if cmd.Ed, err = newExtractData(reader); err != nil {
			return errors.Annotate(err,
				"Create (*daemon.ExtractData) failed")
		}
	} else {
		return errors.Annotatef(ErrIsRunning, "%s: %s",
			"(ErrIsRunning): Cannot (*daemon.GpsPipeCmd).Init()",
			"Is already running")
	}

	return nil
}

func gpsPipePutMessage(cmd *GpsPipeCmd) (string, error) {
	var (
		extractedData *extractedData
		err           error
		buf           bytes.Buffer
	)

	if extractedData, err = cmd.Ed.ExtractData(); err != nil {
		return "", errors.Annotate(err, "Extract data error")
	}

	if extractedData.Message != nil {
		cache.Put(extractedData.Message)
	} else {
		buf.WriteString("gpsd error message@")
		buf.WriteString(extractedData.Time.Format(time.RFC3339))
		buf.WriteString(": ")
		buf.WriteString(extractedData.GpsdErrorMessage)

		return buf.String(), nil
	}

	return "", nil
}

//gpsPipe is a go routine that "speaks" with the external gpspipe program.
func gpsPipe(cmd *GpsPipeCmd) {
	const cutDuration = time.Duration(time.Millisecond * 200)
	var (
		ticker       *time.Ticker
		gpsdErr      string
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
			//A 'cmd.Config.TickerDuration.Duration' duration of time had elapsed ...
			//Read from the gpspipe external executable's standard output (STDOUT):
			if gpsdErr, err = gpsPipePutMessage(cmd); err != nil {
				annotatedErr = errors.Annotatef(err, "%s: %s",
					"go routine: daemon.gpsPipe(cmd *GpsPipeCmd)",
					"daemon.gpsPipePutMessage(*GpsPipeCmd) FATAL error")
				defer writeerror.AndExit(annotatedErr, 4)
				return //kill this go routine
			}
			if len(gpsdErr) > 0 {
				log.Printf("gpsd error: %s", gpsdErr)
			}

			runtime.Gosched() //run another go routine
		default:
			//To avoid a deadlock on single-core microprocessors ...
			runtime.Gosched() //run another go routine
		}
	}
}
