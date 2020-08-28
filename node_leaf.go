package go_kd_segment_tree

import (
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set"
)

type LeafNode struct {
	TreeNode
	Segments []*Segment
}

func (node *LeafNode) Search(p Point) []interface{} {
	if node == nil {
		return nil
	}
	if node.Segments != nil {
		var result = mapset.NewSet()
		for _, seg := range node.Segments {
			if seg.Rect.Contains(p) {
				result = result.Union(seg.Data)
			}
		}
		return result.ToSlice()
	}
	return nil
}

func (node *LeafNode) SearchRect(r Rect) []interface{} {
	if node == nil {
		return nil
	}
	if node.Segments != nil {
		var result = mapset.NewSet()
		for _, seg := range node.Segments {
			if seg.Rect.HasIntersect(r) {
				result = result.Union(seg.Data)
			}
		}
		return result.ToSlice()
	}
	return nil
}

func (node *LeafNode) Insert(seg *Segment) error {
	if node == nil {
		return errors.New("leaf node is nil")
	}
	node.Segments = append(node.Segments, seg)

	return nil
}

func (node *LeafNode) Dumps(prefix string) string {
	if node == nil {
		return ""
	}

	return fmt.Sprintf("%s -leaf:{size=%v}", prefix, len(node.Segments))
}

func MergeSegments(segments []*Segment) []*Segment {
	var newSegments []*Segment
	var uniqMap = make(map[string]*Segment)
	for _, seg := range segments {
		rectKey := seg.Rect.Key()
		if s, ok := uniqMap[rectKey]; ok {
			s.Data = s.Data.Union(seg.Data)
		} else {
			newSegments = append(newSegments, seg)
		}
	}
	return newSegments
}
