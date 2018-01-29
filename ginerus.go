package ginerus

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Ginerus() gin.HandlerFunc {
	return GinerusWithLogger(logrus.StandardLogger(), nil)
}

func GinerusWithLogger(logger *logrus.Logger, config map[string]string) gin.HandlerFunc {
	
	// parse config
	if _, ok := config["format"]; ok {
		if config["format"] == "json" {
			// Log as JSON instead of the default ASCII formatter.
			logrus.SetFormatter(&logrus.JSONFormatter{})
		}
	}
	if _, ok := config["output"]; ok {
		if config["output"] == "stdout" {
			logrus.SetOutput(os.Stdout)
		}
		if config["output"] == "stderr" {
			logrus.SetOutput(os.Stderr)
		}
	}
	
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.String()
		userAgent := c.Request.UserAgent()

		timeFormatted := end.Format("2006-01-02 15:04:05")

		msg := fmt.Sprintf(
			"%s %s \"%s %s\" %d %s %s",
			clientIP,
			timeFormatted,
			method,
			path,
			statusCode,
			latency,
			userAgent,
		)

		logger.WithFields(logrus.Fields{
			"time":       timeFormatted,
			"method":     method,
			"path":       path,
			"latency":    latency,
			"ip":         clientIP,
			"comment":    comment,
			"status":     statusCode,
			"user-agent": userAgent,
		}).Info(msg)
	}
}
