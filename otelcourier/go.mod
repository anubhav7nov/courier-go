module github.com/gojekfarm/courier-go/otelcourier

go 1.16

require (
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v1.0.0
	go.opentelemetry.io/otel/sdk v1.0.0
	go.opentelemetry.io/otel/trace v1.0.0
	github.com/gojekfarm/courier-go v0.6.0
)

replace github.com/gojekfarm/courier-go => ../
