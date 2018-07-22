package httputil

import (
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/kunhou/TMDB/log"
)

type errorMsg struct {
	StatusCode    int64  `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

// ResponseFail responses Fail response and add error information to errors
func ResponseFail(c *gin.Context, httpStatusCode int, code int64, msg string, err error) {
	meta := errorMsg{
		StatusCode:    code,
		StatusMessage: msg,
	}
	c.JSON(httpStatusCode, meta)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		log.WithError(err).WithFields(log.Fields{
			"pc":     pc,
			"file":   file,
			"line":   line,
			"status": code,
		}).Error(msg)
	} else {
		log.WithFields(log.Fields{
			"status": code,
		}).Warning(msg)
	}
}
