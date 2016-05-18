package daemon

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/larsth/go-gpsdjson"
	"github.com/larsth/go-gpsfix"
	"github.com/larsth/rmsggpsd-gpspipe/errors"
)

type (
	ApplicationConfig struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	AddrLoggerConfig struct {
		AddrString     string            `json:"addr"`
		LoggerString   string            `json:"logger"`
		TickerDuration gpsdjson.Duration `json:"ticker-duration"`
	}

	GpsCoordConfig struct {
		Alt     float32        `json:"altitude"`
		Lat     float32        `json:"latitude"`
		Lon     float32        `json:"longitude"`
		FixMode gpsfix.FixMode `json:"fixmode,string"`
	}

	GpsPipeConfig struct {
		ExecFileName   string            `json:"exec-filename"`
		ExecArgs       []string          `json:"exec-args"`
		TickerDuration gpsdjson.Duration `json:"ticker-duration,string"`
		Logger         string            `json:"logger"`
	}

	JsonConfig struct {
		Application ApplicationConfig `json:"application"`
		Httpd       AddrLoggerConfig  `json:"httpd"`
		OtherGps    AddrLoggerConfig  `json:"other-gps"`
		GpsPipe     *GpsPipe          `json:"gpspipe"`
		ThisGps     *GpsCoordConfig   `json:"this-gps"`
	}
)

const (
	envVarName          = `RMSGGPSD_JSONCONF`
	expectedJsonVersion = `1.0.0`
)

func checkIsRmsggpsdJsonDocument(a *ApplicationConfig) error {
	if strings.Compare("rmsggpsd", a.Name) != 0 {
		return errors.Errorf("%s./nWant: \"rmsggpsd\"\nGot: \"%s\"",
			`Invalid application name in the JSON document: `,
			a.Name)
	}

	if strings.Compare(expectedJsonVersion, a.Version) != 0 {
		return errors.Errorf("%s./nWant: \"%s\"\nGot: \"%s\"",
			`Invalid application version in the JSON document: `,
			expectedJsonVersion,
			a.Version)
	}

	return nil
}

func readJsonConfigDocument(c *JsonConfig) error {
	var (
		fileName string
		ok       bool
		err      error
		file     *os.File
		buf      bytes.Buffer
		p        []byte
	)

	if fileName, ok = os.LookupEnv(envVarName); !ok {
		return errors.Errorf("%s '%s' %s.",
			"The enviroment variable", envVarName,
			"does not exists")
	}

	if len(fileName) == 0 {
		return errors.Errorf("%s '%s' %s.",
			"The enviroment variable", envVarName,
			"has no content")
	}

	if file, err = os.Open(fileName); err != nil {
		return errors.Annotate(err,
			`Cannot open the JSON configuration file.`)
	}
	defer file.Close()

	if _, err = io.Copy(&buf, file); err != nil {
		return errors.Annotate(err,
			`Cannot read (copy into memory) the content of the JSON configuration file.`)
	}
	p = append(p, buf.Bytes()...)

	if err = json.Unmarshal(p, c); err != nil {
		return errors.Annotate(err,
			`Cannot unmarshal the JSON configuration file.`)
	}

	if err = checkIsRmsggpsdJsonDocument(&c.Application); err != nil {
		return errors.Trace(err)
	}

	return nil
}
