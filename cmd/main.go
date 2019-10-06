package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPRequestTotal    *prometheus.CounterVec
)

func init() {
	HTTPRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "HTTP Request latencies in seconds",
	}, []string{"method", "path", "status"})
	HTTPRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})
	prometheus.MustRegister(
		HTTPRequestDuration,
		HTTPRequestTotal,
	)
}

func main() {
	engine := gin.Default()
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	engine.Use(httpMetrics())
	engine.GET("/hello", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"hello": "world"})
	})
	log.Fatal(engine.Run(":8000"))
}

func httpMetrics() gin.HandlerFunc {
	return func(context *gin.Context) {
		startRequestTime := time.Now()
		context.Next()

		duration := float64(time.Since(startRequestTime)) / float64(time.Second)
		path := context.Request.URL.Path

		HTTPRequestTotal.With(prometheus.Labels{
			"method": context.Request.Method,
			"path":   path,
			"status": strconv.Itoa(context.Writer.Status()),
		}).Inc()
		HTTPRequestDuration.With(prometheus.Labels{
			"method": context.Request.Method,
			"path":   path,
			"status": strconv.Itoa(context.Writer.Status()),
		}).Observe(duration)
	}

}
