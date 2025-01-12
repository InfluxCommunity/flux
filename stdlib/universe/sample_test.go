package universe_test

import (
	"testing"

	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/execute"
	"github.com/InfluxCommunity/flux/execute/executetest"
	"github.com/InfluxCommunity/flux/stdlib/universe"
)

func TestSample_Process(t *testing.T) {
	testCases := []struct {
		name   string
		data   flux.Table
		want   [][]int
		fromor *universe.SampleSelector
	}{
		{
			fromor: &universe.SampleSelector{
				N:   1,
				Pos: 0,
			},
			name: "everything in separate Do calls",
			data: &executetest.RowWiseTable{
				Table: &executetest.Table{
					KeyCols: []string{"t1"},
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
						{Label: "t1", Type: flux.TString},
						{Label: "t2", Type: flux.TString},
					},
					Data: [][]interface{}{
						{execute.Time(0), 7.0, "a", "y"},
						{execute.Time(10), 5.0, "a", "x"},
						{execute.Time(20), 9.0, "a", "y"},
						{execute.Time(30), 4.0, "a", "x"},
						{execute.Time(40), 6.0, "a", "y"},
						{execute.Time(50), 8.0, "a", "x"},
						{execute.Time(60), 1.0, "a", "y"},
						{execute.Time(70), 2.0, "a", "x"},
						{execute.Time(80), 3.0, "a", "y"},
						{execute.Time(90), 10.0, "a", "x"},
					},
				},
			},
			want: [][]int{
				{0},
				{0},
				{0},
				{0},
				{0},
				{0},
				{0},
				{0},
				{0},
				{0},
			},
		},
		{
			fromor: &universe.SampleSelector{
				N:   1,
				Pos: 0,
			},
			name: "everything in single Do call",
			data: executetest.MustCopyTable(&executetest.Table{
				KeyCols: []string{"t1"},
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "t1", Type: flux.TString},
					{Label: "t2", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(0), 7.0, "a", "y"},
					{execute.Time(10), 5.0, "a", "x"},
					{execute.Time(20), 9.0, "a", "y"},
					{execute.Time(30), 4.0, "a", "x"},
					{execute.Time(40), 6.0, "a", "y"},
					{execute.Time(50), 8.0, "a", "x"},
					{execute.Time(60), 1.0, "a", "y"},
					{execute.Time(70), 2.0, "a", "x"},
					{execute.Time(80), 3.0, "a", "y"},
					{execute.Time(90), 10.0, "a", "x"},
				},
			}),
			want: [][]int{{
				0,
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8,
				9,
			}},
		},
		{
			fromor: &universe.SampleSelector{
				N:   2,
				Pos: 0,
			},
			name: "every-other-even",
			data: executetest.MustCopyTable(&executetest.Table{
				KeyCols: []string{"t1"},
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "t1", Type: flux.TString},
					{Label: "t2", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(0), 7.0, "a", "y"},
					{execute.Time(10), 5.0, "a", "x"},
					{execute.Time(20), 9.0, "a", "y"},
					{execute.Time(30), 4.0, "a", "x"},
					{execute.Time(40), 6.0, "a", "y"},
					{execute.Time(50), 8.0, "a", "x"},
					{execute.Time(60), 1.0, "a", "y"},
					{execute.Time(70), 2.0, "a", "x"},
					{execute.Time(80), 3.0, "a", "y"},
					{execute.Time(90), 10.0, "a", "x"},
				},
			}),
			want: [][]int{{
				0,
				2,
				4,
				6,
				8,
			}},
		},
		{
			fromor: &universe.SampleSelector{
				N:   2,
				Pos: 1,
			},
			name: "every-other-odd",
			data: executetest.MustCopyTable(&executetest.Table{
				KeyCols: []string{"t1"},
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "t1", Type: flux.TString},
					{Label: "t2", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(0), 7.0, "a", "y"},
					{execute.Time(10), 5.0, "a", "x"},
					{execute.Time(20), 9.0, "a", "y"},
					{execute.Time(30), 4.0, "a", "x"},
					{execute.Time(40), 6.0, "a", "y"},
					{execute.Time(50), 8.0, "a", "x"},
					{execute.Time(60), 1.0, "a", "y"},
					{execute.Time(70), 2.0, "a", "x"},
					{execute.Time(80), 3.0, "a", "y"},
					{execute.Time(90), 10.0, "a", "x"},
				},
			}),
			want: [][]int{{
				1,
				3,
				5,
				7,
				9,
			}},
		},
		{
			fromor: &universe.SampleSelector{
				N:   3,
				Pos: 0,
			},
			name: "every-third-0",
			data: executetest.MustCopyTable(&executetest.Table{
				KeyCols: []string{"t1"},
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "t1", Type: flux.TString},
					{Label: "t2", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(0), 7.0, "a", "y"},
					{execute.Time(10), 5.0, "a", "x"},
					{execute.Time(20), 9.0, "a", "y"},
					{execute.Time(30), 4.0, "a", "x"},
					{execute.Time(40), 6.0, "a", "y"},
					{execute.Time(50), 8.0, "a", "x"},
					{execute.Time(60), 1.0, "a", "y"},
					{execute.Time(70), 2.0, "a", "x"},
					{execute.Time(80), 3.0, "a", "y"},
					{execute.Time(90), 10.0, "a", "x"},
				},
			}),
			want: [][]int{{
				0,
				3,
				6,
				9,
			}},
		},
		{
			fromor: &universe.SampleSelector{
				N:   3,
				Pos: 1,
			},
			name: "every-third-1",
			data: executetest.MustCopyTable(&executetest.Table{
				KeyCols: []string{"t1"},
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "t1", Type: flux.TString},
					{Label: "t2", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(0), 7.0, "a", "y"},
					{execute.Time(10), 5.0, "a", "x"},
					{execute.Time(20), 9.0, "a", "y"},
					{execute.Time(30), 4.0, "a", "x"},
					{execute.Time(40), 6.0, "a", "y"},
					{execute.Time(50), 8.0, "a", "x"},
					{execute.Time(60), 1.0, "a", "y"},
					{execute.Time(70), 2.0, "a", "x"},
					{execute.Time(80), 3.0, "a", "y"},
					{execute.Time(90), 10.0, "a", "x"},
				},
			}),
			want: [][]int{{
				1,
				4,
				7,
			}},
		},
		{
			fromor: &universe.SampleSelector{
				N:   3,
				Pos: 2,
			},
			name: "every-third-2",
			data: executetest.MustCopyTable(&executetest.Table{
				KeyCols: []string{"t1"},
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "t1", Type: flux.TString},
					{Label: "t2", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(0), 7.0, "a", "y"},
					{execute.Time(10), 5.0, "a", "x"},
					{execute.Time(20), 9.0, "a", "y"},
					{execute.Time(30), 4.0, "a", "x"},
					{execute.Time(40), 6.0, "a", "y"},
					{execute.Time(50), 8.0, "a", "x"},
					{execute.Time(60), 1.0, "a", "y"},
					{execute.Time(70), 2.0, "a", "x"},
					{execute.Time(80), 3.0, "a", "y"},
					{execute.Time(90), 10.0, "a", "x"},
				},
			}),
			want: [][]int{{
				2,
				5,
				8,
			}},
		},
		{
			fromor: &universe.SampleSelector{
				N:   3,
				Pos: 2,
			},
			name: "every-third-2 in separate Do calls",
			data: &executetest.RowWiseTable{
				Table: &executetest.Table{
					KeyCols: []string{"t1"},
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
						{Label: "t1", Type: flux.TString},
						{Label: "t2", Type: flux.TString},
					},
					Data: [][]interface{}{
						{execute.Time(0), 7.0, "a", "y"},
						{execute.Time(10), 5.0, "a", "x"},
						{execute.Time(20), 9.0, "a", "y"},
						{execute.Time(30), 4.0, "a", "x"},
						{execute.Time(40), 6.0, "a", "y"},
						{execute.Time(50), 8.0, "a", "x"},
						{execute.Time(60), 1.0, "a", "y"},
						{execute.Time(70), 2.0, "a", "x"},
						{execute.Time(80), 3.0, "a", "y"},
						{execute.Time(90), 10.0, "a", "x"},
					},
				},
			},
			want: [][]int{
				nil,
				nil,
				{0},
				nil,
				nil,
				{0},
				nil,
				nil,
				{0},
				nil,
			},
		},
		{
			fromor: &universe.SampleSelector{
				N:   2,
				Pos: 0,
			},
			name: "with-nulls",
			data: executetest.MustCopyTable(&executetest.Table{
				KeyCols: []string{"t1"},
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "t1", Type: flux.TString},
					{Label: "t2", Type: flux.TString},
					{Label: "t3", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(0), 7.0, "a", "y", nil},
					{execute.Time(10), 5.0, "a", "x", nil},
					{execute.Time(20), 9.0, "a", "y", nil},
					{execute.Time(30), 4.0, "a", "x", nil},
					{execute.Time(40), 6.0, "a", "y", nil},
					{execute.Time(50), 8.0, "a", "x", nil},
					{execute.Time(60), 1.0, "a", "y", nil},
					{execute.Time(70), 2.0, "a", "x", nil},
					{execute.Time(80), 3.0, "a", "y", nil},
					{execute.Time(90), 10.0, "a", "x", nil},
				},
			}),
			want: [][]int{{
				0,
				2,
				4,
				6,
				8,
			}},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			executetest.IndexSelectorFuncTestHelper(
				t,
				tc.fromor,
				tc.data,
				tc.want,
			)
		})
	}
}

func BenchmarkSample(b *testing.B) {
	ss := &universe.SampleSelector{
		N:   10,
		Pos: 0,
	}
	executetest.IndexSelectorFuncBenchmarkHelper(b, ss, NormalTable)
}
