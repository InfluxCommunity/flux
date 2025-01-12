package gen_test

import (
	"context"
	"testing"

	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/execute/executetest"
	"github.com/InfluxCommunity/flux/internal/gen"
	"github.com/InfluxCommunity/flux/memory"
)

func TestInput_TableTest(t *testing.T) {
	executetest.RunTableTests(t, executetest.TableTest{
		NewFn: func(ctx context.Context, alloc memory.Allocator) flux.TableIterator {
			schema := gen.Schema{
				Tags: []gen.Tag{
					{Name: "_measurement", Cardinality: 1},
					{Name: "_field", Cardinality: 1},
					{Name: "t0", Cardinality: 100},
				},
				NumPoints: 100,
				Alloc:     alloc,
			}
			tables, err := gen.Input(context.Background(), schema)
			if err != nil {
				t.Fatal(err)
			}
			return tables
		},
		IsDone: func(tbl flux.Table) bool {
			return tbl.(interface{ IsDone() bool }).IsDone()
		},
	})
}

func benchmarkInput(b *testing.B, n int) {
	schema := gen.Schema{
		Tags: []gen.Tag{
			{Name: "_measurement", Cardinality: 1},
			{Name: "_field", Cardinality: 1},
			{Name: "t0", Cardinality: 100},
		},
		NumPoints: n,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ti, err := gen.Input(context.Background(), schema)
		if err != nil {
			b.Fatal(err)
		}

		if err := ti.Do(func(tbl flux.Table) error {
			tbl.Done()
			return nil
		}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInput(b *testing.B) {
	b.Run("1000", func(b *testing.B) {
		benchmarkInput(b, 1000)
	})
	b.Run("100000", func(b *testing.B) {
		benchmarkInput(b, 100000)
	})
	b.Run("1000000", func(b *testing.B) {
		benchmarkInput(b, 1000000)
	})
}
