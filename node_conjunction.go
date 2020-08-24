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

	segments []*Segment

	counter []int

	dimNode map[interface{}]ConjunctionDimNode
}

func (node *ConjunctionNode) Search(p Point) []interface{} {
	for i := 0; i < len(node.counter); i++ {
		node.counter[i] = 0
	}

	for dimName, d := range p {
		if node.dimNode[dimName] == nil {
			continue
		}

		for _, segIndex := range node.dimNode[dimName].Search(d) {
			node.counter[segIndex] += 1
		}
	}

	var result = mapset.NewSet()
	for segIndex, matchNum := range node.counter {
		if len(node.segments[segIndex].Rect) == matchNum {
			result = result.Union(node.segments[segIndex].Data)
		}
	}

	return result.ToSlice()
}

func (node *ConjunctionNode) SearchRect(r Rect) []interface{} {
	for i := 0; i < len(node.counter); i++ {
		node.counter[i] = 0
	}
	for dimName, d := range r {
		if node.dimNode[dimName] == nil {
			continue
		}

		for _, seg := range node.dimNode[dimName].SearchRect(d) {
			node.counter[seg] += 1
		}
	}

	var result = mapset.NewSet()
	for segIndex, matchNum := range node.counter {
		if len(node.segments[segIndex].Rect) == matchNum {
			result = result.Union(node.segments[segIndex].Data)
		}
	}

	return result.ToSlice()
}

func NewConjunctionNode(tree *Tree,
	segments []*Segment,
	dimName interface{},
	decreasePercent float64,
	level int,
) *ConjunctionNode {
	var node = &ConjunctionNode{
		Tree:            tree,
		DimName:         dimName,
		Level:           level,
		segments:segments,
		counter:make([]int, len(segments)),
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
	Search(measure Measure) []int
	SearchRect(rect interface{}) []int
}

type ConjunctionDimRealNode struct {
	ConjunctionDimNode
	dimName interface{}

	splitPoints []Measure

	segments map[string][]int
}

func (dimNode *ConjunctionDimRealNode) Search(measure Measure) []int {
	if dimNode == nil || len(dimNode.splitPoints) == 0 {
		return nil
	}

	start := 0
	end := len(dimNode.splitPoints) - 1
	for start < end {
		mid := (start + end) / 2
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

func (dimNode *ConjunctionDimRealNode) SearchRect(measure interface{}) []int {
	if dimNode == nil || len(dimNode.splitPoints) == 0 {
		return nil
	}

	if _, ok := measure.(Interval); ok == false {
		return nil
	}

	interval := measure.(Interval)

	matchSegments := mapset.NewSet()
	for i, _ := range dimNode.splitPoints[:len(dimNode.splitPoints)-1] {
		if interval.Contains(dimNode.splitPoints[i]) &&
			interval.Contains(dimNode.splitPoints[i+1]) {
			for _, seg := range dimNode.segments[fmt.Sprintf("%v_%v",
				dimNode.splitPoints[i], dimNode.splitPoints[i+1])] {
				matchSegments.Add(seg)
			}
		}
	}

	var result []int
	for _, seg := range matchSegments.ToSlice() {
		result = append(result, seg.(int))
	}
	return result
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
		dimName:     dimName,
		splitPoints: nil,
		segments:    make(map[string][]int),
	}

	for _, t := range allSplit.ToSlice() {
		dimNode.splitPoints = append(dimNode.splitPoints, t.(Measure))
	}

	if len(dimNode.splitPoints) == 0 {
		return nil
	}

	sort.Sort(&sortMeasures{measures: dimNode.splitPoints})

	for index, seg := range segments {
		for i, m := range dimNode.splitPoints[:len(dimNode.splitPoints)-1] {
			nextM := dimNode.splitPoints[i+1]
			if seg.Rect[dimName] == nil {
				continue
			}

			if seg.Rect[dimName].(Interval)[0].SmallerOrEqual(m) &&
				seg.Rect[dimName].(Interval)[1].BiggerOrEqual(nextM) {
				key := fmt.Sprintf("%v_%v", m, nextM)
				dimNode.segments[key] = append(dimNode.segments[key], index)
			}
		}
	}

	return dimNode
}

type ConjunctionDimDiscreteNode struct {
	ConjunctionNode

	dimName interface{}

	segments map[Measure][]int
}

func (node *ConjunctionDimDiscreteNode) Search(measure Measure) []int {
	if node == nil || node.segments == nil {
		return nil
	}

	return node.segments[measure]
}

func (node *ConjunctionDimDiscreteNode) SearchRect(scatters interface{}) []int {
	if node == nil || node.segments == nil {
		return nil
	}
	if _, ok := scatters.(Scatters); ok == false {
		return nil
	}

	matchSegments := mapset.NewSet()
	for _, d := range scatters.(Scatters) {
		for _, seg := range node.segments[d] {
			matchSegments.Add(seg)
		}
	}

	var result []int
	for _, seg := range matchSegments.ToSlice() {
		result = append(result, seg.(int))
	}
	return result
}

func NewDiscreteConjunctionNode(segments []*Segment, dimName interface{}) *ConjunctionDimDiscreteNode {
	node := &ConjunctionDimDiscreteNode{
		dimName:  dimName,
		segments: make(map[Measure][]int),
	}
	for _, seg := range segments {
		if seg.Rect[dimName] == nil {
			continue
		}

		for index, m := range seg.Rect[dimName].(Scatters) {
			node.segments[m] = append(node.segments[m], index)
		}
	}

	if len(node.segments) == 0 {
		return nil
	}

	return node
}
