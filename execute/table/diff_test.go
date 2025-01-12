package table_test

import (
	"testing"

	"github.com/InfluxCommunity/flux/execute/table"
	"github.com/InfluxCommunity/flux/execute/table/static"
)

func TestDiff(t *testing.T) {
	wantI := static.Table{
		static.StringKey("_measurement", "m0"),
		static.StringKey("_field", "f0"),
		static.StringKey("t0", "a"),
		static.Times("_time", "2020-01-01T00:00:00Z", 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110),
		static.Ints("_value", 6, 7, 4, 12, 3, 9, 5, 6, 5, 1, 8, 4),
	}
	gotI := static.Table{
		static.StringKey("_measurement", "m0"),
		static.StringKey("_field", "f0"),
		static.StringKey("t0", "a"),
		static.Times("_time", "2020-01-01T00:00:00Z", 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110),
		static.Ints("_value", 6, 7, 3, 12, 3, 9, 5, 6, 5, 1, 8, 4),
	}
	got := table.Diff(wantI, gotI)

	want := ` # _field=f0,_measurement=m0,t0=a _time=time,_value=int
 _field=f0,_measurement=m0,t0=a _time=2020-01-01T00:00:00Z,_value=6i
 _field=f0,_measurement=m0,t0=a _time=2020-01-01T00:00:10Z,_value=7i
-_field=f0,_measurement=m0,t0=a _time=2020-01-01T00:00:20Z,_value=4i
+_field=f0,_measurement=m0,t0=a _time=2020-01-01T00:00:20Z,_value=3i
 _field=f0,_measurement=m0,t0=a _time=2020-01-01T00:00:30Z,_value=12i
 _field=f0,_measurement=m0,t0=a _time=2020-01-01T00:00:40Z,_value=3i
 _field=f0,_measurement=m0,t0=a _time=2020-01-01T00:00:50Z,_value=9i`
	if got != want {
		t.Errorf("unexpected diff output -want/+got:\n\t- %q\n\t+ %q", want, got)
	}
}
