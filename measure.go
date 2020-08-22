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

type FloatMeasure float64

func (a FloatMeasure) Bigger(b interface{}) bool {
	switch b.(type) {
	case FloatMeasure:
		return float64(a) > float64(b.(FloatMeasure))
	}
	return false
}

func (a FloatMeasure) Smaller(b interface{}) bool {
	switch b.(type) {
	case FloatMeasure:
		return float64(a) < float64(b.(FloatMeasure))
	}
	return false

}

func (a FloatMeasure) Equal(b interface{}) bool {
	switch b.(type) {
	case FloatMeasure:
		return float64(a) == float64(b.(FloatMeasure))
	}
	return false
}

func (a FloatMeasure) BiggerOrEqual(b interface{}) bool {
	switch b.(type) {
	case FloatMeasure:
		return float64(a) >= float64(b.(FloatMeasure))
	}
	return false

}

func (a FloatMeasure) SmallerOrEqual(b interface{}) bool {
	switch b.(type) {
	case FloatMeasure:
		return float64(a) <= float64(b.(FloatMeasure))
	}
	return false
}

type StringMeasure string

func (a StringMeasure) Bigger(b interface{}) bool {
	switch b.(type) {
	case StringMeasure:
		return string(a) > string(b.(StringMeasure))
	}
	return false
}

func (a StringMeasure) Smaller(b interface{}) bool {
	switch b.(type) {
	case StringMeasure:
		return string(a) < string(b.(StringMeasure))
	}
	return false
}

func (a StringMeasure) Equal(b interface{}) bool {
	switch b.(type) {
	case StringMeasure:
		return string(a) == string(b.(StringMeasure))
	}
	return false
}

func (a StringMeasure) BiggerOrEqual(b interface{}) bool {
	switch b.(type) {
	case StringMeasure:
		return string(a) >= string(b.(StringMeasure))
	}
	return false
}

func (a StringMeasure) SmallerOrEqual(b interface{}) bool {
	switch b.(type) {
	case StringMeasure:
		return string(a) <= string(b.(StringMeasure))
	}
	return false
}

type TimeMeasure time.Time

func (f TimeMeasure) Bigger(b interface{}) bool {
	switch b.(type) {
	case TimeMeasure:
		return time.Time(f).After(b.(time.Time))
	}
	return false
}

func (f TimeMeasure) Smaller(b interface{}) bool {
	switch b.(type) {
	case TimeMeasure:
		return time.Time(f).Before(b.(time.Time))
	}
	return false

}

func (f TimeMeasure) Equal(b interface{}) bool {
	switch b.(type) {
	case TimeMeasure:
		return time.Time(f).Equal(b.(time.Time))
	}
	return false
}

func (f TimeMeasure) BiggerOrEqual(b interface{}) bool {
	switch b.(type) {
	case TimeMeasure:
		return time.Time(f).After(b.(time.Time)) || time.Time(f).Equal(b.(time.Time))
	}
	return false
}

func (f TimeMeasure) SmallerOrEqual(b interface{}) bool {
	switch b.(type) {
	case TimeMeasure:
		return time.Time(f).Before(b.(time.Time)) || time.Time(f).Equal(b.(time.Time))
	}
	return false
}

type Interval [2]Measure

func (i Interval) Contains(p Measure) bool {
	if p == nil {
		return true
	}

	return p.BiggerOrEqual(i[0]) && p.SmallerOrEqual(i[1])
}

type Scatters []Measure

func (s Scatters) Contains(p Measure) bool {
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
		case Scatters:
			var newSc Scatters
			for _, s := range d.(Scatters) {
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
		case Scatters:
			dimKeys = append(dimKeys, fmt.Sprintf("%v_%v",
				name, d.(Scatters)))
		}

	}
	return strings.Join(dimKeys, ":")
}

func (rect Rect) Contains(p Point) bool {
	if len(rect) != len(p) {
		return false
	}

	for name, d := range rect {
		switch d.(type) {
		case Interval:
			if d.(Interval).Contains(p[name]) == false {
				return false
			}
		case Scatters:
			if d.(Scatters).Contains(p[name]) == false {
				return false
			}
		}

	}

	return true
}
