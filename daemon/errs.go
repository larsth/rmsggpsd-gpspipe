package daemon

import "github.com/larsth/rmsggpsd-gpspipe/errors"
import "github.com/larsth/rmsggpsd-gpspipe/env"

var (
	ErrNilFilter    = errors.New("The *gpsdfilter.Filer is a nil pointer")
	ErrIsRunning    = errors.New("The go routine is already running")
	ErrNoEnvVarName = errors.New("The enviroment variable `" +
		env.ConfigVarName +
		"` is not defined, or has the empty string as content.")
	ErrGoosWindows = errors.New("Running on GOOS=windows: " +
		"External application \"gpsd\" is required, and it " +
		"is a UNIX/Linux/BSD-only application " +
		"that cannot run on a Windows operating system.")
)
