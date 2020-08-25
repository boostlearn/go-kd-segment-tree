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

	dimNode map[interface{}]ConjunctionDimNode
}

func (node *ConjunctionNode) Search(p Point) []interface{} {
	segCounter := make(map[int]int)
	for dimName, d := range p {
		if node.dimNode[dimName] == nil {
			continue
		}

		for _, segIndex := range node.dimNode[dimName].Search(d) {
			segCounter[segIndex] += 1
		}
	}

	var result = mapset.NewSet()
	for segIndex, matchNum := range segCounter {
		if len(node.segments[segIndex].Rect) == matchNum {
			result = result.Union(node.segments[segIndex].Data)
		}
	}

	return result.ToSlice()
}

func (node *ConjunctionNode) SearchRect(r Rect) []interface{} {
	segCounter :=  make(map[int]int)
	for dimName, d := range r {
		if node.dimNode[dimName] == nil {
			continue
		}

		for _, seg := range node.dimNode[dimName].SearchRect(d) {
			segCounter[seg] += 1
		}
	}

	var result = mapset.NewSet()
	for segIndex, matchNum := range segCounter {
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

func (node *ConjunctionNode) MaxInvertNodeNum() int {
	totalInvertNode := 0
	for _, dimNode := range node.dimNode {
		if dimNode == nil {
			continue
		}

		totalInvertNode += dimNode.MaxInvertNode()
	}
	return totalInvertNode
}

func (node *ConjunctionNode) Dumps(prefix string) string {
	return fmt.Sprintf("%v   conjunction_node{max_invert_node=%v}\n", prefix, node.MaxInvertNodeNum())
}

type ConjunctionDimNode interface {
	Search(measure Measure) []int
	SearchRect(rect interface{}) []int
	MaxInvertNode() int
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
			dimNode.splitPoints[mid+1].Bigger(measure) {
			return dimNode.segments[fmt.Sprintf("%v_%v",
				dimNode.splitPoints[mid], dimNode.splitPoints[mid+1])]
		} else if dimNode.splitPoints[mid+1].SmallerOrEqual(measure) {
			start = mid + 1
		} else if dimNode.splitPoints[mid].Bigger(measure){
			end = mid
		} else {
			break
		}
	}

	return nil
}

func (dimNode *ConjunctionDimRealNode) MaxInvertNode() int {
	if dimNode == nil || len(dimNode.segments) == 0 {
		return 0
	}

	maxNodeNum := 0
	for _, nodes := range dimNode.segments {
		if len(nodes)  > maxNodeNum {
			maxNodeNum = len(nodes)
		}
	}
	return maxNodeNum
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
	var allSplit = []Measure{}
	for _, seg := range segments {
		if seg.Rect[dimName] == nil {
			continue
		}

		foundStart := false
		for _, m := range allSplit {
			if m.Equal(seg.Rect[dimName].(Interval)[0]) {
				foundStart = true
				break
			}
		}
		if foundStart == false {
			allSplit = append(allSplit, seg.Rect[dimName].(Interval)[0])
		}

		foundEnd := false
		for _, m := range allSplit {
			if m.Equal(seg.Rect[dimName].(Interval)[1]) {
				foundEnd = true
				break
			}
		}
		if foundEnd == false {
			allSplit = append(allSplit, seg.Rect[dimName].(Interval)[1])
		}
	}

	var dimNode = &ConjunctionDimRealNode{
		dimName:     dimName,
		splitPoints: allSplit,
		segments:    make(map[string][]int),
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

func (node *ConjunctionDimDiscreteNode) MaxInvertNode() int {
	if node == nil || len(node.segments) == 0 {
		return 0
	}

	maxNodeNum := 0
	for _, nodes := range node.segments {
		if len(nodes)  > maxNodeNum {
			maxNodeNum = len(nodes)
		}
	}
	return maxNodeNum
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
	for segIndex, seg := range segments {
		if seg.Rect[dimName] == nil {
			continue
		}

		for _, m := range seg.Rect[dimName].(Scatters) {
			node.segments[m] = append(node.segments[m], segIndex)
		}
	}

	if len(node.segments) == 0 {
		return nil
	}

	return node
}
