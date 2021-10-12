package server

import (
	"emviwiki/shared/config"
	"fmt"
	"github.com/emvi/logbuch"
	"log"
	"path/filepath"
	"strings"
	"time"
)

const (
	loggingFilenameTimeFormat = "2006-01-02_15:04:05"
	loggingDirnameTimeFormat  = "2006-01-02_15:04:05"
	logFiles                  = 20
	logFileSize               = 1024 * 1024 * 5
	logBufferSize             = 1024 * 4
)

type nameSchema struct {
	filename string
	counter  int
}

func (schema *nameSchema) Name() string {
	schema.counter++

	if schema.counter > 999 {
		schema.counter = 1
	}

	return fmt.Sprintf("%s_%03d_%s.log", time.Now().Format(loggingFilenameTimeFormat), schema.counter, schema.filename)
}

// ConfigureLogging configures the log level by an environment variable.
func ConfigureLogging() (*logbuch.RollingFileAppender, *logbuch.RollingFileAppender) {
	logbuch.Info("Configuring logging...")
	logbuch.SetFormatter(logbuch.NewFieldFormatter(config.Get().Logging.TimeFormat, "\t"))
	level := strings.ToLower(config.Get().Logging.Level)

	if level == "debug" {
		logbuch.SetLevel(logbuch.LevelDebug)
	} else if level == "info" {
		logbuch.SetLevel(logbuch.LevelInfo)
	} else {
		logbuch.SetLevel(logbuch.LevelWarning)
	}

	dir := config.Get().Logging.Dir

	if dir != "" {
		dir = filepath.Join(dir, time.Now().Format(loggingDirnameTimeFormat))
		logbuch.Info("Logging to file, see output directory for logs...", logbuch.Fields{"dir": dir})
		stdName := &nameSchema{filename: "stdout"}
		errName := &nameSchema{filename: "stderr"}
		stdout, err := logbuch.NewRollingFileAppender(logFiles, logFileSize, logBufferSize, dir, stdName)

		if err != nil {
			logbuch.Fatal("Error configuring stdout rolling file appender")
		}

		stderr, err := logbuch.NewRollingFileAppender(logFiles, logFileSize, logBufferSize, dir, errName)

		if err != nil {
			logbuch.Fatal("Error configuring stderr rolling file appender")
		}

		logbuch.SetOutput(stdout, stderr)
		return stdout, stderr
	}

	return nil, nil
}

// CloseLogger closes the rolling file appenders in case they are open (non nil).
func CloseLogger(stdout *logbuch.RollingFileAppender, stderr *logbuch.RollingFileAppender) {
	if stdout != nil {
		log.Println("Closing stdout rolling file appender")

		if err := stdout.Close(); err != nil {
			log.Println("Error closing stdout rolling file appender")
		}
	}

	if stderr != nil {
		log.Println("Closing stderr rolling file appender")

		if err := stderr.Close(); err != nil {
			log.Println("Error closing stderr rolling file appender")
		}
	}
}
