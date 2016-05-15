package daemon

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/larsth/go-gpsfix"
	"github.com/larsth/rmsggpsd-gpspipe/errors"
)

type (
	GpsCoord struct {
		Alt     float32        `json:"altitude"`
		Lat     float32        `json:"latitude"`
		Lon     float32        `json:"longitude"`
		FixMode gpsfix.FixMode `json:"fixmode,string"`
	}

	Application struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	JsonConfig struct {
		Application Application `json:"application"`
		HttpAddr    string      `json:"httpd-addr"`
		GpsPipeCmd  *GpsPipe `json:"gpspipe"`
		Gps         *GpsCoord   `json:"gps-coord"`
	}
)

const (
	envVarName          = `RMSGGPSD_JSONCONF`
	expectedJsonVersion = `1.0.0`
)

func checkIsARmsggpsdJsonDocument(a *Application) error {
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

	if err = checkIsARmsggpsdJsonDocument(&c.Application); err != nil {
		s := `{
	"application":{
		"name":"rmsggpsd",
		"version":"1.0.0",
	},
	"gpspipe":{},
	"httpd-addr":"An IPv6 address goes here"
}`
		return errors.Annotatef(err, "%s\nWant something like: %s\nGot:%s\n",
			`The JSON configuration file is a \"rmsggpsd\" JSON document.`,
			s,
			string(p))
	}

	return nil
}
