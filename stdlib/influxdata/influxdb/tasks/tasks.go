package tasks

import (
	"context"

	"github.com/InfluxCommunity/flux/codes"
	"github.com/InfluxCommunity/flux/execute"
	"github.com/InfluxCommunity/flux/internal/errors"
	"github.com/InfluxCommunity/flux/interpreter"
	"github.com/InfluxCommunity/flux/runtime"
	"github.com/InfluxCommunity/flux/semantic"
	"github.com/InfluxCommunity/flux/values"
)

const (
	pkgpath             = "influxdata/influxdb/tasks"
	zeroTimeName        = "_zeroTime"
	lastSuccessFuncName = "_lastSuccess"
)

func init() {
	runtime.RegisterPackageValue(pkgpath, lastSuccessFuncName, LastSuccessFunction)
	runtime.RegisterPackageValue(pkgpath, zeroTimeName, values.Null)
}

// LastSuccessFunction is a function that calls LastSuccess.
var LastSuccessFunction = makeLastSuccessFunc()

func makeLastSuccessFunc() values.Function {
	sig := runtime.MustLookupBuiltinType(pkgpath, lastSuccessFuncName)
	return values.NewFunction("lastSuccess", sig, LastSuccess, false)
}

// LastSuccess retrieves the last successful run of the task, or returns the value of the
// orTime parameter if the task has never successfully run.
func LastSuccess(ctx context.Context, args values.Object) (values.Value, error) {
	return interpreter.DoFunctionCallContext(func(ctx context.Context, args interpreter.Arguments) (values.Value, error) {
		orTime, err := args.GetRequired("orTime")
		if err != nil {
			return nil, err
		} else if !values.IsTimeable(orTime) {
			return nil, errors.Newf(codes.Invalid, "keyword argument \"orTime\" should be a time or duration, but got %v", orTime.Type().Nature())
		}

		if orTime.Type().Nature() == semantic.Duration {
			deps := execute.GetExecutionDependencies(ctx)
			orTime = values.NewTime(values.ConvertTime(*deps.Now).Add(orTime.Duration()))
		}

		lastSuccess, err := args.GetRequired("lastSuccessTime")
		if err != nil {
			return nil, err
		}

		// If the last success time is null, do not check its nature
		// as the nature isn't valid. Return the orTime.
		if lastSuccess.IsNull() {
			return orTime, nil
		}

		// We are going to return the lastSuccess time so verify that it
		// is the correct type.
		if kind := semantic.Time; lastSuccess.Type().Nature() != kind {
			return nil, errors.Newf(codes.Invalid, "keyword argument \"lastSuccessTime\" should be of kind %v, but got %v", kind, lastSuccess.Type().Nature())
		}
		return lastSuccess, nil
	}, ctx, args)
}
