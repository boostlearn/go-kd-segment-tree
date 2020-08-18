package go_kd_segment_tree

import (
	"fmt"
	"math/rand"
	"sort"
	mapset "github.com/deckarep/golang-set"
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

	var axisCutPoints []float64
	var cutPointsMap = make(map[float64]bool)
	for _, s := range segments {
		start := s.Rect[axis][0]
		if cutPointsMap[start] == false {
			axisCutPoints = append(axisCutPoints, start)
			cutPointsMap[start] = true
		}
		end := s.Rect[axis][1]
		if cutPointsMap[end] == false {
			axisCutPoints = append(axisCutPoints, end)
			cutPointsMap[end] = true
		}
	}
	if len(axisCutPoints) == 0 {
		return nil
	}
	sort.Float64s(axisCutPoints)

	var allSegment []*Segment
	var segmentMap = make(map[string]*Segment)
	for i, _ := range axisCutPoints[:len(axisCutPoints)-1] {
		start := axisCutPoints[i]
		end := axisCutPoints[i+1]
		for _, segment := range segments {
			if segment.Rect[axis][0] <= start && segment.Rect[axis][1] >= end {
				segmentSlice := segment.SliceClone(axis, start, end)
				rectKey := segmentSlice.Rect.Key()
				if dup, ok := segmentMap[rectKey]; ok {
					for _, t := range segmentSlice.Data.ToSlice() {
						dup.Data.Add(t)
					}
					continue
				}
				allSegment = append(allSegment, segmentSlice)
			}
		}
	}

	newAxisSegments := &AxisSegments{
		axis:     axis,
		segments: allSegment,
	}

	return newAxisSegments
}
