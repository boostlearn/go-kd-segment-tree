package go_kd_segment_tree

import (
	"fmt"
	"strings"
	"time"
)

type Measure interface {
	Bigger(b interface{}) bool
	Smaller(b interface{}) bool
	Equal(b interface{}) bool
	BiggerOrEqual(b interface{}) bool
	SmallerOrEqual(b interface{}) bool
}

type MeasureFloat float64

func (a MeasureFloat) Bigger(b interface{}) bool {
	switch b.(type) {
	case MeasureFloat:
		return float64(a) > float64(b.(MeasureFloat))
	}
	return false
}

func (a MeasureFloat) Smaller(b interface{}) bool {
	switch b.(type) {
	case MeasureFloat:
		return float64(a) < float64(b.(MeasureFloat))
	}
	return false

}

func (a MeasureFloat) Equal(b interface{}) bool {
	switch b.(type) {
	case MeasureFloat:
		return float64(a) == float64(b.(MeasureFloat))
	}
	return false
}

func (a MeasureFloat) BiggerOrEqual(b interface{}) bool {
	switch b.(type) {
	case MeasureFloat:
		return float64(a) >= float64(b.(MeasureFloat))
	}
	return false

}

func (a MeasureFloat) SmallerOrEqual(b interface{}) bool {
	switch b.(type) {
	case MeasureFloat:
		return float64(a) <= float64(b.(MeasureFloat))
	}
	return false
}

type MeasureString string

func (a MeasureString) Bigger(b interface{}) bool {
	switch b.(type) {
	case MeasureString:
		return string(a) > string(b.(MeasureString))
	}
	return false
}

func (a MeasureString) Smaller(b interface{}) bool {
	switch b.(type) {
	case MeasureString:
		return string(a) < string(b.(MeasureString))
	}
	return false
}

func (a MeasureString) Equal(b interface{}) bool {
	switch b.(type) {
	case MeasureString:
		return string(a) == string(b.(MeasureString))
	}
	return false
}

func (a MeasureString) BiggerOrEqual(b interface{}) bool {
	switch b.(type) {
	case MeasureString:
		return string(a) >= string(b.(MeasureString))
	}
	return false
}

func (a MeasureString) SmallerOrEqual(b interface{}) bool {
	switch b.(type) {
	case MeasureString:
		return string(a) <= string(b.(MeasureString))
	}
	return false
}

type MeasureTime time.Time

func (f MeasureTime) Bigger(b interface{}) bool {
	switch b.(type) {
	case MeasureTime:
		return time.Time(f).After(b.(time.Time))
	}
	return false
}

func (f MeasureTime) Smaller(b interface{}) bool {
	switch b.(type) {
	case MeasureTime:
		return time.Time(f).Before(b.(time.Time))
	}
	return false

}

func (f MeasureTime) Equal(b interface{}) bool {
	switch b.(type) {
	case MeasureTime:
		return time.Time(f).Equal(b.(time.Time))
	}
	return false
}

func (f MeasureTime) BiggerOrEqual(b interface{}) bool {
	switch b.(type) {
	case MeasureTime:
		return time.Time(f).After(b.(time.Time)) || time.Time(f).Equal(b.(time.Time))
	}
	return false
}

func (f MeasureTime) SmallerOrEqual(b interface{}) bool {
	switch b.(type) {
	case MeasureTime:
		return time.Time(f).Before(b.(time.Time)) || time.Time(f).Equal(b.(time.Time))
	}
	return false
}

type Interval [2]Measure
type Intervals []Interval

func (i Interval) Contains(p Measure) bool {
	if p == nil {
		return true
	}

	return p.BiggerOrEqual(i[0]) && p.SmallerOrEqual(i[1])
}

type Measures []Measure

func (s Measures) Contains(p Measure) bool {
	if p == nil {
		return true
	}

	for _, m := range s {
		if m.Equal(p) {
			return true
		}
	}
	return false
}

type Point map[interface{}]Measure

type Rect map[interface{}]interface{}

func (rect Rect) Clone() Rect {
	var newRect = make(Rect)
	for name, d := range rect {
		switch d.(type) {
		case Interval:
			newRect[name] = Interval{d.(Interval)[0], d.(Interval)[1]}
		case Measures:
			var newSc Measures
			for _, s := range d.(Measures) {
				newSc = append(newSc, s)
			}
			newRect[name] = newSc
		}

	}
	return newRect
}

func (rect Rect) Key() string {
	var dimKeys []string
	for name, d := range rect {
		switch d.(type) {
		case Interval:
			dimKeys = append(dimKeys, fmt.Sprintf("%v=%v_%v",
				name, d.(Interval)[0], d.(Interval)[1]))
		case Measures:
			dimKeys = append(dimKeys, fmt.Sprintf("%v_%v",
				name, d.(Measures)))
		}

	}
	return strings.Join(dimKeys, ":")
}

func (rect Rect) Contains(p Point) bool {
	for name, d := range rect {
		switch d.(type) {
		case Interval:
			if d.(Interval).Contains(p[name]) == false {
				return false
			}
		case Intervals:
			found := false
			for _, interval := range d.(Intervals) {
				if interval.Contains(p[name]) {
					found = true
					break
				}
			}
			if found == false {
				return false
			}
		case Measure:
			if d.(Measure).Equal(p[name]) == false {
				return false
			}
		case Measures:
			found := false
			for _, dD := range d.(Measures) {
				if dD.Equal(p[name]) {
					found = true
					break
				}
			}
			if found == false {
				return false
			}
		}

	}
	return true
}

func (rect Rect) HasIntersect(s Rect) bool {
	for name, d := range rect {
		p := s[name]
		if p == nil {
			continue
		}

		switch d.(type) {
		case Interval:
			switch p.(type) {
			case Interval:
				if d.(Interval).Contains(p.(Interval)[0]) == false && d.(Interval).Contains(p.(Interval)[1]) == false {
					return false
				}
			case Intervals:
				found := false
				for _, pInterval := range p.(Intervals) {
					if d.(Interval).Contains(pInterval[0])  || d.(Interval).Contains(pInterval[1])  {
						found  = true
						break
					}
				}
				if found == false {
					return false
				}
			default:
				return false
			}
		case Intervals:
			switch p.(type) {
			case Interval:
				found := false
				for _, dInterval := range d.(Intervals) {
					if dInterval.Contains(p.(Interval)[0])  || d.(Interval).Contains(p.(Interval)[1])  {
						found  = true
						break
					}
				}
				if found == false {
					return false
				}
			case Intervals:
				found := false
				for _, dInterval := range d.(Intervals) {
					for _, pInterval := range p.(Intervals) {
						if dInterval.Contains(pInterval[0]) || dInterval.Contains(pInterval[1]) {
							found = true
							break
						}
					}
					if found == true {
						break
					}
				}
				if found == false {
					return false
				}
			default:
				return false
			}
		case Measure:
			switch p.(type) {
			case Measure:
				if d.(Measure).Equal(p.(Measure)) == false {
					return false
				}
			case Measures:
				found := false
				for _, pP := range p.(Measures) {
					if d.(Measure).Equal(pP) {
						found = true
						break
					}
				}
				if found == false {
					return false
				}
			default:
				return false
			}
		case Measures:
			switch p.(type) {
			case Measure:
				found := false
				for _, dD := range d.(Measures) {
					if dD.Equal(p.(Measure)) == false {
						found = true
						break
					}
				}
				if found == false {
					return false
				}
			case Measures:
				found := false
				for _, dD := range d.(Measures) {
					for _, pP := range p.(Measures) {
						if dD.Equal(pP) {
							found = true
							break
						}
					}
					if found {
						break
					}
				}
				if found == false {
					return false
				}
			default:
				return false
			}
		}

	}
	return true
}
