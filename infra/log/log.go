package log

import (
	"api-gin/config"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
	"time"
)

type Logger struct {
	*logrus.Entry

	mu          sync.Mutex
	lastTraceId string
	lastSpan    uint
}

// NewLogger 根据配置创建一个新的Logger实例
func NewLogger(config *config.Config) (*Logger, error) {
	if config == nil {
		return nil, fmt.Errorf("[NewLogger] 配置不能为空")
	}
	logConfig := config.Log
	logger := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(logConfig.Level)
	if err != nil {
		return nil, fmt.Errorf("无效的日志级别: %v", logConfig.Level)
	}
	logger.SetLevel(level)

	// 设置日志输出格式
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
	})
	if logConfig.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}

	// 设置日志输出位置
	var output io.Writer
	switch logConfig.Mode {
	case "command":
		output = os.Stdout
	case "file":
		// 按日期切割
		fileName := fmt.Sprintf("%s/%s_output.", logConfig.Path, logConfig.FileName) + "%Y%m%d"
		output, err = rotatelogs.New(fileName,
			rotatelogs.WithMaxAge(time.Duration(logConfig.MaxAge)*24*time.Hour),
			rotatelogs.WithRotationSize(int64(logConfig.MaxSize)*1024*1024),
			//rotatelogs.WithRotationCount(logConfig.MaxBackups), // MaxAge只能留一个
		)
		if err != nil {
			return nil, fmt.Errorf("初始化日志写入器错误：%v", err)
		}
	default:
		output = os.Stdout
	}
	logger.SetOutput(output)
	// TODO 自定义hook
	e := &Logger{Entry: logger.WithField("app", config.Name)}

	return e, nil
}

func (l *Logger) NewLogger(call string) *Logger {
	entry := l.Entry.WithField("caller", call)
	return &Logger{
		Entry: entry,
	}
}

type KeyTraceKey struct{}

// SetTraceId 设置traceId
func SetTraceId(ctx context.Context) context.Context {
	return context.WithValue(ctx, KeyTraceKey{}, uuid.New().String())
}

// Trace 自动增加traceId
func (l *Logger) Trace(ctx context.Context) logrus.Fields {
	l.mu.Lock()
	defer l.mu.Unlock()
	nowTraceIdStr := ""
	if nowTraceId := ctx.Value(KeyTraceKey{}); nowTraceId != nil {
		if nowTraceId, ok := nowTraceId.(string); ok {
			nowTraceIdStr = nowTraceId
		}
	}
	if nowTraceIdStr == "" {
		nowTraceIdStr = uuid.New().String()
	}
	if nowTraceIdStr != l.lastTraceId {
		l.lastTraceId = nowTraceIdStr
		l.lastSpan = 0
	}
	l.lastSpan++
	return logrus.Fields{"trace_id": nowTraceIdStr, "span": l.lastSpan}
}

func (l *Logger) Debug(ctx context.Context, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Debug(args...)
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Debugf(format, args...)
}

func (l *Logger) Info(ctx context.Context, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Info(args...)
}

func (l *Logger) Infof(ctx context.Context, format string, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Infof(format, args...)
}

func (l *Logger) Warn(ctx context.Context, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Warn(args...)
}

func (l *Logger) Warnf(ctx context.Context, format string, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Warnf(format, args...)
}

func (l *Logger) Error(ctx context.Context, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Error(args...)
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Errorf(format, args...)
}

func (l *Logger) Fatal(ctx context.Context, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Fatal(args...)
}

func (l *Logger) Fatalf(ctx context.Context, format string, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Fatalf(format, args...)
}

func (l *Logger) Panic(ctx context.Context, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Panic(args...)
}

func (l *Logger) Panicf(ctx context.Context, format string, args ...any) {
	l.Entry.WithFields(l.Trace(ctx)).Panicf(format, args...)
}
