package utils

import (
	"github.com/charmbracelet/log"
	"os"
	"time"
)

func NewLogger(prefix string) *log.Logger {
	return log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          prefix,
	})
}
