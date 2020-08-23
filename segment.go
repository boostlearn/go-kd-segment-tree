package go_kd_segment_tree

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"math/rand"
	"sort"
)

type Segment struct {
	Rect Rect
	Data mapset.Set
	rnd  float64
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

type sortSegments struct {
	dimName  interface{}
	segments []*Segment
}

func (s *sortSegments) Len() int {
	return len(s.segments)
}

func (s *sortSegments) Less(i, j int) bool {
	iSeg, iSegOk := s.segments[i].Rect[s.dimName]
	jSeg, jSegOk := s.segments[j].Rect[s.dimName]

	if iSegOk == true && jSegOk == false {
		return false
	}

	if iSegOk == false && jSegOk == true {
		return true
	}

	if iSegOk && jSegOk {
		switch iSeg.(type) {
		case Interval:
			if iSeg.(Interval)[0].Equal(jSeg.(Interval)[0]) == false {
				return iSeg.(Interval)[0].Smaller(jSeg.(Interval)[0])
			}
		}
	}

	for s.segments[i].rnd == s.segments[j].rnd {
		s.segments[i].rnd = rand.Float64()
		s.segments[j].rnd = rand.Float64()
	}

	return s.segments[i].rnd < s.segments[j].rnd

}

func (s *sortSegments) Swap(i, j int) {
	s.segments[i], s.segments[j] =
		s.segments[j], s.segments[i]
}

type sortMeasures struct {
	measures []Measure
}

func (s *sortMeasures) Len() int {
	return len(s.measures)
}

func (s *sortMeasures) Less(i, j int) bool {
	return s.measures[i].Smaller(s.measures[j])
}

func (s *sortMeasures) Swap(i, j int) {
	s.measures[i], s.measures[j] =
		s.measures[j], s.measures[i]
}

func getRealDimSegmentsDecrease(segments []*Segment, dimName interface{}) (int, Measure) {
	var dimSegments []*Segment
	for _, seg := range segments {
		if seg.Rect[dimName] != nil {
			dimSegments = append(dimSegments, seg)
		}
	}
	if len(dimSegments) == 0 {
		return 0, nil
	}

	sort.Sort(&sortSegments{dimName: dimName, segments: dimSegments})

	var starts []Measure
	var ends []Measure
	for _, seg := range dimSegments {
		starts = append(starts, seg.Rect[dimName].(Interval)[0])
		ends = append(ends, seg.Rect[dimName].(Interval)[1])
	}
	sort.Sort(&sortMeasures{measures:starts})
	sort.Sort(&sortMeasures{measures:ends})
	pos := 0
	for pos < len(starts) - 1 && ends[pos].Smaller(starts[len(ends) - 1 - pos]) {
		pos += 1
	}
	midMeasure := ends[pos]

	leftNum := 0
	rightNum := 0
	for _, seg := range dimSegments {
		if seg.Rect[dimName].(Interval)[1].Smaller(midMeasure) {
			leftNum += 1
		} else if seg.Rect[dimName].(Interval)[0].BiggerOrEqual(midMeasure) {
			rightNum += 1
		}
	}

	if leftNum < rightNum {
		return leftNum, midMeasure
	} else {
		return rightNum, midMeasure
	}

}

func getDiscreteDimSegmentsDecrease(segments []*Segment, dimName interface{}) (int, Measure) {
	var dimSegments []*Segment
	for _, seg := range segments {
		if seg.Rect[dimName] != nil {
			dimSegments = append(dimSegments, seg)
		}
	}
	if len(dimSegments) == 0 {
		return 0, nil
	}

	var scatterMap = make(map[Measure]int)
	for _, seg := range dimSegments {
		for _, s := range seg.Rect[dimName].(Scatters) {
			scatterMap[s] = scatterMap[s] + 1
		}
	}
	var maxCounter = 0
	var maxMeasure Measure
	for m, n := range scatterMap {
		if n > maxCounter {
			maxCounter = n
			maxMeasure = m
		}
	}

	return len(dimSegments) - maxCounter, maxMeasure
}
