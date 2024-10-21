package metrics

import (
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/origadmin/toolkits/errors"
	toolmetrics "github.com/origadmin/toolkits/metrics"
	"github.com/origadmin/toolkits/runtime/kratos/transport/gins"
	"go.opentelemetry.io/otel"
)

const (
	SideClient = "client"
	SideServer = "server"
)

func Middleware(side string, config *MetricConfig, logger log.Logger) (middleware.Middleware, error) {
	var (
		m   middleware.Middleware
		err error
	)
	switch side {
	case SideServer:
		m, err = ServerMiddleware(config, logger)
	case SideClient:
		m, err = ClientMiddleware(config, logger)
	default:
		return nil, errors.New("unknown metrics side")
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func ServerMiddleware(config *MetricConfig, logger log.Logger) (middleware.Middleware, error) {
	meter := otel.Meter(config.Name)
	opts := make([]metrics.Option, 0, len(config.Includes))
	if slices.Contains(config.Includes, "requests") {
		metricRequests, err := metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
		if err != nil {
			return nil, err
		}
		opts = append(opts, metrics.WithRequests(metricRequests))
	}
	if slices.Contains(config.Includes, "seconds") {
		metricSeconds, err := metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
		if err != nil {
			return nil, err
		}
		opts = append(opts, metrics.WithSeconds(metricSeconds))
	}
	return metrics.Server(opts...), nil
}

func ClientMiddleware(config *MetricConfig, logger log.Logger) (middleware.Middleware, error) {
	meter := otel.Meter(config.Name)
	opts := make([]metrics.Option, 0, len(config.Includes))
	if slices.Contains(config.Includes, "requests") {
		metricRequests, err := metrics.DefaultRequestsCounter(meter, metrics.DefaultClientRequestsCounterName)
		if err != nil {
			return nil, err
		}
		opts = append(opts, metrics.WithRequests(metricRequests))
	}
	if slices.Contains(config.Includes, "seconds") {
		metricSeconds, err := metrics.DefaultSecondsHistogram(meter, metrics.DefaultClientSecondsHistogramName)
		if err != nil {
			return nil, err
		}
		opts = append(opts, metrics.WithSeconds(metricSeconds))
	}
	return metrics.Client(opts...), nil
}

func WithMetrics(metrics toolmetrics.Metrics) gins.HandlerFunc {
	if !metrics.Enabled() {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}
	return func(ctx *gin.Context) {
		start := time.Now()
		recv := int64(0)
		if ctx.Request.ContentLength > 0 {
			recv = ctx.Request.ContentLength
		}
		ctx.Next()
		code := ctx.Writer.Status()
		send := int64(ctx.Writer.Size())
		metrics.Observe(ctx.Request.Context(), toolmetrics.MetricData{
			Endpoint: ctx.FullPath(),
			Method:   ctx.Request.Method,
			Code:     code,
			RecvSize: recv,
			SendSize: send,
			Latency:  time.Since(start).Seconds(),
			Succeed:  code < 400,
		})
	}
}
