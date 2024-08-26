package ginzap

import (
	"errors"
	"fx-vote-server/common/app"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	LogRequest *bool
	LogStack   *bool
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return CustomRecoveryWithZap(defaultHandleRecovery)
}

func defaultHandleRecovery(c *gin.Context, _ any) {
	c.AbortWithStatus(http.StatusInternalServerError)
}

func CustomRecoveryWithZap(recovery gin.RecoveryFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				if brokenPipe {
					{
						app.GinLog.Error(c.Request.URL.Path, zap.Any("error", err))
						app.GinLog.Error(string(httpRequest))
					}

					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) //nolint: errcheck
					c.Abort()
					return
				}

				{
					app.GinLog.Error("[Recovery from panic]", zap.Any("error", err))

					if *LogRequest {
						app.GinLog.Error(string(httpRequest))
					}
					if *LogStack {
						app.GinLog.Error(string(debug.Stack()))
					}
				}

				recovery(c, err)
			}
		}()
		c.Next()
	}
}
