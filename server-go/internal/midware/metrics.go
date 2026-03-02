package midware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var apiCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "httpcat_http_api_qps_counter",
	Help: "Httpcat Http API QPS",
}, []string{"handle", "code"})

var apiHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "httpcat_http_api_histogram",
	Help:    "Httpcat Http API Histogram",
	Buckets: prometheus.DefBuckets,
}, []string{"handle", "code"})

var initOnce = &sync.Once{}

func Metrics() gin.HandlerFunc {
	initOnce.Do(func() {
		prometheus.MustRegister(apiCounter)
		prometheus.MustRegister(apiHistogram)
	})
	return func(c *gin.Context) {
		begin := time.Now()

		handle := c.HandlerName()
		c.Next()
		code := fmt.Sprint(c.Writer.Status())
		apiCounter.With(prometheus.Labels{"handle": handle, "code": code}).Inc()
		apiHistogram.With(prometheus.Labels{"handle": handle, "code": code}).Observe(time.Since(begin).Seconds())
	}
}
