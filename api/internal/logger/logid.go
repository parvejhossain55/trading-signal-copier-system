package logger

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"log/slog"

	"github.com/google/uuid"
)

type logIDKeyType struct{}

var LogIDKey = logIDKeyType{}

const logIDString = "log_id"

var (
	sequentialCounter uint64
	counterMutex      sync.Mutex
)

type LogIDGenerator struct{}

func NewLogIDGenerator() *LogIDGenerator {
	return &LogIDGenerator{}
}

func (g *LogIDGenerator) GenerateUUID() string {
	return uuid.New().String()
}

func (g *LogIDGenerator) GenerateSequential() string {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	sequentialCounter++
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("log_%s_%06d", timestamp, sequentialCounter)
}

func (g *LogIDGenerator) GenerateRandom() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("log_%x", bytes)
}

func (g *LogIDGenerator) GenerateTimestamped() string {
	timestamp := time.Now().Format("20060102150405.000")
	return fmt.Sprintf("log_%s", timestamp)
}

func (g *LogIDGenerator) GenerateWithPrefix(prefix string) string {
	timestamp := time.Now().Format("20060102150405")
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	return fmt.Sprintf("%s_%s_%x", prefix, timestamp, randomBytes)
}

func (g *LogIDGenerator) GenerateFromContext(ctx context.Context) string {
	if logID, ok := ctx.Value(LogIDKey).(string); ok && logID != "" {
		return logID
	}
	return g.GenerateUUID()
}

func ContextWithLogID(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, LogIDKey, logID)
}

func ContextWithNewLogID(ctx context.Context) context.Context {
	generator := NewLogIDGenerator()
	logID := generator.GenerateUUID()
	return context.WithValue(ctx, LogIDKey, logID)
}

func GetLogIDFromContext(ctx context.Context) string {
	if logID, ok := ctx.Value(LogIDKey).(string); ok {
		return logID
	}
	return ""
}

func LogID(id string) slog.Attr {
	return slog.String(logIDString, id)
}

func LogIDFromContext(ctx context.Context) slog.Attr {
	logID := GetLogIDFromContext(ctx)
	if logID == "" {
		logID = NewLogIDGenerator().GenerateUUID()
	}
	return slog.String(logIDString, logID)
}

// Global convenience functions
var (
	defaultGenerator = NewLogIDGenerator()
)

func GenerateLogID() string {
	return defaultGenerator.GenerateUUID()
}

func GenerateSequentialLogID() string {
	return defaultGenerator.GenerateSequential()
}

func GenerateRandomLogID() string {
	return defaultGenerator.GenerateRandom()
}

func GenerateTimestampedLogID() string {
	return defaultGenerator.GenerateTimestamped()
}

func GenerateLogIDWithPrefix(prefix string) string {
	return defaultGenerator.GenerateWithPrefix(prefix)
}
