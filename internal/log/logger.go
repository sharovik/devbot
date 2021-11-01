package log

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var loggerInstance LoggerInstance

// Config configuration required by Logger
type Config struct {
	Env               string
	Output            string
	Level             string
	FieldContext      string
	FieldLevelName    string
	FieldErrorMessage string
}

// LoggerInstance shared state
type LoggerInstance struct {
	config      Config
	env         string
	context     map[string]interface{}
	initialized bool
}

//default constants
const (
	FieldContext      = "context"
	FieldLevelName    = "level_name"
	FieldErrorMessage = "error_message"

	appEnvTesting = "testing"

	Dbg   = "debug"
	Inf   = "info"
	Err   = "error"
	Warn  = "warn"
	Fatal = "fatal"
	Panic = "panic"

	OutputConsole = "console"
	OutputJSON    = "json"
)

func (l *LoggerInstance) setLogLevel() {
	switch l.config.Level {
	case Dbg:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		break
	case Inf:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		break
	case Err:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		break
	case Warn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		break
	case Fatal:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
		break
	case Panic:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
		break
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		break
	}
}

func (l *LoggerInstance) getOutput() zerolog.Context {
	switch l.config.Output {
	case OutputConsole:
		return log.Output(zerolog.NewConsoleWriter()).With()
	case OutputJSON:
		return log.With()
	}

	return log.With()
}

// Init initializes the logger (required before use)
func Init(config Config) error {
	loggerInstance = LoggerInstance{
		config:      config,
		env:         config.Env,
		context:     make(map[string]interface{}),
		initialized: true,
	}

	return nil
}

//Refresh refreshes the logger instance
func Refresh() {
	loggerInstance = LoggerInstance{}
}

//IsInitialized function retrieves current status of logger instance
func IsInitialized() bool {
	return loggerInstance.initialized
}

//Logger returns a pointer to the singleton Logger loggerInstance
func Logger() *LoggerInstance {
	if !loggerInstance.initialized {
		panic("logger not initialized")
	}
	return &loggerInstance
}

//AppendGlobalContext for setting global context
func (l *LoggerInstance) AppendGlobalContext(context map[string]interface{}) {
	if l.context == nil {
		l.context = context
	}

	for field, value := range context {
		l.context[field] = value
	}

	l.Debug().Interface("context_changes", context).Msg("Append new global context")
}

//GlobalContext method retrieve the GlobalContext variable
func (l *LoggerInstance) GlobalContext() map[string]interface{} {
	return l.context
}

//DestroyGlobalContext method for global context destroy
func (l *LoggerInstance) DestroyGlobalContext() {
	l.context = make(map[string]interface{})
}

//AddError for correct error messages parse
func (l *LoggerInstance) AddError(err error) *zerolog.Event {
	err = errors.Wrap(err, err.Error())
	return l.Error().Stack().Err(err)
}

//DefaultContext method which returns Logger with default context
func (l *LoggerInstance) DefaultContext() *zerolog.Logger {
	l.setLogLevel()
	var context = zerolog.Context{}
	switch l.config.Env {
	case appEnvTesting:
		//For testing environment we need to disable the logs
		context = log.Output(ioutil.Discard).With()
		break
	default:
		context = l.getOutput()
	}

	zerolog.TimestampFieldName = "@timestamp"
	zerolog.LevelFieldName = l.config.FieldLevelName
	zerolog.ErrorFieldName = l.config.FieldErrorMessage
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000000"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := context.
		Interface(l.config.FieldContext, l.context).
		Logger()

	return &logger
}

//Debug method for messages with level DEBUG
func (l *LoggerInstance) Debug() *zerolog.Event {
	return l.DefaultContext().Debug()
}

//Info method for messages with level INFO
func (l *LoggerInstance) Info() *zerolog.Event {
	return l.DefaultContext().Info()
}

//Error method for messages with level ERROR
func (l *LoggerInstance) Error() *zerolog.Event {
	return l.DefaultContext().Error()
}

//Warn method for messages with level WARNING
func (l *LoggerInstance) Warn() *zerolog.Event {
	return l.DefaultContext().Warn()
}

//StartMessage adds message with START postfix
func (l *LoggerInstance) StartMessage(msg string) {
	l.DefaultContext().Info().Msg(fmt.Sprintf("%s: %s", msg, "START"))
}

//FinishMessage adds message with FINISH postfix
func (l *LoggerInstance) FinishMessage(msg string) {
	l.DefaultContext().Info().Msg(fmt.Sprintf("%s: %s", msg, "FINISH"))
}
