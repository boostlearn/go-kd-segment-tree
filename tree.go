package go_kd_segment_tree

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"sync"
)

type Point []Measure
type Rect [][2]Measure

const DefaultTreeLevelMax = 16
const DefaultLeafDataMin = 16
const DefaultBranchGiniMin = 0.2

type Tree struct {
	mu       sync.RWMutex
	updateMu sync.Mutex

	treeLevelMax     int
	leafNodeMin      int
	branchingGiniMin float64

	segments []*Segment
	root     *TreeNode
}

type TreeOptions struct {
	TreeLevelMax     int
	LeafNodeMin      int
	BranchingGiniMin float64
}

func NewTree(opts *TreeOptions) *Tree {
	if opts == nil {
		opts = &TreeOptions{}
	}

	if opts.TreeLevelMax == 0 {
		opts.TreeLevelMax = DefaultTreeLevelMax
	}

	if opts.LeafNodeMin == 0 {
		opts.LeafNodeMin = DefaultLeafDataMin
	}

	if opts.BranchingGiniMin == 0.0 {
		opts.BranchingGiniMin = DefaultBranchGiniMin
	}
	return &Tree{
		treeLevelMax:     opts.TreeLevelMax,
		leafNodeMin:      opts.LeafNodeMin,
		branchingGiniMin: opts.BranchingGiniMin,
	}
}

func (tree *Tree) Search(p Point) []interface{} {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	return tree.root.Search(p)
}

func (tree *Tree) String() string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	return fmt.Sprintf("%v", tree.root.Dump(""))
}

func (tree *Tree) Add(rect Rect, data interface{}) {
	tree.updateMu.Lock()
	defer tree.updateMu.Unlock()

	seg := &Segment{
		Rect: rect.Clone(),
		Data: mapset.NewSet(data),
	}

	tree.segments = append(tree.segments, seg)
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

	newNode := NewNode(tree.segments, tree.treeLevelMax, tree.leafNodeMin, tree.branchingGiniMin)

	tree.mu.Lock()
	tree.root = newNode
	tree.mu.Unlock()
}
