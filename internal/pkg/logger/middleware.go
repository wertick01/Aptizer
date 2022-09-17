package logger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinMiddlewareLog(l ILogger, loqRequestBody bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		d := new(Detail)

		if loqRequestBody && c.Request.Body != nil {
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(c.Request.Body)
			b, _ := ioutil.ReadAll(c.Request.Body)
			d.Request = string(b)

			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		duration := time.Since(start)

		d.Address = c.Request.URL.Path
		d.Response = blw.body.String()
		d.ResponseCode = c.Writer.Status()
		d.Fields = map[string]any{
			"method":     c.Request.Method,
			"query":      c.Request.URL.RawQuery,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
			"duration":   duration,
		}

		f := &Fields{
			Detail: d,
			UserID: c.Request.Header.Get("X-User-Id"),
			Code:   "api_server",
			Message: fmt.Sprintf(
				"completed with %d %s in %v",
				d.ResponseCode,
				http.StatusText(d.ResponseCode),
				duration,
			),
		}

		if c.Writer.Status() >= 400 {
			l.Warn(f)
		} else {
			l.Info(f)
		}
	}
}
