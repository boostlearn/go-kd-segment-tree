package go_kd_segment_tree

import (
	"fmt"
	"sort"
)

type BinaryNode struct {
	TreeNode

	Tree *Tree

	AxisName interface{}
	Level    int
	DecreasePercent     float64

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

	if _, ok := p[node.AxisName]; ok == false {
		return nil
	}

	x := p[node.AxisName]
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

	return fmt.Sprintf("%s -bnode{axis:%d, decreasePercent:%v, mid:%v, min:%v, max:%v}\n%v\n%v\n", prefix,
		node.AxisName, node.DecreasePercent, node.Mid, node.Min, node.Max,
		node.Left.Dumps(prefix+"    "), node.Right.Dumps(prefix+"    "))
}

func NewBinaryNode(tree *Tree,
	segments []*Segment,
	axisName interface{},
	decreasePercent float64,
	level int,
) (*BinaryNode, []*Segment, []*Segment) {
	sort.Stable(&sortSegments{axisName:axisName, segments:segments})

	_, midMeasure := RealDimSegmentsDecrease(segments, axisName)
	if midMeasure == nil {
		return nil, nil, nil
	}

	node := &BinaryNode{
		Tree:tree,
		AxisName: axisName,
		Level:    level,
		DecreasePercent:     decreasePercent,
		Mid:      midMeasure,
	}

	var left []*Segment
	var right []*Segment
	for _, seg := range segments {
		if seg.Rect[axisName] == nil {
			left = append(left, seg)
			right = append(right, seg)
			continue
		}

		mRange := seg.Rect[axisName].(Interval)
		if mRange[0].Smaller(node.Min) {
			node.Min = mRange[0]
		}
		if mRange[1].Bigger(node.Max) {
			node.Max = mRange[1]
		}

		if mRange[1].Smaller(midMeasure) {
			left = append(left, seg)
			continue
		}

		if mRange[0].Bigger(midMeasure) {
			right = append(right, seg)
			continue
		}

		left = append(left, seg)
		right = append(right, seg)

	}

	return node, left, right
}