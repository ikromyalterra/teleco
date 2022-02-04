package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

type WorkerLogger struct {
	WorkerLogger zerolog.Logger
}

var WorkerLog WorkerLogger

func init() {
	WorkerLog = WorkerLogger{Logger}
}

func (log *WorkerLogger) Info(args ...interface{}) {
	log.WorkerLogger.Info().Msg(fmt.Sprintf("%v", args...))
}

func (log *WorkerLogger) Infof(format string, args ...interface{}) {
	log.WorkerLogger.Printf(format, args...)
}

func (log *WorkerLogger) Debug(args ...interface{}) {
	log.WorkerLogger.Debug().Msg(fmt.Sprintf("%v", args...))
}

func (log *WorkerLogger) Debugf(format string, args ...interface{}) {
	log.WorkerLogger.Debug().Msgf(format, args...)
}

func (log *WorkerLogger) Error(args ...interface{}) {
	log.WorkerLogger.Error().Msg(fmt.Sprintf("%v", args...))
}

func (log *WorkerLogger) Errorf(format string, args ...interface{}) {
	log.WorkerLogger.Error().Msgf(format, args...)
}
