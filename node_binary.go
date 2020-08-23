package go_kd_segment_tree

import (
	"fmt"
	"sort"
)

type BinaryNode struct {
	TreeNode

	Tree *Tree

	DimName         interface{}
	Level           int
	DecreasePercent float64

	Mid Measure

	Left  TreeNode
	Right TreeNode
}

func (node *BinaryNode) Search(p Point) []interface{} {
	if node == nil {
		return nil
	}

	if _, ok := p[node.DimName]; ok == false {
		return nil
	}

	x := p[node.DimName]

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

	return fmt.Sprintf("%s -bnode{dim:%d, decreasePercent:%v,mid:%v}\n%v\n%v\n", prefix,
		node.DimName, node.DecreasePercent, node.Mid,
		node.Left.Dumps(prefix+"    "), node.Right.Dumps(prefix+"    "))
}

func NewBinaryNode(tree *Tree,
	segments []*Segment,
	dimName interface{},
	decreasePercent float64,
	level int,
) (*BinaryNode, []*Segment, []*Segment) {
	sort.Stable(&sortSegments{dimName: dimName, segments: segments})

	_, midMeasure := getRealDimSegmentsDecrease(segments, dimName)
	if midMeasure == nil {
		return nil, nil, nil
	}

	node := &BinaryNode{
		Tree:            tree,
		DimName:         dimName,
		Level:           level,
		DecreasePercent: decreasePercent,
		Mid:             midMeasure,
	}

	var left []*Segment
	var right []*Segment
	for _, seg := range segments {
		if seg.Rect[dimName] == nil {
			left = append(left, seg)
			right = append(right, seg)
			continue
		}

		mRange := seg.Rect[dimName].(Interval)

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
