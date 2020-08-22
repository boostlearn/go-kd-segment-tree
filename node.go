package go_kd_segment_tree

type TreeNode interface {
	Search(p Point) []interface{}
	Dumps(prefix string) string
}

func NewNode(segments []*Segment,
	tree *Tree,
	level int,
) TreeNode {
	if len(segments) == 0 {
		return nil
	}

	if len(segments) < tree.options.LeafNodeMin || level > tree.options.TreeLevelMax {
		return &LeafNode{
			Segments: MergeSegments(segments),
		}
	}

	axisName, decreasePercent := GetBestBranchingAxisName(segments, tree.dimTypes)
	if decreasePercent < tree.options.BranchingDecreasePercentMin {
		return &LeafNode{
			Segments: MergeSegments(segments),
		}
	}

	switch tree.dimTypes[axisName].Type {
	case DimTypeReal.Type:
		node, left, right := NewBinaryNode(tree, segments, axisName, decreasePercent, level)
		if len(left) > 0 {
			node.Left = NewNode(left, tree, level+1)
		}
		if len(right) > 0 {
			node.Right = NewNode(right, tree, level+1)
		}
		return node
	case DimTypeDiscrete.Type:
		node, children := NewHashNode(tree, segments, axisName, decreasePercent, level)
		for childKey, childSegments := range children {
			node.child[childKey] = NewNode(childSegments, tree, level+1)
		}
		return node
	}
	return nil
}



func GetBestBranchingAxisName (
	segments []*Segment,
	dimTypes map[interface{}]DimType,
) (interface{}, float64) {
	if len(segments) == 0 {
		return nil, 0
	}

	var maxDecreaseAxisName interface{}
	var maxDecrease int
	for axisName, dimType := range dimTypes {
		if dimType.Type == DimTypeReal.Type {
			decreaseC,_ := RealDimSegmentsDecrease(segments, axisName)
			if decreaseC > maxDecrease {
				maxDecrease = decreaseC
				maxDecreaseAxisName = axisName
			}
		} else if dimType.Type == DimTypeDiscrete.Type {
			decreaseC,_ := DiscreteDimSegmentsDecrease(segments, axisName)
			if decreaseC > maxDecrease {
				maxDecrease = decreaseC
				maxDecreaseAxisName = axisName
			}
		}
	}

	p := float64(maxDecrease) * 1.0 / float64(len(segments))
	return maxDecreaseAxisName, p
}