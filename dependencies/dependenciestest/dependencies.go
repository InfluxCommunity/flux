package dependenciestest

import (
	"context"
	"io"
	"net/http"

	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/dependencies/filesystem"
	"github.com/InfluxCommunity/flux/dependencies/influxdb"
	"github.com/InfluxCommunity/flux/dependencies/mqtt"
	"github.com/InfluxCommunity/flux/dependencies/url"
	"github.com/InfluxCommunity/flux/dependency"
	"github.com/InfluxCommunity/flux/execute"
	"github.com/InfluxCommunity/flux/mock"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

var StatusOK int = 200

func defaultTestFunction(req *http.Request) *http.Response {
	body := (*req).Body
	// Test request parameters
	return &http.Response{
		StatusCode: StatusOK,
		Status:     "Body generated by test client",

		// Send response to be tested
		Body: io.NopCloser(body),

		// Must be set to non-nil value or it panics
		Header: make(http.Header),
	}
}

type Deps struct {
	flux.Deps
	influxdb influxdb.Dependency
	mqtt     mqtt.Dependency
}

func (d Deps) Inject(ctx context.Context) context.Context {
	ctx = d.Deps.Inject(ctx)
	ctx = d.influxdb.Inject(ctx)
	return d.mqtt.Inject(ctx)
}

func Default() Deps {
	var deps flux.Deps

	deps.Deps.HTTPClient = &http.Client{
		Transport: RoundTripFunc(defaultTestFunction),
	}
	deps.Deps.SecretService = &mock.SecretService{
		"password": "mysecretpassword",
		"token":    "mysecrettoken",
	}
	deps.Deps.FilesystemService = filesystem.SystemFS
	deps.Deps.URLValidator = url.PassValidator{}
	return Deps{
		Deps: deps,
		influxdb: influxdb.Dependency{
			Provider: influxdb.HttpProvider{},
		},
		mqtt: mqtt.Dependency{
			Dialer: mqtt.DefaultDialer{},
		},
	}
}

func ExecutionDefault() execute.ExecutionDependencies {
	return execute.DefaultExecutionDependencies()
}

// Injects all default dependencies into the context
func InjectAllDeps(ctx context.Context) (context.Context, *dependency.Span) {
	deps := Default()
	execDeps := ExecutionDefault()
	return dependency.Inject(ctx, deps, execDeps)
}
