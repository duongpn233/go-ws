package trace

import "github.com/getsentry/sentry-go"

type SentryConfig struct {
	DSN           string
	Debug         bool
	EnableTracing bool
	Environment   string
}

func InitSentry(cfg *SentryConfig) error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:           cfg.DSN,
		Debug:         cfg.Debug,
		EnableTracing: cfg.EnableTracing,
		Environment:   cfg.Environment,
	})
}

func CaptureException(err error) {
	sentry.CaptureException(err)
}

func CaptureMessage(msg string) {
	sentry.CaptureMessage(msg)
}

func CaptureEvent(lv sentry.Level, msg string, log map[string]interface{}) {
	event := sentry.Event{
		Level:   lv,
		Message: msg,
		Extra:   log,
	}
	sentry.CaptureEvent(&event)
}

const (
	ERROR   = sentry.LevelError
	DEBUG   = sentry.LevelDebug
	INFO    = sentry.LevelInfo
	FATAL   = sentry.LevelFatal
	WARNING = sentry.LevelWarning
)
