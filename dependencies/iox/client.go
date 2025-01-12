package iox

import (
	"context"

	"github.com/InfluxCommunity/flux/codes"
	"github.com/InfluxCommunity/flux/dependencies/influxdb"
	"github.com/InfluxCommunity/flux/internal/errors"
	"github.com/InfluxCommunity/flux/memory"
	"github.com/apache/arrow/go/v7/arrow/array"
	influxdbiox "github.com/influxdata/influxdb-iox-client-go/v2"
)

type key int

const clientKey key = iota

type Config = influxdb.Config

// Dependency holds the iox.Dependency to be injected.
type Dependency struct {
	Provider Provider
}

// Inject will inject the iox dependency into the dependency chain.
func (d Dependency) Inject(ctx context.Context) context.Context {
	return context.WithValue(ctx, clientKey, d.Provider)
}

// Provider provides access to a Client with the given configuration.
type Provider interface {
	// ClientFor will return a client with the given configuration.
	ClientFor(ctx context.Context, conf Config) (Client, error)
}

// GetProvider retrieves the iox Provider.
func GetProvider(ctx context.Context) Provider {
	p := ctx.Value(clientKey)
	if p == nil {
		return ErrorProvider{}
	}
	return p.(Provider)
}

// RecordReader is similar to the RecordReader interface provided by Arrow's array
// package, but includes a method for detecting errors that are sent mid-stream.
type RecordReader interface {
	array.RecordReader
	Err() error
}

// Client provides a way to query an iox instance.
type Client interface {
	// Query will initiate a query using the given query string, parameters, and memory allocator
	// against the iox instance. It returns an array.RecordReader from the arrow flight api.
	Query(ctx context.Context, query string, params []interface{}, mem memory.Allocator) (RecordReader, error)

	// GetSchema will retrieve a schema for the given table if this client supports that capability.
	// If this Client doesn't support this capability, it should return a flux error with the code
	// codes.Unimplemented.
	GetSchema(ctx context.Context, table string) (map[string]influxdbiox.ColumnType, error)
}

// ErrorProvider is an implementation of the Provider that returns an error.
type ErrorProvider struct{}

func (u ErrorProvider) ClientFor(ctx context.Context, conf Config) (Client, error) {
	return nil, errors.New(codes.Invalid, "iox client has not been configured")
}
