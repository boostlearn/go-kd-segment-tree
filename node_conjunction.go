package go_kd_segment_tree

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"sort"
)

type ConjunctionNode struct {
	TreeNode

	Tree *Tree

	DimName         interface{}
	Level           int
	DecreasePercent float64

	dimNode map[interface{}]ConjunctionDimNode
}

func (node *ConjunctionNode) Search(p Point) []interface{} {
	matchSegments := make(map[*Segment]int)
	for dimName, d := range p {
		for _, seg := range node.dimNode[dimName].Search(d) {
			matchSegments[seg] = matchSegments[seg] + 1
		}
	}

	var result = mapset.NewSet()
	for seg, matchNum := range matchSegments {
		if len(seg.Rect) == matchNum {
			result = result.Union(seg.Data)
		}
	}

	return result.ToSlice()
}

func NewLogicNode(tree *Tree,
	segments []*Segment,
	dimName interface{},
	decreasePercent float64,
	level int,
) *ConjunctionNode {
	var node = &ConjunctionNode{
		Tree:           tree,
		DimName:         dimName,
		Level:           level,
		DecreasePercent: decreasePercent,
		dimNode:         make(map[interface{}]ConjunctionDimNode),
	}

	for dimName, dimType := range tree.dimTypes {
		switch dimType.Type {
		case DimTypeDiscrete.Type:
			node.dimNode[dimName] = NewDiscreteConjunctionNode(segments, dimName)
		case DimTypeReal.Type:
			node.dimNode[dimName] = NewConjunctionRealNode(segments, dimName)
		}
	}

	return node
}

func (node *ConjunctionNode) Dumps(prefix string) string {
	return "conjunction_node"
}

type ConjunctionDimNode interface {
	Search(measure Measure) []*Segment
}

type ConjunctionDimRealNode struct {
	ConjunctionDimNode
	dimName interface{}

	splitPoints []Measure

	segments map[string][]*Segment
}

func (dimNode *ConjunctionDimRealNode) Search(measure Measure) []*Segment {
	if dimNode == nil || len(dimNode.splitPoints) == 0 {
		return nil
	}

	start := 0
	end := len(dimNode.splitPoints) - 1
	for start < end {
		mid := (start + end)/2
		if dimNode.splitPoints[mid].SmallerOrEqual(measure) &&
			dimNode.splitPoints[mid+1].BiggerOrEqual(measure) {
			return dimNode.segments[fmt.Sprintf("%v_%v",
				dimNode.splitPoints[mid], dimNode.splitPoints[mid+1])]
		} else if dimNode.splitPoints[mid].Bigger(measure) {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}

	return nil
}

func NewConjunctionRealNode(segments []*Segment, dimName interface{}) *ConjunctionDimRealNode {
	var allSplit = mapset.NewSet()
	for _, seg := range segments {
		if seg.Rect[dimName] == nil {
			continue
		}

		allSplit.Add(seg.Rect[dimName].(Interval)[0])
		allSplit.Add(seg.Rect[dimName].(Interval)[1])
	}

	var dimNode = &ConjunctionDimRealNode{
		dimName:      dimName,
		splitPoints:       nil,
		segments:     make(map[string][]*Segment),
	}

	for _, t := range allSplit.ToSlice() {
		dimNode.splitPoints = append(dimNode.splitPoints, t.(Measure))
	}

	sort.Sort(&sortMeasures{measures:dimNode.splitPoints})

	for _, seg := range segments {
		for i, m := range dimNode.splitPoints {
			nextM := dimNode.splitPoints[i+1]
			if seg.Rect[dimName].(Interval)[0].SmallerOrEqual(m) &&
				seg.Rect[dimName].(Interval)[1].BiggerOrEqual(nextM) {
				key := fmt.Sprintf("%v_%v", m, nextM)
				dimNode.segments[key] = append(dimNode.segments[key], seg)
			}
		}
	}

	return dimNode
}

type ConjunctionDimDiscreteNode struct {
	ConjunctionNode

	dimName interface{}

	segments map[Measure][]*Segment
}

func (node *ConjunctionDimDiscreteNode) Search(measure Measure) []*Segment {
	return node.segments[measure]
}

func NewDiscreteConjunctionNode(segments []*Segment, dimName interface{}) *ConjunctionDimDiscreteNode {
	node := &ConjunctionDimDiscreteNode{
		dimName:   dimName,
		segments:     make(map[Measure][]*Segment),
	}
	for _, seg := range segments {
		for _, m := range seg.Rect[dimName].(Scatters) {
			node.segments[m] = append(node.segments[m], seg)
		}
	}

	return node
}