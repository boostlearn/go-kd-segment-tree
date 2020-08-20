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

type SegmentBranching struct {
	axis     int
	segments []*Segment

	mid int
	min Measure
	max Measure

	GiniCoefficient float64

	midSeg *Segment
	left   []*Segment
	right  []*Segment
}

func (b *SegmentBranching) Len() int {
	return len(b.segments)
}

func (b *SegmentBranching) Less(i, j int) bool {
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

func (b *SegmentBranching) Swap(i, j int) {
	b.segments[i], b.segments[j] = b.segments[j], b.segments[i]
}

func NewSegmentBranch(segments []*Segment, minJini float64) *SegmentBranching {
	if len(segments) == 0 {
		return nil
	}

	segmentBranch := &SegmentBranching{
		segments: segments,
	}

	maxGiniAxis := 0
	maxGiniMid := 0
	maxGiniCoefficient := -1.0
	for axis, _ := range segments[0].Rect {
		segmentBranch.axis = axis
		sort.Sort(segmentBranch)

		start := 0
		end := len(segmentBranch.segments) - 1
		for start <= end && segmentBranch.segments[start].Rect[axis][0].Equal(MeasureMin{}) {
			start += 1
		}
		for start <= end && segmentBranch.segments[start].Rect[axis][0].Equal(MeasureMax{}) {
			end -= 1
		}
		if start == end {
			continue
		}

		mid := (start + end)/2
		for mid > 0 {
			if segmentBranch.segments[mid-1].Rect[axis][0].Equal(segmentBranch.segments[mid].Rect[axis][0]) {
				mid -= 1
				continue
			}
			break
		}
		midSeg := segmentBranch.segments[mid]

		leftNum := 0
		rightNum := 0
		for _, seg := range segments {
			if seg.Rect[axis][1].Smaller(midSeg.Rect[axis][0]) {
				leftNum += 1
			} else if seg.Rect[axis][0].BiggerOrEqual(midSeg.Rect[axis][0]) {
				rightNum += 1
			} else {
				leftNum += 1
				rightNum += 1
			}
		}

		p := 0.0
		if leftNum < rightNum {
			p = float64(rightNum) * 1.0 / float64(len(segments))
		} else {
			p = float64(leftNum) * 1.0 / float64(len(segments))
		}
		axisJini := 1 - p*p - (1-p)*(1-p)

		if axisJini > maxGiniCoefficient {
			maxGiniCoefficient = axisJini
			maxGiniAxis = axis
			maxGiniMid = mid
		}
	}

	if maxGiniCoefficient < minJini {
		return nil
	}

	segmentBranch.axis = maxGiniAxis
	segmentBranch.GiniCoefficient = maxGiniCoefficient

	sort.Sort(segmentBranch)
	segmentBranch.mid = maxGiniMid
	segmentBranch.midSeg = segmentBranch.segments[segmentBranch.mid]
	segmentBranch.min = segmentBranch.midSeg.Rect[maxGiniAxis][0]
	segmentBranch.max = segmentBranch.midSeg.Rect[maxGiniAxis][1]

	for _, seg := range segments {
		if seg.Rect[maxGiniAxis][1].Smaller(segmentBranch.midSeg.Rect[maxGiniAxis][0]) {
			segmentBranch.left = append(segmentBranch.left, seg)
		} else if seg.Rect[maxGiniAxis][0].BiggerOrEqual(segmentBranch.midSeg.Rect[maxGiniAxis][0]) {
			segmentBranch.right = append(segmentBranch.right, seg)
		} else {
			leftSlice := seg.SliceClone(maxGiniAxis, seg.Rect[maxGiniAxis][0], segmentBranch.midSeg.Rect[maxGiniAxis][0])
			segmentBranch.left = append(segmentBranch.left, leftSlice)

			rightSlice := seg.SliceClone(maxGiniAxis, segmentBranch.midSeg.Rect[maxGiniAxis][0], seg.Rect[maxGiniAxis][1])
			segmentBranch.right = append(segmentBranch.right, rightSlice)
		}

		if seg.Rect[maxGiniAxis][0].Smaller(segmentBranch.min) {
			segmentBranch.min = seg.Rect[maxGiniAxis][0]
		}

		if seg.Rect[maxGiniAxis][1].Bigger(segmentBranch.max) {
			segmentBranch.max = seg.Rect[maxGiniAxis][1]
		}
	}

	return segmentBranch
}

func MergeSegments(segments []*Segment) []*Segment {
	var newSegments []*Segment
	var uniqMap = make(map[string]*Segment)
	for _, seg := range segments {
		rectKey := seg.Rect.Key()
		if s, ok := uniqMap[rectKey]; ok {
			s.Data = s.Data.Union(seg.Data)
		}else {
			newSegments = append(newSegments, seg)
		}
	}
	return newSegments
}