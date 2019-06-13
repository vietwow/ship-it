package ecrconsumer

import (
	"context"
	"time"

	"github.com/Wattpad/sqsconsumer"
	"github.com/Wattpad/sqsconsumer/middleware"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

func NewSQSConsumer(logger log.Logger, hist metrics.Histogram, name string, svc sqsconsumer.SQSAPI) (*sqsconsumer.Consumer, error) {
	service, err := sqsconsumer.NewSQSService(name, svc)
	if err != nil {
		return nil, err
	}

	track := dataDogTimeTracker(hist)
	wrappedLogger := loggerMiddleware(logger)
	handler := middleware.ApplyDecoratorsToHandler(processMessage, track, wrappedLogger)

	consumer := sqsconsumer.NewConsumer(service, handler)

	return consumer, nil
}

func processMessage(ctx context.Context, msg string) error {
	return nil
}

func dataDogTimeTracker(hist metrics.Histogram) middleware.MessageHandlerDecorator {
	return func(fn sqsconsumer.MessageHandlerFunc) sqsconsumer.MessageHandlerFunc {
		return func(ctx context.Context, msg string) error {
			start := time.Now()

			err := fn(ctx, msg)

			var status string
			if err != nil {
				status = "failure"
			} else {
				status = "success"
			}
			hist.With("status", status).Observe(float64(time.Since(start).Seconds() * 1000))

			return err
		}
	}
}

func loggerMiddleware(logger log.Logger) middleware.MessageHandlerDecorator {
	return func(fn sqsconsumer.MessageHandlerFunc) sqsconsumer.MessageHandlerFunc {
		return func(ctx context.Context, msg string) error {
			err := fn(ctx, msg)
			if err != nil {
				logger.Log("error", err)
			}
			return err
		}
	}
}