package go_kd_segment_tree

import (
	"fmt"
	"strings"
)

type HashNode struct {
	TreeNode

	Tree            *Tree
	DimName         interface{}
	Level           int
	DecreasePercent float64

	child map[Measure]TreeNode
}

func (node *HashNode) Search(p Point) []interface{} {
	if node == nil {
		return nil
	}

	if _, ok := p[node.DimName]; ok == false {
		return nil
	}

	x := p[node.DimName]
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
	msgs = append(msgs, fmt.Sprintf("%s -hnode{dim:%d, decreasePercent:%v}", prefix, node.DimName, node.DecreasePercent))
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
) (*HashNode, map[Measure][]*Segment) {
	hashSegments := make(map[Measure][]*Segment)

	for _, seg := range segments {
		if seg.Rect[dimName] == nil {
			continue
		}
		for _, key := range seg.Rect[dimName].(Scatters) {
			hashSegments[key] = append(hashSegments[key], seg)
		}
	}

	for _, seg := range segments {
		if seg.Rect[dimName] == nil {
			for key := range hashSegments {
				hashSegments[key] = append(hashSegments[key], seg)
			}
		}
	}

	node := &HashNode{
		Tree:            tree,
		DimName:         dimName,
		Level:           level,
		DecreasePercent: decreasePercent,
		child:           make(map[Measure]TreeNode),
	}

	return node, hashSegments
}
