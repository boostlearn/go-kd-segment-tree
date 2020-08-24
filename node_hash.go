package go_kd_segment_tree

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"strings"
)

type HashNode struct {
	TreeNode

	Tree            *Tree
	DimName         interface{}
	Level           int
	DecreasePercent float64

	child map[Measure]TreeNode
	pass TreeNode
}

func (node *HashNode) Search(p Point) []interface{} {
	if node == nil {
		return nil
	}

	if _, ok := p[node.DimName]; ok == false {
		return nil
	}

	x := p[node.DimName]

	var defaultResult []interface{}
	if node.pass != nil {
		 defaultResult = node.pass.Search(p)
	}

	var childResult []interface{}
	if child, ok := node.child[x]; ok {
		childResult =  child.Search(p)
	}

	if len(defaultResult) == 0 {
		return childResult
	} else if len(childResult) == 0 {
		return defaultResult
	} else {
		return mapset.NewSet(defaultResult...).Union(mapset.NewSet(childResult...)).ToSlice()
	}
}

func (node *HashNode) Dumps(prefix string) string {
	if node == nil {
		return ""
	}

	var msgs []string
	msgs = append(msgs, fmt.Sprintf("%s -hnode{dim:%d, decreasePercent:%v}", prefix, node.DimName, node.DecreasePercent))
	msgs = append(msgs, node.pass.Dumps(fmt.Sprintf("%v    %v:", prefix, "<PASS>")))
	for childKey, child := range node.child {
		msgs = append(msgs, child.Dumps(fmt.Sprintf("%v    %v:", prefix, childKey)))
	}
	return strings.Join(msgs, "\n")
}

func NewHashNode(tree *Tree,
	segments []*Segment,
	dimName interface{},
	decreasePercent float64,
	level int,
) (*HashNode, []*Segment, map[Measure][]*Segment) {
	hashSegments := make(map[Measure][]*Segment)

	var passSegments []*Segment
	for _, seg := range segments {
		if seg.Rect[dimName] == nil {
			passSegments = append(passSegments, seg)
			continue
		}
		for _, key := range seg.Rect[dimName].(Scatters) {
			hashSegments[key] = append(hashSegments[key], seg)
		}
	}

	node := &HashNode{
		Tree:            tree,
		DimName:         dimName,
		Level:           level,
		DecreasePercent: decreasePercent,
		child:           make(map[Measure]TreeNode),
	}

	return node, passSegments, hashSegments
}
