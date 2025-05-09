package logs

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Print(level zerolog.Level, args ...any) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

// Debug DebugLevel is the lowest level of logging.
// Debug logs are intended for debugging and development purposes.
func (logger *Logger) Debug(args ...any) {
	logger.Print(zerolog.DebugLevel, args...)
}

// Info InfoLevel is used for general informational log messages.
func (logger *Logger) Info(args ...any) {
	logger.Print(zerolog.InfoLevel, args...)
}

// Warn WarnLevel is used for undesired but relatively expected events, which may indicate a problem.
func (logger *Logger) Warn(args ...any) {
	logger.Print(zerolog.WarnLevel, args...)
}

// Error ErrorLevel is used for undesired and unexpected events that the program can recover from.
func (logger *Logger) Error(args ...any) {
	logger.Print(zerolog.ErrorLevel, args...)
}

func (logger *Logger) Errorf(format string, args ...any) {
	log.WithLevel(zerolog.ErrorLevel).Msgf(format, args...)
}

// Fatal FatalLevel is used for undesired and unexpected events that the program cannot recover from.
func (logger *Logger) Fatal(args ...any) {
	logger.Print(zerolog.FatalLevel, args...)
}

// Note: reserving value zero to differentiate unspecified case.
//level_unspecified LogLevel = iota

// Printf only for  internal log system
func (logger *Logger) Infof(format string, v ...any) {
	log.WithLevel(zerolog.InfoLevel).Msgf(format, v...)
}
