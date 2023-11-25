package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

var Log = *log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller: true,
})
