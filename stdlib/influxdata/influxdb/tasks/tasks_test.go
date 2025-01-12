package tasks_test

import (
	"context"
	"testing"
	"time"

	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/execute"
	"github.com/InfluxCommunity/flux/stdlib/influxdata/influxdb/tasks"
	"github.com/InfluxCommunity/flux/values"
	"github.com/google/go-cmp/cmp"
)

func TestLastSuccess(t *testing.T) {
	for _, tt := range []struct {
		name string
		args map[string]values.Value
		now  time.Time
		want values.Value
	}{
		{
			name: "orTime",
			args: map[string]values.Value{
				"orTime":          values.NewTime(10),
				"lastSuccessTime": values.Null,
			},
			want: values.NewTime(10),
		},
		{
			name: "lastSuccessTime",
			args: map[string]values.Value{
				"orTime":          values.NewTime(10),
				"lastSuccessTime": values.NewTime(5),
			},
			want: values.NewTime(5),
		},
		{
			name: "implied orTime",
			args: map[string]values.Value{
				"orTime":          values.NewDuration(flux.ConvertDuration(-5 * time.Minute)),
				"lastSuccessTime": values.NewTime(5),
			},
			want: values.NewTime(5),
		},
		{
			name: "implied orTime with null",
			args: map[string]values.Value{
				"orTime":          values.NewDuration(flux.ConvertDuration(-5 * time.Nanosecond)),
				"lastSuccessTime": values.Null,
			},
			now:  time.Unix(0, 10),
			want: values.NewTime(5),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			args := values.NewObjectWithValues(tt.args)
			deps := execute.DefaultExecutionDependencies()

			if !tt.now.IsZero() {
				deps = execute.NewExecutionDependencies(nil, &tt.now, nil)
			}

			got, err := tasks.LastSuccess(deps.Inject(context.Background()), args)
			if err != nil {
				t.Fatal(err)
			} else if !cmp.Equal(tt.want, got) {
				t.Fatalf("unexpected value -want/+got:\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}
