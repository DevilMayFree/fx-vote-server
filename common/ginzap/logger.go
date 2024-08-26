package ginzap

import (
	"fx-vote-server/common/app"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {

	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

func Logger() gin.HandlerFunc {
	return LoggerWithConfig(&LoggerConfig{})
}

func LoggerWithConfig(conf *LoggerConfig) gin.HandlerFunc {
	notLogged := conf.SkipPaths
	skipPaths := make(map[string]bool, len(notLogged))

	for _, path := range notLogged {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		// Process request
		c.Next()
		// Log only when path is not being skipped
		if _, ok := skipPaths[path]; !ok {

			end := time.Now()
			latency := end.Sub(start)

			if raw != "" {
				path = path + "?" + raw
			}

			fields := []zapcore.Field{
				zap.Int("status", c.Writer.Status()),
				zap.Duration("latency", latency),
				zap.String("ip", c.ClientIP()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", raw),
			}

			if len(c.Errors) > 0 {
				for _, e := range c.Errors.Errors() {
					app.GinLog.Error(e, fields...)
				}
			} else {
				app.GinLog.Debug("[GIN]", fields...)
			}
		}
	}

}
