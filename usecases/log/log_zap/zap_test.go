package logzap

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZap_Usage(t *testing.T) {
	var buf bytes.Buffer

	// 创建一个 zapcore.Core，将日志输出到 buf
	writer := zapcore.AddSync(&buf)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)

	// 创建 logger
	logger := zap.New(core)
	defer logger.Sync()

	zap.NewProduction()

	// 使用 logger
	logger.Info("Hello, Zap!")
	logger.With(zap.String("key", "value")).Info("This is an info message")

	// 检查 buf 的内容
	t.Log("Logged output:", buf.String())
	assert.Contains(t, buf.String(), "Hello, Zap!")
	assert.Contains(t, buf.String(), "This is an info message")
}
