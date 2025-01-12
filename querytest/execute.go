package querytest

import (
	"context"
	"io"

	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/dependency"
	"github.com/InfluxCommunity/flux/execute/executetest"
	"github.com/InfluxCommunity/flux/memory"
	"github.com/InfluxCommunity/flux/runtime"
)

type Querier struct{}

func (q *Querier) Query(ctx context.Context, w io.Writer, c flux.Compiler, d flux.Dialect) (int64, error) {
	program, err := c.Compile(ctx, runtime.Default)
	if err != nil {
		return 0, err
	}
	ctx, deps := dependency.Inject(ctx, executetest.NewTestExecuteDependencies())
	defer deps.Finish()
	query, err := program.Start(ctx, memory.DefaultAllocator)
	if err != nil {
		return 0, err
	}
	results := flux.NewResultIteratorFromQuery(query)
	defer results.Release()

	encoder := d.Encoder()
	return encoder.Encode(w, results)
}

func NewQuerier() *Querier {
	return &Querier{}
}
