package main

import (
	"log"
	"os"

	"github.com/larsth/rmsggpsd/daemon"
	"github.com/larsth/rmsggpsd/errors"
)

var c daemon.JsonConfig

func main() {
	var (
		err          error
		annotatedErr error
		s            string
	)

	log.SetPrefix("rmsggpsd ")
	log.SetFlags(log.LUTC)
	errors.RemoveFilePath = true

	if err = daemon.Start(&c); err != nil {
		annotatedErr = errors.Annotate(err, "FATAL ERROR")
		s = errors.ErrorStack(annotatedErr)
		log.Println(s)
		os.Exit(-1)
	}
	os.Exit(0)
}
