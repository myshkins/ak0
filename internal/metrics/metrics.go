package metrics

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

  r, err := resource.Merge(
    resource.Default(),
    resource.NewWithAttributes(
      semconv.SchemaURL,
      semconv.ServiceName("ak0_2"),
      ),
    )
  if err != nil {
    return nil, err
  }

	// Set up meter provider.
	meterProvider, err := newMeterProvider(r)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newMeterProvider(r *resource.Resource) (*metric.MeterProvider, error) {
  // todo: double check this
  metricExporter, err := otlpmetrichttp.New(context.Background())
	if err != nil {
		return nil, err
	}

  var view metric.View = func(i metric.Instrument) (metric.Stream, bool) {
	// In a custom View function, you need to explicitly copy
	// the name, description, and unit.
	s := metric.Stream{
      Name: strings.ReplaceAll(i.Name, ".", "_"),
      Description: i.Description,
      Unit: i.Unit,
    }
	return s, true
  }

	meterProvider := metric.NewMeterProvider(
    metric.WithResource(r),
		metric.WithReader(metric.NewPeriodicReader(
      metricExporter,
      // todo: should maybe be 1m, change this
			metric.WithInterval(5*time.Second))),
    metric.WithView(view),
	)
	return meterProvider, nil
}

