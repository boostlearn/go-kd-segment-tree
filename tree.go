package go_kd_segment_tree

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"sync"
)

type Point []Measure
type Rect [][2]Measure

const DefaultTreeLevelMax = 12
const DefaultLeafDataMin = 16

type Tree struct {
	mu       sync.RWMutex
	updateMu sync.Mutex

	levelMax        int
	leafDataSizeMin int

	segments []*Segment
	root     *TreeNode
}

func NewTree(maxLevel int, minLeafDataSize int) *Tree {
	if maxLevel == 0 {
		maxLevel = DefaultTreeLevelMax
	}

	if minLeafDataSize == 0 {
		minLeafDataSize = DefaultLeafDataMin
	}
	return &Tree{
		levelMax:        maxLevel,
		leafDataSizeMin: minLeafDataSize,
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

	newNode := NewNode(tree.segments, tree.levelMax, tree.leafDataSizeMin)

	tree.mu.Lock()
	tree.root = newNode
	tree.mu.Unlock()
}
