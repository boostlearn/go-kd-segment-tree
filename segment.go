package go_kd_segment_tree

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"math/rand"
	"sort"
	"strings"
)

func (rect Rect) Clone() Rect {
	var newRect Rect
	for _, d := range rect {
		switch d.(type) {
		case Interval:
			newRect = append(newRect, Interval{d.(Interval)[0], d.(Interval)[1]})
		case Scatters:
			var newSc Scatters
			for _, s := range d.(Scatters) {
				newSc = append(newSc, s)
			}
			newRect = append(newRect, newSc)
		}

	}
	return newRect
}

func (rect Rect) Key() string {
	var dimKeys []string
	for _, d := range rect {
		switch d.(type) {
		case Interval:
			dimKeys = append(dimKeys, fmt.Sprintf("%v_%v", d.(Interval)[0], d.(Interval)[1]))
		case Scatters:
			dimKeys = append(dimKeys, fmt.Sprintf("%v", d.(Scatters)))
		}

	}
	return strings.Join(dimKeys, ":")
}

type Segment struct {
	Rect Rect
	Data mapset.Set
	rnd float64
}

func (rect Rect) Contains(p Point) bool {
	if len(rect) != len(p) {
		return false
	}

	for axis, d := range rect {
		switch d.(type) {
		case Interval:
			if d.(Interval).Contains(p[axis]) == false {
				return false
			}
		case Scatters:
			if d.(Scatters).Contains(p[axis]) == false {
				return false
			}
		}

	}

	return true
}

func (s *Segment) String() string {
	return fmt.Sprintf("{%v, %v}", s.Rect, s.Data)
}

func (s *Segment) Clone() *Segment {
	newSegment := &Segment{
		Rect: s.Rect.Clone(),
		Data: s.Data.Clone(),
	}
	return newSegment
}

type SegmentBranching struct {
	axis     int
	gini float64
	segments []*Segment

	mid int
	min Measure
	max Measure
	midSeg *Segment
	left   []*Segment
	right  []*Segment

	hashSegments map[Measure][]*Segment
}

func (branching *SegmentBranching) Len() int {
	return len(branching.segments)
}

func (branching *SegmentBranching) Less(i, j int) bool {
	switch branching.segments[i].Rect[branching.axis].(type) {
	case Interval:
		if branching.segments[i].Rect[branching.axis].(Interval)[0].Equal(branching.segments[j].Rect[branching.axis].(Interval)[0]) == false {
			return branching.segments[i].Rect[branching.axis].(Interval)[0].Smaller(branching.segments[j].Rect[branching.axis].(Interval)[0])
		}
	}

	if branching.segments[i].rnd == 0 {
		branching.segments[i].rnd = rand.Float64()
	}

	if branching.segments[j].rnd == 0 {
		branching.segments[j].rnd = rand.Float64()
	}

	return branching.segments[i].rnd < branching.segments[j].rnd

}

func (branching *SegmentBranching) Swap(i, j int) {
	branching.segments[i], branching.segments[j] = branching.segments[j], branching.segments[i]
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
		switch segments[0].Rect[axis].(type) {
		case Interval:
			segmentBranch.axis = axis
			sort.Sort(segmentBranch)

			start := 0
			end := len(segmentBranch.segments) - 1
			for start <= end && segmentBranch.segments[start].Rect[axis].(Interval)[0].Equal(MeasureMin{}) {
				start += 1
			}
			for start <= end && segmentBranch.segments[start].Rect[axis].(Interval)[0].Equal(MeasureMax{}) {
				end -= 1
			}
			if start == end {
				continue
			}

			mid := (start + end) / 2
			for mid > 0 {
				if segmentBranch.segments[mid-1].Rect[axis].(Interval)[0].Equal(segmentBranch.segments[mid].Rect[axis].(Interval)[0]) {
					mid -= 1
					continue
				}
				break
			}
			midSeg := segmentBranch.segments[mid]

			leftNum := 0
			rightNum := 0
			for _, seg := range segments {
				if seg.Rect[axis].(Interval)[1].Smaller(midSeg.Rect[axis].(Interval)[0]) {
					leftNum += 1
				} else if seg.Rect[axis].(Interval)[0].BiggerOrEqual(midSeg.Rect[axis].(Interval)[0]) {
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
		case Scatters:
			var scatterMap = make(map[Measure]int)
			for _, seg := range segmentBranch.segments {
				for _, s := range seg.Rect[axis].(Scatters) {
					scatterMap[s] = scatterMap[s] + 1
				}
			}
			var maxCounter = 0
			for _, n := range scatterMap {
				if n > maxCounter {
					maxCounter = n
				}
			}

			p := float64(maxCounter) * 1.0 / float64(len(segments))
			axisJini := 1 - p*p - (1-p)*(1-p)

			if axisJini > maxGiniCoefficient {
				maxGiniCoefficient = axisJini
				maxGiniAxis = axis
			}
		}
	}

	if maxGiniCoefficient < minJini {
		return nil
	}

	segmentBranch.axis = maxGiniAxis
	segmentBranch.gini = maxGiniCoefficient

	switch segments[0].Rect[segmentBranch.axis].(type) {
	case Interval:
		sort.Sort(segmentBranch)
		segmentBranch.mid = maxGiniMid
		segmentBranch.midSeg = segmentBranch.segments[segmentBranch.mid]
		segmentBranch.min = segmentBranch.midSeg.Rect[maxGiniAxis].(Interval)[0]
		segmentBranch.max = segmentBranch.midSeg.Rect[maxGiniAxis].(Interval)[1]

		for _, seg := range segments {
			if seg.Rect[maxGiniAxis].(Interval)[1].Smaller(segmentBranch.midSeg.Rect[maxGiniAxis].(Interval)[0]) {
				segmentBranch.left = append(segmentBranch.left, seg)
			} else if seg.Rect[maxGiniAxis].(Interval)[0].BiggerOrEqual(segmentBranch.midSeg.Rect[maxGiniAxis].(Interval)[0]) {
				segmentBranch.right = append(segmentBranch.right, seg)
			} else {
				segmentBranch.left = append(segmentBranch.left, seg.Clone())
				segmentBranch.right = append(segmentBranch.right, seg)
			}

			if seg.Rect[maxGiniAxis].(Interval)[0].Smaller(segmentBranch.min) {
				segmentBranch.min = seg.Rect[maxGiniAxis].(Interval)[0]
			}

			if seg.Rect[maxGiniAxis].(Interval)[1].Bigger(segmentBranch.max) {
				segmentBranch.max = seg.Rect[maxGiniAxis].(Interval)[1]
			}
		}
	case Scatters:
		segmentBranch.hashSegments = make(map[Measure][]*Segment)
		for _, seg := range segmentBranch.segments {
			for _, key := range seg.Rect[segmentBranch.axis].(Scatters) {
				if _, ok := segmentBranch.hashSegments[key]; ok {
					segmentBranch.hashSegments[key] = append(segmentBranch.hashSegments[key], seg)
				} else {
					segmentBranch.hashSegments[key] = []*Segment{seg}
				}
			}
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
		} else {
			newSegments = append(newSegments, seg)
		}
	}
	return newSegments
}
