package top

import (
	"reflect"
	"testing"
)

func TestGraph(t *testing.T) {
	type pairs []struct{ from, to string }
	type result []interface{}
	for _, test := range []struct {
		name   string
		pairs  pairs
		result result
		err    error
	}{
		{
			"Simple",
			pairs{
				{"A", "B"},
			},
			result{"A", "B"},
			nil,
		},
		{
			"Transitive",
			pairs{
				{"A", "B"},
				{"B", "C"},
			},
			result{"A", "B", "C"},
			nil,
		},
		{
			"Cycle",
			pairs{
				{"A", "B"},
				{"B", "A"},
			},
			nil,
			ErrCyclicGraph,
		},
		{
			"CycleTrans",
			pairs{
				{"A", "B"},
				{"B", "C"},
				{"C", "A"},
			},
			nil,
			ErrCyclicGraph,
		},
		{
			"Stutter",
			pairs{
				{"A", "B"},
				{"A", "B"},
			},
			result{"A", "B"},
			nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			g := new(Graph)
			for _, p := range test.pairs {
				g.Link(p.from, p.to)
			}
			out, err := g.Sort()
			if err != test.err {
				t.Errorf("unexpected error, expected %v, got %v", test.err, err)
			}
			if !reflect.DeepEqual(result(out), test.result) {
				t.Errorf("expected %v, got %v", test.result, out)
			}
		})
	}
}
