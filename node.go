package go_kd_segment_tree

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"sort"
)

type TreeNode struct {
	Axis int

	Level int

	Mid float64
	Min float64
	Max float64

	Left  *TreeNode
	Right *TreeNode

	Segments []*Segment
}

func NewNode(segments []*Segment, axis int, level int, levelMax int, leafDataSizeMin int) *TreeNode {
	if len(segments) == 0 {
		return nil
	}

	if len(segments) < leafDataSizeMin || level > levelMax {
		return &TreeNode{
			Segments:  segments,
		}
	}

	axisSegments := NewSegments(axis, segments)
	sort.Sort(axisSegments)
	mid := len(axisSegments.segments) / 2
	for mid > 0 {
		if axisSegments.segments[mid-1].Rect[axis][0] == axisSegments.segments[mid].Rect[axis][0] {
			mid -= 1
			continue
		}
		break
	}

	nextDim := (axis + 1) % len(segments[0].Rect)
	return &TreeNode{
		Axis:  axis,
		Mid:   axisSegments.segments[mid].Rect[axis][0],
		Min:   axisSegments.segments[0].Rect[axis][0],
		Max:   axisSegments.segments[len(axisSegments.segments)-1].Rect[axis][1],
		Left:  NewNode(axisSegments.segments[:mid], nextDim, level+1, levelMax, leafDataSizeMin),
		Right: NewNode(axisSegments.segments[mid:], nextDim, level+1, levelMax, leafDataSizeMin),
	}
}

func (node *TreeNode) Search(p Point, axis int) []interface{} {
	if node == nil {
		return nil
	}

	if node.Segments != nil {
		var result  = mapset.NewSet()
		for _, seg := range node.Segments {
			if seg.Rect.Contains(p) {
				for _, d := range seg.Data.ToSlice() {
					result.Add(d)
				}
			}
		}
		return result.ToSlice()
	}

	x := p[axis]
	if x < node.Min || x > node.Max {
		return nil
	}

	if x < node.Mid {
		if node.Left == nil {
			return nil
		}
		return node.Left.Search(p, (axis+1)%len(p))
	} else {
		if node.Right == nil {
			return nil
		}
		return node.Right.Search(p, (axis+1)%len(p))
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
