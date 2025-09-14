package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		log.Infof("[GIN] | %v | %s | %s | %s | %s",
			c.Writer.Status(),
			latency,
			c.Request.Method,
			path,
			raw,
		)
	}
}
