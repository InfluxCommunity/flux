package universe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/dependencies/dependenciestest"
	"github.com/InfluxCommunity/flux/dependency"
	"github.com/InfluxCommunity/flux/execute"
	"github.com/InfluxCommunity/flux/execute/executetest"
	"github.com/InfluxCommunity/flux/internal/gen"
	"github.com/InfluxCommunity/flux/internal/operation"
	"github.com/InfluxCommunity/flux/interpreter"
	"github.com/InfluxCommunity/flux/memory"
	"github.com/InfluxCommunity/flux/querytest"
	"github.com/InfluxCommunity/flux/runtime"
	"github.com/InfluxCommunity/flux/semantic"
	"github.com/InfluxCommunity/flux/stdlib/influxdata/influxdb"
	"github.com/InfluxCommunity/flux/stdlib/universe"
	"github.com/InfluxCommunity/flux/values"
	"github.com/InfluxCommunity/flux/values/valuestest"
)

func TestMap_NewQuery(t *testing.T) {
	tests := []querytest.NewQueryTestCase{
		{
			Name: "simple static map",
			Raw:  `from(bucket:"mybucket") |> map(fn: (r) => r._value + 1)`,
			Want: &operation.Spec{
				Operations: []*operation.Node{
					{
						ID: "from0",
						Spec: &influxdb.FromOpSpec{
							Bucket: influxdb.NameOrID{Name: "mybucket"},
						},
					},
					{
						ID: "map1",
						Spec: &universe.MapOpSpec{
							Fn: interpreter.ResolvedFunction{
								Fn:    executetest.FunctionExpression(t, "(r) => r._value + 1"),
								Scope: valuestest.Scope(),
							},
						},
					},
				},
				Edges: []operation.Edge{
					{Parent: "from0", Child: "map1"},
				},
			},
		},
		{
			Name: "simple static map mergeKey=true",
			Raw:  `from(bucket:"mybucket") |> map(fn: (r) => r._value + 1, mergeKey:true)`,
			Want: &operation.Spec{
				Operations: []*operation.Node{
					{
						ID: "from0",
						Spec: &influxdb.FromOpSpec{
							Bucket: influxdb.NameOrID{Name: "mybucket"},
						},
					},
					{
						ID: "map1",
						Spec: &universe.MapOpSpec{
							MergeKey: true,
							Fn: interpreter.ResolvedFunction{
								Fn:    executetest.FunctionExpression(t, "(r) => r._value + 1"),
								Scope: valuestest.Scope(),
							},
						},
					},
				},
				Edges: []operation.Edge{
					{Parent: "from0", Child: "map1"},
				},
			},
		},
		{
			Name: "resolve map",
			Raw:  `x = 2 from(bucket:"mybucket") |> map(fn: (r) => r._value + x)`,
			Want: &operation.Spec{
				Operations: []*operation.Node{
					{
						ID: "from0",
						Spec: &influxdb.FromOpSpec{
							Bucket: influxdb.NameOrID{Name: "mybucket"},
						},
					},
					{
						ID: "map1",
						Spec: &universe.MapOpSpec{
							Fn: interpreter.ResolvedFunction{
								Fn: executetest.FunctionExpression(t, "(r) => r._value + 2"),
								Scope: func() values.Scope {
									scope := valuestest.Scope()
									scope.Set("x", values.NewInt(2))
									return scope
								}(),
							},
						},
					},
				},
				Edges: []operation.Edge{
					{Parent: "from0", Child: "map1"},
				},
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			querytest.NewQueryTestHelper(t, tc)
		})
	}
}

func TestMap_Process(t *testing.T) {
	builtIns := runtime.Prelude()
	testCases := []struct {
		name    string
		spec    *universe.MapProcedureSpec
		data    []flux.Table
		want    []*executetest.Table
		wantErr error
	}{
		{
			name: `overwrite groupkey`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({r with _stop: "a"})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				KeyCols: []string{"_stop"},
				ColMeta: []flux.ColMeta{
					{Label: "_stop", Type: flux.TTime},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(6), execute.Time(1), 1.0},
					{execute.Time(6), execute.Time(2), 2.0},
					{execute.Time(6), execute.Time(3), 3.0},
					{execute.Time(6), execute.Time(4), 4.0},
					{execute.Time(6), execute.Time(5), 5.0},
				},
			}},
			want: []*executetest.Table{{
				KeyCols: []string{"_stop"},
				ColMeta: []flux.ColMeta{
					{Label: "_stop", Type: flux.TString},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{"a", execute.Time(1), 1.0},
					{"a", execute.Time(2), 2.0},
					{"a", execute.Time(3), 3.0},
					{"a", execute.Time(4), 4.0},
					{"a", execute.Time(5), 5.0},
				},
			}},
		},
		{
			// Parameter mergeKey is now deprecated and has no effect.
			// `map` is evaluated as if it was false.
			name: `overwrite mergekey`,
			spec: &universe.MapProcedureSpec{
				MergeKey: true,
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_stop: "a", _value: "b"})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				KeyCols: []string{"_start", "_stop"},
				ColMeta: []flux.ColMeta{
					{Label: "_start", Type: flux.TTime},
					{Label: "_stop", Type: flux.TTime},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(0), execute.Time(6), execute.Time(1), 1.0},
					{execute.Time(0), execute.Time(6), execute.Time(2), 2.0},
					{execute.Time(0), execute.Time(6), execute.Time(3), 3.0},
					{execute.Time(0), execute.Time(6), execute.Time(4), 4.0},
					{execute.Time(0), execute.Time(6), execute.Time(5), 5.0},
				},
			}},
			want: []*executetest.Table{{
				KeyCols: []string{"_stop"},
				ColMeta: []flux.ColMeta{
					{Label: "_stop", Type: flux.TString},
					{Label: "_value", Type: flux.TString},
				},
				Data: [][]interface{}{
					{"a", "b"},
					{"a", "b"},
					{"a", "b"},
					{"a", "b"},
					{"a", "b"},
				},
			}},
		},
		{
			name: `identity function`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, "(r) => r"),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 1.0},
					{execute.Time(2), 6.0},
					{execute.Time(3), nil},
					{execute.Time(4), 7.0},
					{execute.Time(5), nil},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 1.0},
					{execute.Time(2), 6.0},
					{execute.Time(3), nil},
					{execute.Time(4), 7.0},
					{execute.Time(5), nil},
				},
			}},
		},
		{
			name: `_value+5`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, "(r) => ({_time: r._time, _value: r._value + 5.0})"),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 1.0},
					{execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 6.0},
					{execute.Time(2), 11.0},
				},
			}},
		},
		{
			name: `_value+5 custom scope`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: func() values.Scope {
						scope := values.NewScope()
						scope.Set("x", values.NewFloat(5))
						return scope
					}(),
					Fn: executetest.FunctionExpression(t, `
x = 5.0
f = (r) => ({_time: r._time, _value: r._value + x})
f
`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 1.0},
					{execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 6.0},
					{execute.Time(2), 11.0},
				},
			}},
		},
		{
			name: `_value+5 drop cols`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, "(r) => ({_time: r._time, _value: r._value + 5.0})"),
				},
			},
			data: []flux.Table{&executetest.Table{
				KeyCols: []string{"_measurement", "host"},
				ColMeta: []flux.ColMeta{
					{Label: "_measurement", Type: flux.TString},
					{Label: "host", Type: flux.TString},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{"m", "A", execute.Time(1), 1.0},
					{"m", "A", execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 6.0},
					{execute.Time(2), 11.0},
				},
			}},
		},
		{
			name: `_value+5 drop cols regroup`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, "(r) => ({_time: r._time, _value: r._value + 5.0})"),
				},
			},
			data: []flux.Table{
				&executetest.Table{
					KeyCols: []string{"_measurement", "host"},
					ColMeta: []flux.ColMeta{
						{Label: "_measurement", Type: flux.TString},
						{Label: "host", Type: flux.TString},
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{"m", "A", execute.Time(1), 1.0},
						{"m", "A", execute.Time(2), 6.0},
					},
				},
				&executetest.Table{
					KeyCols: []string{"_measurement", "host"},
					ColMeta: []flux.ColMeta{
						{Label: "_measurement", Type: flux.TString},
						{Label: "host", Type: flux.TString},
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{"m", "B", execute.Time(1), 1.5},
						{"m", "B", execute.Time(2), 6.5},
					},
				},
			},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 6.0},
					{execute.Time(2), 11.0},
					{execute.Time(1), 6.5},
					{execute.Time(2), 11.5},
				},
			}},
		},
		{
			name: `with _value+5 regroup`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({r with host: r.host + ".local", _value: r._value + 5.0})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				KeyCols: []string{"_measurement", "host"},
				ColMeta: []flux.ColMeta{
					{Label: "_measurement", Type: flux.TString},
					{Label: "host", Type: flux.TString},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{"m", "A", execute.Time(1), 1.0},
					{"m", "A", execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{{
				KeyCols: []string{"_measurement", "host"},
				ColMeta: []flux.ColMeta{
					{Label: "_measurement", Type: flux.TString},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "host", Type: flux.TString},
				},
				Data: [][]interface{}{
					{"m", execute.Time(1), 6.0, "A.local"},
					{"m", execute.Time(2), 11.0, "A.local"},
				},
			}},
		},
		{
			name: `with _value+5 regroup fan out`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({r with host: r.host + "." + r.domain, _value: r._value + 5.0})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				KeyCols: []string{"_measurement", "host"},
				ColMeta: []flux.ColMeta{
					{Label: "_measurement", Type: flux.TString},
					{Label: "host", Type: flux.TString},
					{Label: "domain", Type: flux.TString},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{"m", "A", "example.com", execute.Time(1), 1.0},
					{"m", "A", "www.example.com", execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{
				{
					KeyCols: []string{"_measurement", "host"},
					ColMeta: []flux.ColMeta{
						{Label: "_measurement", Type: flux.TString},
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
						{Label: "domain", Type: flux.TString},
						{Label: "host", Type: flux.TString},
					},
					Data: [][]interface{}{
						{"m", execute.Time(1), 6.0, "example.com", "A.example.com"},
					},
				},
				{
					KeyCols: []string{"_measurement", "host"},
					ColMeta: []flux.ColMeta{
						{Label: "_measurement", Type: flux.TString},
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
						{Label: "domain", Type: flux.TString},
						{Label: "host", Type: flux.TString},
					},
					Data: [][]interface{}{
						{"m", execute.Time(2), 11.0, "www.example.com", "A.www.example.com"},
					},
				},
			},
		},
		{
			name: `with _value+5 with nulls`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({r with _value: r._value + 5.0})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				KeyCols: []string{"_measurement", "host"},
				ColMeta: []flux.ColMeta{
					{Label: "_measurement", Type: flux.TString},
					{Label: "host", Type: flux.TString},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{"m", nil, execute.Time(1), 1.0},
					{"m", nil, execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{{
				KeyCols: []string{"_measurement", "host"},
				ColMeta: []flux.ColMeta{
					{Label: "_measurement", Type: flux.TString},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "host", Type: flux.TString},
				},
				Data: [][]interface{}{
					{"m", execute.Time(1), 6.0, nil},
					{"m", execute.Time(2), 11.0, nil},
				},
			}},
		},
		{
			// Parameter mergeKey is now deprecated and has no effect.
			// `map` is evaluated as if it was false.
			name: `_value+5 mergeKey=true regroup`,
			spec: &universe.MapProcedureSpec{
				MergeKey: true,
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_time: r._time, host: r.host + ".local", _value: r._value + 5.0})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				KeyCols: []string{"_measurement", "host"},
				ColMeta: []flux.ColMeta{
					{Label: "_measurement", Type: flux.TString},
					{Label: "host", Type: flux.TString},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{"m", "A", execute.Time(1), 1.0},
					{"m", "A", execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{{
				KeyCols: []string{"host"},
				ColMeta: []flux.ColMeta{
					{Label: "host", Type: flux.TString},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{"A.local", execute.Time(1), 6.0},
					{"A.local", execute.Time(2), 11.0},
				},
			}},
		},
		{
			// Parameter mergeKey is now deprecated and has no effect.
			// `map` is evaluated as if it was false.
			name: `_value+5 mergeKey=true regroup fan out`,
			spec: &universe.MapProcedureSpec{
				MergeKey: true,
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_time: r._time, host: r.host + "." + r.domain, _value: r._value + 5.0})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				KeyCols: []string{"_measurement", "host"},
				ColMeta: []flux.ColMeta{
					{Label: "_measurement", Type: flux.TString},
					{Label: "host", Type: flux.TString},
					{Label: "domain", Type: flux.TString},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{"m", "A", "example.com", execute.Time(1), 1.0},
					{"m", "A", "www.example.com", execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{
				{
					KeyCols: []string{"host"},
					ColMeta: []flux.ColMeta{
						{Label: "host", Type: flux.TString},
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{"A.example.com", execute.Time(1), 6.0},
					},
				},
				{
					KeyCols: []string{"host"},
					ColMeta: []flux.ColMeta{
						{Label: "host", Type: flux.TString},
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{"A.www.example.com", execute.Time(2), 11.0},
					},
				},
			},
		},
		{
			// Parameter mergeKey is now deprecated and has no effect.
			// `map` is evaluated as if it was false.
			name: `_value+5 mergeKey=true with nulls`,
			spec: &universe.MapProcedureSpec{
				MergeKey: true,
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_time: r._time, _value: r._value + 5.0, })`),
				},
			},
			data: []flux.Table{&executetest.Table{
				KeyCols: []string{"_measurement", "host"},
				ColMeta: []flux.ColMeta{
					{Label: "_measurement", Type: flux.TString},
					{Label: "host", Type: flux.TString},
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{"m", nil, execute.Time(1), 1.0},
					{"m", nil, execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 6.0},
					{execute.Time(2), 11.0},
				},
			}},
		},
		{
			name: `_value*_value`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_time: r._time, _value: r._value * r._value})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 1.0},
					{execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 1.0},
					{execute.Time(2), 36.0},
				},
			}},
		},
		{
			name: "float(r._value) int",
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_time: r._time, _value: float(v: r._value)})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), int64(1)},
					{execute.Time(2), int64(6)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 1.0},
					{execute.Time(2), 6.0},
				},
			}},
		},
		{
			name: "float(r._value) uint",
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_time: r._time, _value: float(v: r._value)})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TUInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), uint64(1)},
					{execute.Time(2), uint64(6)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 1.0},
					{execute.Time(2), 6.0},
				},
			}},
		},
		{
			name: `error returning array property`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_value: ["str"]})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TUInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), uint64(1)},
				},
			}},
			wantErr: errors.New(`map object property "_value" is array type which is not supported in a flux table`),
		},
		{
			name: `error returning regexp property`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_value: /ab?/})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TUInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), uint64(1)},
				},
			}},
			wantErr: errors.New(`map object property "_value" is regexp type which is not supported in a flux table`),
		},
		{
			name: `error returning function property`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_value: (p) => 1})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TUInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), uint64(1)},
				},
			}},
			wantErr: errors.New(`map object property "_value" is function type which is not supported in a flux table`),
		},
		{
			name: `error returning duration property`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_value: 1d})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TUInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), uint64(1)},
				},
			}},
			wantErr: errors.New(`map object property "_value" is duration type which is not supported in a flux table`),
		},
		{
			name: `error returning object property`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_value: {v: 1}})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TUInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), uint64(1)},
				},
			}},
			wantErr: errors.New(`map object property "_value" is object type which is not supported in a flux table`),
		},
		{
			name: `error returning byte property`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_value: bytes(v: "123")})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TUInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), uint64(1)},
				},
			}},
			wantErr: errors.New(`map object property "_value" is bytes type which is not supported in a flux table`),
		},
		{
			name: `float("foo") produces error`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_time: r._time, _value: float(v: "foo")})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TUInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), uint64(1)},
				},
			}},
			wantErr: errors.New(`failed to evaluate map function: cannot convert string "foo" to float due to invalid syntax`),
		},
		{
			name: `with null record`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({value: r._value + 5.0})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), nil},
					{execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{nil},
					{11.0},
				},
			}},
		},
		{
			name: `with null column`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({value: r._value, missing: r._missing})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 1.0},
					{execute.Time(2), 6.0},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{1.0},
					{6.0},
				},
			}},
		},
		{
			name: `scoped labels different types`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({r with _value: float(v: r._value + 1)})`),
				},
			},
			data: []flux.Table{&executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TInt},
				},
				Data: [][]interface{}{
					{execute.Time(1), int64(1)},
					{execute.Time(2), int64(6)},
				},
			}},
			want: []*executetest.Table{{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
				},
				Data: [][]interface{}{
					{execute.Time(1), 2.0},
					{execute.Time(2), 7.0},
				},
			}},
		},
		{
			name: `mismatched types`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_value: r._value, _time: r._time})`),
				},
			},
			data: []flux.Table{
				&executetest.Table{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TInt},
					},
					Data: [][]interface{}{
						{execute.Time(1), int64(1)},
						{execute.Time(2), int64(6)},
					},
				},
				&executetest.Table{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 6.0},
					},
				},
			},
			wantErr: errors.New(`schema collision detected: column "_value" is both of type float and int`),
		},
		{
			name: `mismatched types with zero value`,
			spec: &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Scope: builtIns,
					Fn:    executetest.FunctionExpression(t, `(r) => ({_value: r._value, _time: r._time})`),
				},
			},
			data: []flux.Table{
				&executetest.Table{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TString},
					},
					Data: [][]interface{}{
						{execute.Time(1), ""},
						{execute.Time(2), "some words"},
					},
				},
				&executetest.Table{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 0.0},
					},
				},
			},
			wantErr: errors.New(`schema collision detected: column "_value" is both of type float and string`),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			executetest.ProcessTestHelper2(
				t,
				tc.data,
				tc.want,
				tc.wantErr,
				func(id execute.DatasetID, alloc memory.Allocator) (execute.Transformation, execute.Dataset) {
					ctx, deps := dependency.Inject(context.Background(), dependenciestest.Default())
					defer deps.Finish()
					xform, dataset, err := universe.NewMapTransformation(ctx, id, tc.spec, alloc)
					if err != nil {
						t.Fatal(err)
					}
					return xform, dataset
				},
			)
		})
	}
}

func BenchmarkMap_Process(b *testing.B) {
	genSource := func(alloc memory.Allocator) (flux.TableIterator, error) {
		return gen.Input(context.Background(), gen.Schema{
			Tags: []gen.Tag{
				{Name: "t0", Cardinality: 1},
				{Name: "t1", Cardinality: 1},
			},
			NumPoints: 1000,
			Alloc:     alloc,
		})
	}

	genTransformation := func(fn func(b *testing.B) *semantic.FunctionExpression) func(id execute.DatasetID, alloc memory.Allocator) (execute.Transformation, execute.Dataset) {
		return func(id execute.DatasetID, alloc memory.Allocator) (execute.Transformation, execute.Dataset) {
			spec := &universe.MapProcedureSpec{
				Fn: interpreter.ResolvedFunction{
					Fn:    fn(b),
					Scope: valuestest.Scope(),
				},
			}

			tr, d, err := universe.NewMapTransformation(context.Background(), id, spec, alloc)
			if err != nil {
				b.Fatal(err)
			}
			return tr, d
		}
	}

	const source = `(r) => ({r with v: r._value})`
	b.Run("Row", func(b *testing.B) {
		executetest.ProcessBenchmarkHelper(b,
			genSource,
			genTransformation(func(b *testing.B) *semantic.FunctionExpression {
				fn := executetest.FunctionExpression(b, source)
				fn.Vectorized = nil
				return fn
			}),
		)
	})

	b.Run("Vectorized", func(b *testing.B) {
		executetest.ProcessBenchmarkHelper(b,
			genSource,
			genTransformation(func(b *testing.B) *semantic.FunctionExpression {
				return executetest.FunctionExpression(b, source)
			}),
		)
	})
}
