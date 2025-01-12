package regexp

import (
	"context"
	"regexp"

	"github.com/InfluxCommunity/flux/codes"
	"github.com/InfluxCommunity/flux/internal/errors"
	"github.com/InfluxCommunity/flux/runtime"
	"github.com/InfluxCommunity/flux/semantic"
	"github.com/InfluxCommunity/flux/values"
)

var SpecialFns map[string]values.Function

func init() {

	SpecialFns = map[string]values.Function{
		"compile": values.NewFunction(
			"compile",
			runtime.MustLookupBuiltinType("regexp", "compile"),
			func(ctx context.Context, args values.Object) (values.Value, error) {
				v, ok := args.Get("v")
				if !ok {
					return nil, errors.New(codes.Invalid, "missing argument v")
				}

				if !v.IsNull() && v.Type().Nature() == semantic.String {
					re, err := regexp.Compile(v.Str())
					if err != nil {
						return nil, err
					}
					return values.NewRegexp(re), err
				}
				return nil, errors.Newf(codes.Invalid, "cannot convert argument v of type %v value %v to Regex", v.Type().Nature(), v)
			},
			false,
		),
		"quoteMeta": values.NewFunction(
			"quoteMeta",
			runtime.MustLookupBuiltinType("regexp", "quoteMeta"),
			func(ctx context.Context, args values.Object) (values.Value, error) {
				v, ok := args.Get("v")
				if !ok {
					return nil, errors.New(codes.Invalid, "missing argument v")
				}

				if !v.IsNull() && v.Type().Nature() == semantic.String {
					value := regexp.QuoteMeta(v.Str())
					return values.NewString(value), nil
				}
				return nil, errors.Newf(codes.Invalid, "cannot escape all regular expression metacharacters inside argument v of type %v value %v", v.Type().Nature(), v)
			},
			false,
		),
		"findString": values.NewFunction(
			"findString",
			runtime.MustLookupBuiltinType("regexp", "findString"),
			func(ctx context.Context, args values.Object) (values.Value, error) {
				v, ok := args.Get("v")
				r, okk := args.Get("r")
				if !ok || !okk {
					return nil, errors.New(codes.Invalid, "missing argument")
				}

				if !v.IsNull() && !r.IsNull() && v.Type().Nature() == semantic.String && r.Type().Nature() == semantic.Regexp {
					value := r.Regexp().FindString(v.Str())
					return values.NewString(value), nil
				}
				return nil, errors.Newf(codes.Invalid, "cannot execute function containing argument r of type %v value %v and argument v of type %v value %v", r.Type().Nature(), r, v.Type().Nature(), v)
			},
			false,
		),
		"findStringIndex": values.NewFunction(
			"findStringIndex",
			runtime.MustLookupBuiltinType("regexp", "findStringIndex"),
			func(ctx context.Context, args values.Object) (values.Value, error) {
				v, ok := args.Get("v")
				r, okk := args.Get("r")
				if !ok || !okk {
					return nil, errors.New(codes.Invalid, "missing argument")
				}

				if !v.IsNull() && !r.IsNull() && v.Type().Nature() == semantic.String && r.Type().Nature() == semantic.Regexp {
					value := r.Regexp().FindStringIndex(v.Str())
					arr := values.NewArray(semantic.NewArrayType(semantic.BasicInt))
					for _, z := range value {
						arr.Append(values.NewInt(int64(z)))
					}
					return arr, nil
				}
				return nil, errors.Newf(codes.Invalid, "cannot execute function containing argument r of type %v value %v and argument v of type %v value %v", r.Type().Nature(), r, v.Type().Nature(), v)
			},
			false,
		),
		"matchRegexpString": values.NewFunction(
			"matchRegexpString",
			runtime.MustLookupBuiltinType("regexp", "matchRegexpString"),
			func(ctx context.Context, args values.Object) (values.Value, error) {
				v, ok := args.Get("v")
				r, okk := args.Get("r")
				if !ok || !okk {
					return nil, errors.New(codes.Invalid, "missing argument")
				}

				if !v.IsNull() && !r.IsNull() && v.Type().Nature() == semantic.String && r.Type().Nature() == semantic.Regexp {
					value := r.Regexp().MatchString(v.Str())
					return values.NewBool(value), nil
				}
				return nil, errors.Newf(codes.Invalid, "cannot execute function containing argument r of type %v value %v and argument v of type %v value %v", r.Type().Nature(), r, v.Type().Nature(), v)
			},
			false,
		),
		"replaceAllString": values.NewFunction(
			"replaceAllString",
			runtime.MustLookupBuiltinType("regexp", "replaceAllString"),
			func(ctx context.Context, args values.Object) (values.Value, error) {
				r, ok := args.Get("r")
				v, okk := args.Get("v")
				t, okkk := args.Get("t")
				if !ok || !okk || !okkk {
					return nil, errors.New(codes.Invalid, "missing argument")
				}

				if !v.IsNull() && !t.IsNull() && !r.IsNull() && v.Type().Nature() == semantic.String && t.Type().Nature() == semantic.String && r.Type().Nature() == semantic.Regexp {
					value := r.Regexp().ReplaceAllString(v.Str(), t.Str())
					return values.NewString(value), nil
				}
				return nil, errors.Newf(codes.Invalid, "cannot execute function containing argument r of type %v value %v, argument v of type %v value %v, and argument t of type %v value %v", r.Type().Nature(), r, v.Type().Nature(), v, t.Type().Nature(), t)
			},
			false,
		),
		"splitRegexp": values.NewFunction(
			"splitRegexp",
			runtime.MustLookupBuiltinType("regexp", "splitRegexp"),
			func(ctx context.Context, args values.Object) (values.Value, error) {
				r, ok := args.Get("r")
				v, okk := args.Get("v")
				i, okkk := args.Get("i")
				if !ok || !okk || !okkk {
					return nil, errors.New(codes.Invalid, "missing argument")
				}

				if !v.IsNull() && !i.IsNull() && !r.IsNull() && v.Type().Nature() == semantic.String && i.Type().Nature() == semantic.Int && r.Type().Nature() == semantic.Regexp {
					value := r.Regexp().Split(v.Str(), int(i.Int()))
					arr := values.NewArray(semantic.NewArrayType(semantic.BasicString))
					for _, z := range value {
						arr.Append(values.NewString(z))
					}
					return arr, nil
				}
				return nil, errors.Newf(codes.Invalid, "cannot execute function containing argument r of type %v value %v, argument v of type %v value %v, and argument i of type %v value %v", r.Type().Nature(), r, v.Type().Nature(), v, i.Type().Nature(), i)
			},
			false,
		),
		"getString": values.NewFunction(
			"getString",
			runtime.MustLookupBuiltinType("regexp", "getString"),
			func(ctx context.Context, args values.Object) (values.Value, error) {
				r, ok := args.Get("r")
				if !ok {
					return nil, errors.New(codes.Invalid, "missing argument")
				}

				if !r.IsNull() && r.Type().Nature() == semantic.Regexp {
					value := r.Regexp().String()
					return values.NewString(value), nil
				}
				return nil, errors.Newf(codes.Invalid, "cannot execute function containing argument r of type %v value %v", r.Type().Nature(), r)
			},
			false,
		),
	}

	runtime.RegisterPackageValue("regexp", "compile", SpecialFns["compile"])
	runtime.RegisterPackageValue("regexp", "quoteMeta", SpecialFns["quoteMeta"])
	runtime.RegisterPackageValue("regexp", "findString", SpecialFns["findString"])
	runtime.RegisterPackageValue("regexp", "findStringIndex", SpecialFns["findStringIndex"])
	runtime.RegisterPackageValue("regexp", "matchRegexpString", SpecialFns["matchRegexpString"])
	runtime.RegisterPackageValue("regexp", "replaceAllString", SpecialFns["replaceAllString"])
	runtime.RegisterPackageValue("regexp", "splitRegexp", SpecialFns["splitRegexp"])
	runtime.RegisterPackageValue("regexp", "getString", SpecialFns["getString"])
}
