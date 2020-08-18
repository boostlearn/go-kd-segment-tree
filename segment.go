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
		newRect = append(newRect, [2]float64{r[0], r[1]})
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

func (s *Segment) Range(axis int) (float64, float64) {
	return s.Rect[axis][0], s.Rect[axis][1]
}

func (r Rect) Contains(p Point) bool {
	if len(r) != len(p) {
		return false
	}

	for i, r := range r {
		if r[0] > p[i] || r[1] < p[i] {
			return false
		}
	}

	return true
}

func (s *Segment) String() string {
	return fmt.Sprintf("{%v, %v}", s.Rect, s.Data)
}

func (s *Segment) SliceClone(axis int, start float64, end float64) *Segment {
	newSegment := &Segment{
		Rect: s.Rect.Clone(),
		Data: s.Data.Clone(),
	}

	newSegment.Rect[axis][0] = start
	newSegment.Rect[axis][1] = end
	return newSegment
}

type AxisSegments struct {
	axis     int
	segments []*Segment
	mid      int
	midSeg   *Segment
	left     []*Segment
	right    []*Segment

	min float64
	max float64
}

func (b *AxisSegments) Len() int {
	return len(b.segments)
}

func (b *AxisSegments) Less(i, j int) bool {
	if b.segments[i].Rect[b.axis][0] != b.segments[j].Rect[b.axis][0] {
		return b.segments[i].Rect[b.axis][0] < b.segments[j].Rect[b.axis][0]
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

func (b *AxisSegments) Swap(i, j int) {
	b.segments[i], b.segments[j] = b.segments[j], b.segments[i]
}

func NewSegments(axis int, segments []*Segment) *AxisSegments {
	if len(segments) == 0 {
		return nil
	}

	newAxisSegments := &AxisSegments{
		axis:     axis,
		segments: segments,
	}

	sort.Sort(newAxisSegments)
	newAxisSegments.mid = len(newAxisSegments.segments) / 2
	for newAxisSegments.mid > 0 {
		if newAxisSegments.segments[newAxisSegments.mid-1].Rect[axis][0] ==
			newAxisSegments.segments[newAxisSegments.mid].Rect[axis][0] {
			newAxisSegments.mid -= 1
			continue
		}
		break
	}
	newAxisSegments.midSeg = newAxisSegments.segments[newAxisSegments.mid]
	newAxisSegments.min = newAxisSegments.midSeg.Rect[axis][0]
	newAxisSegments.max = newAxisSegments.midSeg.Rect[axis][1]

	for _, seg := range segments {
		if seg.Rect[axis][1] < newAxisSegments.midSeg.Rect[axis][0] {
			newAxisSegments.left = append(newAxisSegments.left, seg)
		} else if seg.Rect[axis][0] >= newAxisSegments.midSeg.Rect[axis][0] {
			newAxisSegments.right = append(newAxisSegments.right, seg)
		} else {
			leftSlice := seg.SliceClone(axis, seg.Rect[axis][0], newAxisSegments.midSeg.Rect[axis][0])
			newAxisSegments.left = append(newAxisSegments.left, leftSlice)

			rightSlice := seg.SliceClone(axis, newAxisSegments.midSeg.Rect[axis][0], seg.Rect[axis][1])
			newAxisSegments.right = append(newAxisSegments.right, rightSlice)
		}

		if seg.Rect[axis][0] < newAxisSegments.min {
			newAxisSegments.min = seg.Rect[axis][0]
		}

		if seg.Rect[axis][1] > newAxisSegments.max {
			newAxisSegments.max = seg.Rect[axis][1]
		}
	}

	return newAxisSegments
}
