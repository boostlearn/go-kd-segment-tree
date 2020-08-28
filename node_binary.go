package go_kd_segment_tree

import (
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set"
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

	Pass TreeNode
}

func (node *BinaryNode) Search(p Point) []interface{} {
	if node == nil {
		return nil
	}

	if _, ok := p[node.DimName]; ok == false {
		return nil
	}

	x := p[node.DimName]

	var passResult []interface{}
	if node.Pass != nil {
		passResult = node.Pass.Search(p)
	}

	var childResult []interface{}
	if x.Smaller(node.Mid) {
		if node.Left == nil {
			return nil
		}
		childResult = node.Left.Search(p)
	} else {
		if node.Right == nil {
			return nil
		}
		childResult = node.Right.Search(p)
	}

	if len(passResult) == 0 {
		return childResult
	} else if len(childResult) == 0 {
		return passResult
	} else {
		return mapset.NewSet(passResult...).Union(mapset.NewSet(childResult...)).ToSlice()
	}
}

func (node *BinaryNode) SearchRect(r Rect) []interface{} {
	if node == nil {
		return nil
	}

	if _, ok := r[node.DimName]; ok == false {
		return nil
	}

	if _, ok := r[node.DimName].(Interval); ok == false {
		return nil
	}

	dimInterval := r[node.DimName].(Interval)

	var passResult []interface{}
	if node.Pass != nil {
		passResult = node.Pass.SearchRect(r)
	}

	var childResult []interface{}
	if dimInterval[1].Smaller(node.Mid) {
		if node.Left == nil {
			return nil
		}
		childResult = node.Left.SearchRect(r)
	} else if dimInterval[0].Bigger(node.Mid) {
		if node.Right == nil {
			return nil
		}
		childResult = node.Right.SearchRect(r)
	} else {
		childResult = append(childResult, node.Left.SearchRect(r)...)
		childResult = append(childResult, node.Right.SearchRect(r)...)
	}

	if len(passResult) == 0 {
		return childResult
	} else if len(childResult) == 0 {
		return passResult
	} else {
		return mapset.NewSet(passResult...).Union(mapset.NewSet(childResult...)).ToSlice()
	}
}

func (node *BinaryNode) Insert(seg *Segment) error {
	if node == nil || seg == nil {
		return errors.New("binary node is None")
	}

	if _, ok := seg.Rect[node.DimName]; ok == false {
		if node.Pass != nil {
			return node.Pass.Insert(seg)
		} else {
			node.Pass = &LeafNode{
				Segments: []*Segment{seg},
			}
			return nil
		}
	}

	if _, ok := seg.Rect[node.DimName].(Interval); ok == false {
		return errors.New(fmt.Sprintf("wrong binary range: %v", node.DimName))
	}

	dimInterval := seg.Rect[node.DimName].(Interval)
	if dimInterval[1].Smaller(node.Mid) {
		if node.Left == nil {
			node.Left = &LeafNode{
				Segments: []*Segment{seg},
			}
			return nil
		}
		return node.Left.Insert(seg)
	} else if dimInterval[0].Bigger(node.Mid) {
		if node.Right == nil {
			node.Right = &LeafNode{
				Segments: []*Segment{seg},
			}
			return nil
		}
		return node.Right.Insert(seg)
	} else {
		if node.Left == nil {
			node.Left = &LeafNode{
				Segments: []*Segment{seg},
			}
			return nil
		}
		err := node.Left.Insert(seg)
		if err != nil {
			return err
		}
		if node.Right == nil {
			node.Right = &LeafNode{
				Segments: []*Segment{seg},
			}
			return nil
		}
		return node.Right.Insert(seg)
	}
}

func (node *BinaryNode) Dumps(prefix string) string {
	if node == nil {
		return ""
	}

	return fmt.Sprintf("%s -bnode{dim:%d, decreasePercent:%v, mid:%v}\n%v\n%v\n", prefix,
		node.DimName, node.DecreasePercent, node.Mid,
		node.Left.Dumps(prefix+"    "), node.Right.Dumps(prefix+"    "))
}

func NewBinaryNode(tree *Tree,
	segments []*Segment,
	dimName interface{},
	decreasePercent float64,
	level int,
) (*BinaryNode, []*Segment, []*Segment, []*Segment) {
	sort.Stable(&sortSegments{dimName: dimName, segments: segments})

	_, midMeasure := getRealDimSegmentsDecrease(segments, dimName)
	if midMeasure == nil {
		return nil, nil, nil, nil
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
	var pass []*Segment
	for _, seg := range segments {
		if seg.Rect[dimName] == nil {
			pass = append(pass, seg)
			continue
		}

		mRange := seg.Rect[dimName].(Interval)

		if mRange[1].Smaller(midMeasure) {
			left = append(left, seg)
			continue
		} else if mRange[0].Bigger(midMeasure) {
			right = append(right, seg)
			continue
		} else {
			pass = append(pass, seg)
		}
	}

	return node, pass, left, right
}
