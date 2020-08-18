package go_kd_segment_tree

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"math/rand"
	"sort"
	"strings"
)

func (r Rect) Clone() Rect {
	var newRect Rect
	for _, r := range r {
		newRect = append(newRect, [2]Measure{r[0], r[1]})
	}
	return newRect
}

func (r Rect) Key() string {
	var dimKeys []string
	for _, r := range r {
		dimKeys = append(dimKeys, fmt.Sprintf("%v_%v", r[0], r[1]))
	}
	return strings.Join(dimKeys, ":")
}

type Segment struct {
	Rect Rect
	Data mapset.Set

	rnd float64
}

func (s *Segment) Range(axis int) (Measure, Measure) {
	return s.Rect[axis][0], s.Rect[axis][1]
}

func (r Rect) Contains(p Point) bool {
	if len(r) != len(p) {
		return false
	}

	for i, r := range r {
		if r[0].Bigger(p[i]) || r[1].Smaller(p[i]) {
			return false
		}
	}

	return true
}

func (s *Segment) String() string {
	return fmt.Sprintf("{%v, %v}", s.Rect, s.Data)
}

func (s *Segment) SliceClone(axis int, start Measure, end Measure) *Segment {
	newSegment := &Segment{
		Rect: s.Rect.Clone(),
		Data: s.Data.Clone(),
	}

	newSegment.Rect[axis][0] = start
	newSegment.Rect[axis][1] = end
	return newSegment
}

type SplitSegments struct {
	axis     int
	segments []*Segment
	mid      int
	midSeg   *Segment
	left     []*Segment
	right    []*Segment

	min Measure
	max Measure
}

func (b *SplitSegments) Len() int {
	return len(b.segments)
}

func (b *SplitSegments) Less(i, j int) bool {
	if b.segments[i].Rect[b.axis][0] != b.segments[j].Rect[b.axis][0] {
		return b.segments[i].Rect[b.axis][0].Smaller(b.segments[j].Rect[b.axis][0])
	} else {
		if b.segments[i].rnd == 0 {
			b.segments[i].rnd = rand.Float64()
		}

		if b.segments[j].rnd == 0 {
			b.segments[j].rnd = rand.Float64()
		}

		return b.segments[i].rnd < b.segments[j].rnd
	}
}

func (b *SplitSegments) Swap(i, j int) {
	b.segments[i], b.segments[j] = b.segments[j], b.segments[i]
}

func NewSplitSegments(segments []*Segment) *SplitSegments {
	if len(segments) == 0 {
		return nil
	}

	splitAxis := 0
	maxKeySize := 0
	for axis, _ := range segments[0].Rect {
		keyMap := mapset.NewSet()
		for _, seg := range segments {
			keyMap.Add(seg.Rect[axis][0])
		}
		keys := keyMap.ToSlice()
		if len(keys) > maxKeySize {
			splitAxis = axis
			maxKeySize = len(keys)
		}
	}

	newAxisSegments := &SplitSegments{
		axis:     splitAxis,
		segments: segments,
	}

	sort.Sort(newAxisSegments)
	newAxisSegments.mid = len(newAxisSegments.segments) / 2
	for newAxisSegments.mid > 0 {
		if newAxisSegments.segments[newAxisSegments.mid-1].Rect[splitAxis][0] ==
			newAxisSegments.segments[newAxisSegments.mid].Rect[splitAxis][0] {
			newAxisSegments.mid -= 1
			continue
		}
		break
	}
	newAxisSegments.midSeg = newAxisSegments.segments[newAxisSegments.mid]
	newAxisSegments.min = newAxisSegments.midSeg.Rect[splitAxis][0]
	newAxisSegments.max = newAxisSegments.midSeg.Rect[splitAxis][1]

	for _, seg := range segments {
		if seg.Rect[splitAxis][1].Smaller(newAxisSegments.midSeg.Rect[splitAxis][0]) {
			newAxisSegments.left = append(newAxisSegments.left, seg)
		} else if seg.Rect[splitAxis][0].BiggerOrEqual(newAxisSegments.midSeg.Rect[splitAxis][0]) {
			newAxisSegments.right = append(newAxisSegments.right, seg)
		} else {
			leftSlice := seg.SliceClone(splitAxis, seg.Rect[splitAxis][0], newAxisSegments.midSeg.Rect[splitAxis][0])
			newAxisSegments.left = append(newAxisSegments.left, leftSlice)

			rightSlice := seg.SliceClone(splitAxis, newAxisSegments.midSeg.Rect[splitAxis][0], seg.Rect[splitAxis][1])
			newAxisSegments.right = append(newAxisSegments.right, rightSlice)
		}

		if seg.Rect[splitAxis][0].Smaller(newAxisSegments.min) {
			newAxisSegments.min = seg.Rect[splitAxis][0]
		}

		if seg.Rect[splitAxis][1].Bigger(newAxisSegments.max) {
			newAxisSegments.max = seg.Rect[splitAxis][1]
		}
	}

	return newAxisSegments
}
