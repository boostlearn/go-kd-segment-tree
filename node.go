package go_kd_segment_tree

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
)

type TreeNode struct {
	Axis int

	Level int

	Mid Measure
	Min Measure
	Max Measure

	Left  *TreeNode
	Right *TreeNode

	Segments []*Segment
}

func NewNode(segments []*Segment, level int, leafDataSizeMin int) *TreeNode {
	if len(segments) == 0 {
		return nil
	}

	if len(segments) < leafDataSizeMin || level <= 0 {
		return &TreeNode{
			Segments: segments,
		}
	}

	axisSegments := NewSplitSegments(segments)
	return &TreeNode{
		Axis:  axisSegments.axis,
		Level:level,
		Mid:   axisSegments.midSeg.Rect[axisSegments.axis][0],
		Min:   axisSegments.min,
		Max:   axisSegments.max,
		Left:  NewNode(axisSegments.left, level - 1, leafDataSizeMin),
		Right: NewNode(axisSegments.right, level - 1, leafDataSizeMin),
	}
}

func (node *TreeNode) Search(p Point) []interface{} {
	if node == nil {
		return nil
	}

	if node.Segments != nil {
		var result = mapset.NewSet()
		for _, seg := range node.Segments {
			if seg.Rect.Contains(p) {
				for _, d := range seg.Data.ToSlice() {
					result.Add(d)
				}
			}
		}
		return result.ToSlice()
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

func (node *TreeNode) String() string {
	if node.Segments != nil {
		return fmt.Sprintf("data:%v", node.Segments)
	} else {
		return fmt.Sprintf("{ %v node{axis:%d, mid:%v, min:%v, max:%v}  %v}", node.Left,
			node.Axis, node.Mid, node.Min, node.Max,
			node.Right)
	}
}

func (node *TreeNode) Dump(prefix string) string {
	if node.Segments != nil {
		return fmt.Sprintf("%s -data:%v", prefix, node.Segments)
	} else {
		return fmt.Sprintf("%s -node{axis:%d, mid:%v, min:%v, max:%v}\n%v\n%v\n", prefix,
			node.Axis, node.Mid, node.Min, node.Max,
			node.Left.Dump(prefix+"    "), node.Right.Dump(prefix+"    "))
	}
}
