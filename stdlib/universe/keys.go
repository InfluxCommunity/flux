package universe

import (
	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/arrow"
	"github.com/InfluxCommunity/flux/codes"
	"github.com/InfluxCommunity/flux/execute"
	"github.com/InfluxCommunity/flux/internal/errors"
	"github.com/InfluxCommunity/flux/plan"
	"github.com/InfluxCommunity/flux/runtime"
)

const KeysKind = "keys"

type KeysOpSpec struct {
	Column string `json:"column"`
}

func init() {
	keysSignature := runtime.MustLookupBuiltinType("universe", "keys")

	runtime.RegisterPackageValue("universe", KeysKind, flux.MustValue(flux.FunctionValue(KeysKind, createKeysOpSpec, keysSignature)))
	plan.RegisterProcedureSpec(KeysKind, newKeysProcedure, KeysKind)
	execute.RegisterTransformation(KeysKind, createKeysTransformation)
}

func createKeysOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	spec := new(KeysOpSpec)

	if col, found, err := args.GetString("column"); err != nil {
		return nil, err
	} else if found {
		spec.Column = col
	} else {
		spec.Column = execute.DefaultValueColLabel
	}

	return spec, nil
}

func (s *KeysOpSpec) Kind() flux.OperationKind {
	return KeysKind
}

type KeysProcedureSpec struct {
	plan.DefaultCost
	Column string
}

func newKeysProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*KeysOpSpec)
	if !ok {
		return nil, errors.Newf(codes.Internal, "invalid spec type %T", qs)
	}

	return &KeysProcedureSpec{
		Column: spec.Column,
	}, nil
}

func (s *KeysProcedureSpec) Kind() plan.ProcedureKind {
	return KeysKind
}

func (s *KeysProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(KeysProcedureSpec)
	*ns = *s
	return ns
}

// TriggerSpec implements plan.TriggerAwareProcedureSpec
func (s *KeysProcedureSpec) TriggerSpec() plan.TriggerSpec {
	return plan.NarrowTransformationTriggerSpec{}
}

func createKeysTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*KeysProcedureSpec)
	if !ok {
		return nil, nil, errors.Newf(codes.Internal, "invalid spec type %T", spec)
	}
	cache := execute.NewTableBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)
	t := NewKeysTransformation(d, cache, s)
	return t, d, nil
}

type keysTransformation struct {
	execute.ExecutionNode
	d     execute.Dataset
	cache execute.TableBuilderCache

	column string
}

func NewKeysTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec *KeysProcedureSpec) *keysTransformation {
	return &keysTransformation{
		d:      d,
		cache:  cache,
		column: spec.Column,
	}
}

func (t *keysTransformation) RetractTable(id execute.DatasetID, key flux.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *keysTransformation) Process(id execute.DatasetID, tbl flux.Table) error {
	builder, created := t.cache.TableBuilder(tbl.Key())
	if !created {
		return errors.Newf(codes.FailedPrecondition, "keys found duplicate table with key: %v", tbl.Key())
	}

	keys := make([]string, 0, len(tbl.Cols()))
	for _, c := range tbl.Key().Cols() {
		keys = append(keys, c.Label)
	}

	// Add the key to this table.
	if err := execute.AddTableKeyCols(tbl.Key(), builder); err != nil {
		return err
	}

	colIdx, err := builder.AddCol(flux.ColMeta{Label: t.column, Type: flux.TString})
	if err != nil {
		return err
	}

	// Append the key values repeatedly to the table.
	for i := 0; i < len(keys); i++ {
		if err := execute.AppendKeyValues(tbl.Key(), builder); err != nil {
			return err
		}
	}

	keysArrow := arrow.NewString(keys, nil)
	defer keysArrow.Release()
	if err := builder.AppendStrings(colIdx, keysArrow); err != nil {
		return err
	}

	// TODO: this is a hack
	return tbl.Do(func(flux.ColReader) error {
		return nil
	})
}

func (t *keysTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) error {
	return t.d.UpdateWatermark(mark)
}
func (t *keysTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}
func (t *keysTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}
