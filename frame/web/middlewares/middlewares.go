package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	htResult "github.com/qq754174349/ht/ht-frame/common/result"
	"github.com/qq754174349/ht/ht-frame/config"
	"github.com/qq754174349/ht/ht-frame/logger"
	"github.com/qq754174349/ht/ht-frame/web/prometheus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body string
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	w.body = s // 记录响应内容
	return w.ResponseWriter.WriteString(s)
}

// TraceIDMiddleware 负责生成或透传 TraceID
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader("X-Trace-Id")
		if traceId == "" {
			traceId = uuid.New().String()
			c.Request.Header.Set("X-Trace-Id", traceId)
		}
		c.Set("traceID", traceId)

		c.Next()
	}
}

// HeaderPropagationMiddleware 负责透传特定的请求头到响应头
func HeaderPropagationMiddleware(headers ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, h := range headers {
			if v := c.GetHeader(h); v != "" {
				c.Header(h, v)
			}
		}
		c.Next()
	}
}

// RequestInfoLogger 请求参数打印中间件
func RequestInfoLogger(SkipPaths []string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(SkipPaths); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range SkipPaths {
			skip[path] = struct{}{}
		}
	}
	env := config.GetAppCfg().Active

	return func(c *gin.Context) {
		start := time.Now()
		raw := c.Request.URL.RawQuery
		path := c.Request.URL.Path
		if raw != "" {
			path = path + "?" + raw
		}

		if _, ok := skip[path]; ok {
			return
		}
		writer := &bodyLogWriter{
			body:           "",
			ResponseWriter: c.Writer,
		}
		c.Writer = writer
		c.Writer.Header().Set("content-type", "application/json")

		c.Next()

		// 获取请求头并转化为 JSON
		reqHeaders, err := json.Marshal(c.Request.Header)
		if err != nil {
			reqHeaders = []byte("{}")
		}

		// 保存请求体数据，并且不影响后续中间件
		var body string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.WithTraceID(c.GetString("traceID")).Errorf("Failed to read request body: %v", err)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				return
			}
			body = string(bodyBytes)
			if env == "pro" && len(body) > 100 {
				body = body[:100]
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		// 获取 traceID
		traceID, _ := c.Get("traceID")

		latency := time.Now().Sub(start)
		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}
		respHeaders, err := json.Marshal(c.Writer.Header())
		if err != nil {
			respHeaders = []byte("{}")
		}

		respBody := ""
		if c.Writer.Header().Get("Content-Type") == "application/json" {
			respBody = writer.body
			if env == "pro" && len(respBody) > 100 {
				respBody = respBody[:100]
			}
		}

		// 记录请求信息
		logger.WithTraceID(traceID.(string)).Infof(
			"[GIN] Method=%s, Path=%s, StatusCode=%3d, Cost=%s, Ip=%s, ReqHeaders=%s, Request=%s, Response=%s, RespHeaders=%s",
			c.Request.Method, path, c.Writer.Status(), latency, c.ClientIP(), string(reqHeaders), body, respBody, respHeaders,
		)

	}
}

// ErrorHandler 全局异常处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Errors == nil {
			return
		}
		htResult.FailByErr(c, c.Errors.Last().Err)
	}
}

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// Prometheus 普罗米修斯监控请求
func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 设置路由
		route := c.FullPath()

		// 继续处理请求
		c.Next()

		// 记录请求持续时间
		prometheus.Duration.WithLabelValues(c.Request.Method, route).Observe(time.Since(start).Seconds())

		// 更新请求计数器
		prometheus.Requests.WithLabelValues(c.Request.Method, route).Inc()
	}
}
