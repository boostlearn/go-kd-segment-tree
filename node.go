package go_kd_segment_tree

type TreeNode interface {
	Search(p Point) []interface{}
	SearchRect(rect Rect) []interface{}
	Dumps(prefix string) string
}

func NewNode(segments []*Segment,
	tree *Tree,
	level int,
) TreeNode {
	if len(segments) == 0 {
		return nil
	}

	conjunctionNode := NewConjunctionNode(tree, segments, nil, 1.0, level+1)
	conjunctionNodeDecreasePercent := float64(len(segments) - conjunctionNode.MaxInvertNodeNum()) /float64(len(segments))
	if conjunctionNodeDecreasePercent < 0 {
		conjunctionNodeDecreasePercent = 0
	}

	if len(segments) < tree.options.LeafNodeMin || level >= tree.options.TreeLevelMax {
		mergedSegments := MergeSegments(segments)
		if conjunctionNodeDecreasePercent > 0 {
			return conjunctionNode
		} else {
			return &LeafNode{
				Segments: mergedSegments,
			}
		}
	}

	dimName, decreasePercent := findBestBranchingDim(segments, tree.dimTypes)
	if decreasePercent < tree.options.BranchingDecreasePercentMin || decreasePercent < conjunctionNodeDecreasePercent {
		mergedSegments := MergeSegments(segments)
		if conjunctionNodeDecreasePercent > 0 {
			return conjunctionNode
		} else {
			return &LeafNode{
				Segments: mergedSegments,
			}
		}
	}

	switch tree.dimTypes[dimName].Type {
	case DimTypeReal.Type:
		node, pass, left, right := NewBinaryNode(tree, segments, dimName, decreasePercent, level)
		if len(pass) > 0 {
			node.Pass = NewNode(pass, tree, level+1)
		}
		if len(left) > 0 {
			node.Left = NewNode(left, tree, level+1)
		}
		if len(right) > 0 {
			node.Right = NewNode(right, tree, level+1)
		}
		return node
	case DimTypeDiscrete.Type:
		node, passSegments, children := NewHashNode(tree, segments, dimName, decreasePercent, level)
		for childKey, childSegments := range children {
			node.child[childKey] = NewNode(childSegments, tree, level+1)
		}
		if len(passSegments) > 0 {
			node.pass = NewNode(passSegments, tree, level+1)
		}
		return node
	}
	return nil
}

func findBestBranchingDim(
	segments []*Segment,
	dimTypes map[interface{}]DimType,
) (interface{}, float64) {
	if len(segments) == 0 {
		return nil, 0
	}

	var maxDecreaseDimName interface{}
	var maxDecrease int
	for dimName, dimType := range dimTypes {
		switch dimType.Type {
		case DimTypeReal.Type:
			decreaseC, _ := getRealDimSegmentsDecrease(segments, dimName)
			if decreaseC > maxDecrease {
				maxDecrease = decreaseC
				maxDecreaseDimName = dimName
			}
		case DimTypeDiscrete.Type:
			decreaseC, _ := getDiscreteDimSegmentsDecrease(segments, dimName)
			if decreaseC > maxDecrease {
				maxDecrease = decreaseC
				maxDecreaseDimName = dimName
			}
		}
	}

	p := float64(maxDecrease) * 1.0 / float64(len(segments))
	return maxDecreaseDimName, p
}
