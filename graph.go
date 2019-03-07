package top

import (
	"errors"
)

var ErrCyclicGraph = errors.New("cyclic graph")

type Graph struct {
	nodes map[interface{}]*node
}

type mark byte

const (
	unmarked mark = iota
	temporary
	permanent
)

type node struct {
	of       interface{}
	mark     mark
	pointsTo []*node
}

func (g *Graph) Link(from, to interface{}) {
	fromN := g.nodeFor(from)
	toN := g.nodeFor(to)
	fromN.linkTo(toN)
}

func (g *Graph) Sort() ([]interface{}, error) {
	sorted := make([]interface{}, 0, len(g.nodes))
	for _, n := range g.nodes {
		if n.mark != unmarked {
			continue
		}
		if err := n.visit(&sorted); err != nil {
			return nil, err
		}
	}
	return reverse(sorted), nil
}

func (g *Graph) nodeFor(x interface{}) *node {
	if g.nodes == nil {
		g.nodes = map[interface{}]*node{}
	}
	n := g.nodes[x]
	if n == nil {
		n = &node{of: x}
		g.nodes[x] = n
	}
	return n
}

func (n *node) linkTo(m *node) {
	for _, l := range n.pointsTo {
		if l == m {
			return
		}
	}
	n.pointsTo = append(n.pointsTo, m)
}

func (n *node) visit(sorted *[]interface{}) error {
	if n.mark == permanent {
		return nil
	}
	if n.mark == temporary {
		return ErrCyclicGraph
	}
	n.mark = temporary
	for _, n := range n.pointsTo {
		if err := n.visit(sorted); err != nil {
			return err
		}
	}
	n.mark = permanent
	*sorted = append(*sorted, n.of)
	return nil
}

func reverse(xs []interface{}) []interface{} {
	l := len(xs)
	for i := 0; i < l/2; i++ {
		xs[i], xs[l-i-1] = xs[l-i-1], xs[i]
	}
	return xs
}
