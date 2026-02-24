package logger

import (
	"fmt"
	"log/slog"
)

type StartupLogger struct {
	logger *slog.Logger
}

func NewStartupLogger(logger *slog.Logger) *StartupLogger {
	return &StartupLogger{
		logger: logger,
	}
}

func (s *StartupLogger) Success(message string, args ...any) {
	icon := "✅"
	if len(args) > 0 {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message), args...)
	} else {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message))
	}
}

func (s *StartupLogger) Info(message string, args ...any) {
	icon := "ℹ️"
	if len(args) > 0 {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message), args...)
	} else {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message))
	}
}

func (s *StartupLogger) Warning(message string, args ...any) {
	icon := "⚠️"
	if len(args) > 0 {
		s.logger.Warn(fmt.Sprintf("%s %s", icon, message), args...)
	} else {
		s.logger.Warn(fmt.Sprintf("%s %s", icon, message))
	}
}

func (s *StartupLogger) Error(message string, args ...any) {
	icon := "❌"
	if len(args) > 0 {
		s.logger.Error(fmt.Sprintf("%s %s", icon, message), args...)
	} else {
		s.logger.Error(fmt.Sprintf("%s %s", icon, message))
	}
}

func (s *StartupLogger) Service(message string, args ...any) {
	icon := "🚀"
	if len(args) > 0 {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message), args...)
	} else {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message))
	}
}

func (s *StartupLogger) Database(message string, args ...any) {
	icon := "🗄️"
	if len(args) > 0 {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message), args...)
	} else {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message))
	}
}

func (s *StartupLogger) Cache(message string, args ...any) {
	icon := "⚡"
	if len(args) > 0 {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message), args...)
	} else {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message))
	}
}

func (s *StartupLogger) Messaging(message string, args ...any) {
	icon := "📨"
	if len(args) > 0 {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message), args...)
	} else {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message))
	}
}

func (s *StartupLogger) Stream(message string, args ...any) {
	icon := "🌊"
	if len(args) > 0 {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message), args...)
	} else {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message))
	}
}

func (s *StartupLogger) Subscription(message string, args ...any) {
	icon := "📡"
	if len(args) > 0 {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message), args...)
	} else {
		s.logger.Info(fmt.Sprintf("%s %s", icon, message))
	}
}
