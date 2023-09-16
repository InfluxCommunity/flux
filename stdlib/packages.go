// DO NOT EDIT: This file is autogenerated via the builtin command.
//
// The imports in this file ensures that all the init functions runs and registers
// the builtins for the flux runtime

package stdlib

import (
	_ "github.com/InfluxCommunity/flux/stdlib/array"
	_ "github.com/InfluxCommunity/flux/stdlib/bitwise"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/RohanSreerama5/naiveBayesClassifier"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/anaisdg/anomalydetection"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/anaisdg/statsmodels"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/bonitoo-io/alerta"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/bonitoo-io/hex"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/bonitoo-io/servicenow"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/bonitoo-io/tickscript"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/bonitoo-io/victorops"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/bonitoo-io/zenoss"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/chobbs/discord"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/jsternberg/influxdb"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/qxip/clickhouse"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/qxip/hash"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/qxip/iox"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/qxip/logql"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/rhajek/bigpanda"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/sranka/opsgenie"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/sranka/sensu"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/sranka/teams"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/sranka/telegram"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/sranka/webexteams"
	_ "github.com/InfluxCommunity/flux/stdlib/contrib/tomhollingworth/events"
	_ "github.com/InfluxCommunity/flux/stdlib/csv"
	_ "github.com/InfluxCommunity/flux/stdlib/date"
	_ "github.com/InfluxCommunity/flux/stdlib/date/boundaries"
	_ "github.com/InfluxCommunity/flux/stdlib/dict"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/aggregate"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/array"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/bigtable"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/bitwise"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/csv"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/date/boundaries"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/dynamic"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/geo"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/http"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/http/requests"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/influxdb"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/iox"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/json"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/mqtt"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/oee"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/polyline"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/prometheus"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/query"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/record"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/table"
	_ "github.com/InfluxCommunity/flux/stdlib/experimental/usage"
	_ "github.com/InfluxCommunity/flux/stdlib/generate"
	_ "github.com/InfluxCommunity/flux/stdlib/http"
	_ "github.com/InfluxCommunity/flux/stdlib/http/requests"
	_ "github.com/InfluxCommunity/flux/stdlib/influxdata/influxdb"
	_ "github.com/InfluxCommunity/flux/stdlib/influxdata/influxdb/monitor"
	_ "github.com/InfluxCommunity/flux/stdlib/influxdata/influxdb/sample"
	_ "github.com/InfluxCommunity/flux/stdlib/influxdata/influxdb/schema"
	_ "github.com/InfluxCommunity/flux/stdlib/influxdata/influxdb/secrets"
	_ "github.com/InfluxCommunity/flux/stdlib/influxdata/influxdb/tasks"
	_ "github.com/InfluxCommunity/flux/stdlib/influxdata/influxdb/v1"
	_ "github.com/InfluxCommunity/flux/stdlib/internal/boolean"
	_ "github.com/InfluxCommunity/flux/stdlib/internal/debug"
	_ "github.com/InfluxCommunity/flux/stdlib/internal/gen"
	_ "github.com/InfluxCommunity/flux/stdlib/internal/influxql"
	_ "github.com/InfluxCommunity/flux/stdlib/internal/location"
	_ "github.com/InfluxCommunity/flux/stdlib/internal/promql"
	_ "github.com/InfluxCommunity/flux/stdlib/internal/testing"
	_ "github.com/InfluxCommunity/flux/stdlib/internal/testutil"
	_ "github.com/InfluxCommunity/flux/stdlib/interpolate"
	_ "github.com/InfluxCommunity/flux/stdlib/join"
	_ "github.com/InfluxCommunity/flux/stdlib/json"
	_ "github.com/InfluxCommunity/flux/stdlib/kafka"
	_ "github.com/InfluxCommunity/flux/stdlib/math"
	_ "github.com/InfluxCommunity/flux/stdlib/pagerduty"
	_ "github.com/InfluxCommunity/flux/stdlib/planner"
	_ "github.com/InfluxCommunity/flux/stdlib/profiler"
	_ "github.com/InfluxCommunity/flux/stdlib/pushbullet"
	_ "github.com/InfluxCommunity/flux/stdlib/regexp"
	_ "github.com/InfluxCommunity/flux/stdlib/runtime"
	_ "github.com/InfluxCommunity/flux/stdlib/sampledata"
	_ "github.com/InfluxCommunity/flux/stdlib/slack"
	_ "github.com/InfluxCommunity/flux/stdlib/socket"
	_ "github.com/InfluxCommunity/flux/stdlib/sql"
	_ "github.com/InfluxCommunity/flux/stdlib/strings"
	_ "github.com/InfluxCommunity/flux/stdlib/system"
	_ "github.com/InfluxCommunity/flux/stdlib/testing"
	_ "github.com/InfluxCommunity/flux/stdlib/testing/basics"
	_ "github.com/InfluxCommunity/flux/stdlib/testing/chronograf"
	_ "github.com/InfluxCommunity/flux/stdlib/testing/expect"
	_ "github.com/InfluxCommunity/flux/stdlib/testing/influxql"
	_ "github.com/InfluxCommunity/flux/stdlib/testing/kapacitor"
	_ "github.com/InfluxCommunity/flux/stdlib/testing/pandas"
	_ "github.com/InfluxCommunity/flux/stdlib/testing/prometheus"
	_ "github.com/InfluxCommunity/flux/stdlib/testing/promql"
	_ "github.com/InfluxCommunity/flux/stdlib/testing/usage"
	_ "github.com/InfluxCommunity/flux/stdlib/timezone"
	_ "github.com/InfluxCommunity/flux/stdlib/types"
	_ "github.com/InfluxCommunity/flux/stdlib/universe"
)
