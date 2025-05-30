package logger

import (
	"github.com/charmbracelet/log"
	"os"
	"time"
)

func New(prefix string) *log.Logger {
	return log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          prefix,
	})
}
