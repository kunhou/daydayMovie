package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/kunhou/TMDB/config"
	"github.com/sirupsen/logrus"
)

var cfg = config.GetConfig()

type Fields = logrus.Fields

type mlyticsFormatter struct {
	logrus.JSONFormatter
}

func (mf mlyticsFormatter) Format(e *logrus.Entry) ([]byte, error) {
	mf.JSONFormatter.TimestampFormat = "2006/01/02 15:04:05"
	e.Time = e.Time.UTC()
	return mf.JSONFormatter.Format(e)
}

func init() {
	logrus.SetFormatter(&mlyticsFormatter{})
	lv, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(lv)
	hook := newContextHook()
	logrus.AddHook(hook)
}

func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

func WithFields(fields Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Warning(args ...interface{}) {
	logrus.Warning(args)
}

func GetLevel() string {
	return logrus.GetLevel().String()
}

type contextHook struct {
	Field  string
	Skip   int
	levels []logrus.Level
}

func newContextHook(levels ...logrus.Level) logrus.Hook {
	hook := contextHook{
		Field:  "source",
		Skip:   5,
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels implement levels
func (hook contextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implement fire
func (hook contextHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}
