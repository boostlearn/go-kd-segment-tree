package go_kd_segment_tree

import (
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"sync"
)

const DefaultTreeLevelMax = 16
const DefaultLeafDataMin = 16
const DefaultBranchDecreasePercentMin = 0.4

type DimType struct{ Type int }

var DimTypeDiscrete = DimType{Type: 0}
var DimTypeReal = DimType{Type: 1}

type Tree struct {
	mu       sync.RWMutex
	updateMu sync.Mutex

	dimTypes map[interface{}]DimType

	options *TreeOptions

	segments []*Segment
	root     TreeNode
}

type TreeOptions struct {
	TreeLevelMax                int
	LeafNodeMin                 int
	BranchingDecreasePercentMin float64
}

func NewTree(dimTypes map[interface{}]DimType, opts *TreeOptions) *Tree {
	if opts == nil {
		opts = &TreeOptions{}
	}

	if opts.TreeLevelMax == 0 {
		opts.TreeLevelMax = DefaultTreeLevelMax
	}

	if opts.LeafNodeMin == 0 {
		opts.LeafNodeMin = DefaultLeafDataMin
	}

	if opts.BranchingDecreasePercentMin == 0.0 {
		opts.BranchingDecreasePercentMin = DefaultBranchDecreasePercentMin
	}
	return &Tree{
		dimTypes: dimTypes,
		options:  opts,
	}
}

func (tree *Tree) Search(p Point) []interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	if tree.root == nil {
		return nil
	}
	return tree.root.Search(p)
}

func (tree *Tree) SearchRect(r Rect) []interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	if tree.root == nil {
		return nil
	}
	return tree.root.SearchRect(r)
}

func (tree *Tree) Dumps() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	return fmt.Sprintf("%v", tree.root.Dumps(""))
}

func (tree *Tree) Add(rect Rect, data interface{}) error {
	tree.updateMu.Lock()
	defer tree.updateMu.Unlock()

	for name, d := range rect {
		switch d.(type) {
		case Interval:
			if tree.dimTypes[name] != DimTypeReal {
				return errors.New(fmt.Sprintf("dim type error:%v", name))
			}
		case Scatters:
			if tree.dimTypes[name] != DimTypeDiscrete {
				return errors.New(fmt.Sprintf("dim type error:%v", name))
			}
		}
	}

	seg := &Segment{
		Rect: rect.Clone(),
		Data: mapset.NewSet(data),
	}

	tree.segments = append(tree.segments, seg)
	return nil
}

func (tree *Tree) Remove(data interface{}) {
	tree.updateMu.Lock()
	defer tree.updateMu.Unlock()

	var newSegments []*Segment
	for _, seg := range tree.segments {
		if seg.Data.Contains(data) {
			seg.Data.Remove(data)
			if len(seg.Data.ToSlice()) > 0 {
				newSegments = append(newSegments, seg)
			}
		} else {
			newSegments = append(newSegments, seg)
		}
	}

	tree.segments = newSegments
}

func (tree *Tree) Build() {
	tree.updateMu.Lock()
	defer tree.updateMu.Unlock()

	if len(tree.segments) == 0 {
		return
	}

	newNode := NewNode(tree.segments, tree, 1)

	tree.mu.Lock()
	tree.root = newNode
	tree.mu.Unlock()
}
