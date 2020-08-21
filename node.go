package go_kd_segment_tree

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"strings"
)

type TreeNode interface {
	Search(p Point) []interface{}
	Dumps(prefix string) string
}

type LeafNode struct {
	TreeNode
	Segments []*Segment
}

func (node *LeafNode) Search(p Point) []interface{} {
	if node == nil {
		return nil
	}
	if node.Segments != nil {
		var result = mapset.NewSet()
		for _, seg := range node.Segments {
			if seg.Rect.Contains(p) {
				result = result.Union(seg.Data)
			}
		}
		return result.ToSlice()
	}
	return nil
}

func (node *LeafNode) Dumps(prefix string) string {
	if node == nil {
		return ""
	}

	return fmt.Sprintf("%s -leaf:{size=%v}", prefix, len(node.Segments))
}

type BinaryNode struct {
	TreeNode

	Axis  int
	Level int
	Gini  float64

	Mid Measure
	Min Measure
	Max Measure

	Left  TreeNode
	Right TreeNode
}

func (node *BinaryNode) Search(p Point) []interface{} {
	if node == nil {
		return nil
	}

	x := p[node.Axis]
	if x.Smaller(node.Min) || x.Bigger(node.Max) {
		return nil
	}

	if x.Smaller(node.Mid) {
		if node.Left == nil {
			return nil
		}
		return node.Left.Search(p)
	} else {
		if node.Right == nil {
			return nil
		}
		return node.Right.Search(p)
	}
}

func (node *BinaryNode) Dumps(prefix string) string {
	if node == nil {
		return ""
	}

	return fmt.Sprintf("%s -bnode{axis:%d, gini:%v, mid:%v, min:%v, max:%v}\n%v\n%v\n", prefix,
		node.Axis, node.Gini, node.Mid, node.Min, node.Max,
		node.Left.Dumps(prefix+"    "), node.Right.Dumps(prefix+"    "))
}

type HashNode struct {
	TreeNode

	Axis  int
	Level int
	Gini  float64

	child map[Measure]TreeNode
}

func (node *HashNode) Search(p Point) []interface{} {
	if node == nil {
		return nil
	}

	x := p[node.Axis]
	if child, ok := node.child[x]; ok {
		return child.Search(p)
	}
	return nil
}

func (node *HashNode) Dumps(prefix string) string {
	if node == nil {
		return ""
	}

	var msgs []string
	msgs = append(msgs, fmt.Sprintf("%s -hnode{axis:%d, gini:%v}", prefix, node.Axis, node.Gini))
	for childKey, child := range node.child {
		msgs = append(msgs, child.Dumps(fmt.Sprintf("%v    %v:", prefix, childKey)))
	}
	return strings.Join(msgs, "\n")
}

func NewNode(segments []*Segment, level int, leafNodeMin int, miniJini float64) TreeNode {
	if len(segments) == 0 {
		return nil
	}

	if len(segments) < leafNodeMin || level <= 0 {
		return &LeafNode{
			Segments: MergeSegments(segments),
		}
	}

	axisSegments := NewSegmentBranch(segments, miniJini)
	if axisSegments == nil {
		return &LeafNode{
			Segments: MergeSegments(segments),
		}
	}

	switch segments[0].Rect[axisSegments.axis].(type) {
	case Interval:
		return &BinaryNode{
			Axis:  axisSegments.axis,
			Level: level,
			Gini:  axisSegments.gini,
			Mid:   axisSegments.midSeg.Rect[axisSegments.axis].(Interval)[0],
			Min:   axisSegments.min,
			Max:   axisSegments.max,
			Left:  NewNode(axisSegments.left, level-1, leafNodeMin, miniJini),
			Right: NewNode(axisSegments.right, level-1, leafNodeMin, miniJini),
		}
	case Measure:
		node := &HashNode{
			Axis:  axisSegments.axis,
			Level: level,
			Gini:  axisSegments.gini,
			child: make(map[Measure]TreeNode),
		}
		for childKey, childSegments := range axisSegments.hashSegments {
			node.child[childKey] = NewNode(childSegments, level-1, leafNodeMin, miniJini)
		}
	}

	return nil
}
