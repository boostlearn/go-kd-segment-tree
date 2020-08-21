package go_kd_segment_tree

import (
	"fmt"
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
	case MeasureMax:
		return false
	case MeasureMin:
		return true
	}
	return false
}

func (a FloatMeasure) Smaller(b interface{}) bool {
	switch b.(type) {
	case FloatMeasure:
		return float64(a) < float64(b.(FloatMeasure))
	case MeasureMax:
		return true
	case MeasureMin:
		return false
	}
	return false

}

func (a FloatMeasure) Equal(b interface{}) bool {
	switch b.(type) {
	case FloatMeasure:
		return float64(a) == float64(b.(FloatMeasure))
	case MeasureMax:
		return false
	case MeasureMin:
		return false
	}
	return false
}

func (a FloatMeasure) BiggerOrEqual(b interface{}) bool {
	switch b.(type) {
	case FloatMeasure:
		return float64(a) >= float64(b.(FloatMeasure))
	case MeasureMax:
		return false
	case MeasureMin:
		return true
	}
	return false

}

func (a FloatMeasure) SmallerOrEqual(b interface{}) bool {
	switch b.(type) {
	case FloatMeasure:
		return float64(a) <= float64(b.(FloatMeasure))
	case MeasureMax:
		return true
	case MeasureMin:
		return false
	}
	return false
}

type StringMeasure string

func (a StringMeasure) Bigger(b interface{}) bool {
	switch b.(type) {
	case StringMeasure:
		return string(a) > string(b.(StringMeasure))
	case MeasureMax:
		return false
	case MeasureMin:
		return true
	}
	return false
}

func (a StringMeasure) Smaller(b interface{}) bool {
	switch b.(type) {
	case StringMeasure:
		return string(a) < string(b.(StringMeasure))
	case MeasureMax:
		return true
	case MeasureMin:
		return false
	}
	return false
}

func (a StringMeasure) Equal(b interface{}) bool {
	switch b.(type) {
	case StringMeasure:
		return string(a) == string(b.(StringMeasure))
	case MeasureMax:
		return false
	case MeasureMin:
		return false
	}
	return false
}

func (a StringMeasure) BiggerOrEqual(b interface{}) bool {
	switch b.(type) {
	case StringMeasure:
		return string(a) >= string(b.(StringMeasure))
	case MeasureMax:
		return false
	case MeasureMin:
		return true
	}
	return false
}

func (a StringMeasure) SmallerOrEqual(b interface{}) bool {
	switch b.(type) {
	case StringMeasure:
		return string(a) <= string(b.(StringMeasure))
	case MeasureMax:
		return true
	case MeasureMin:
		return false
	}
	return false
}

type TimeMeasure time.Time

func (f TimeMeasure) Bigger(b interface{}) bool {
	switch b.(type) {
	case TimeMeasure:
		return time.Time(f).After(b.(time.Time))
	case MeasureMax:
		return false
	case MeasureMin:
		return true
	}
	return false
}

func (f TimeMeasure) Smaller(b interface{}) bool {
	switch b.(type) {
	case TimeMeasure:
		return time.Time(f).Before(b.(time.Time))
	case MeasureMax:
		return true
	case MeasureMin:
		return false
	}
	return false

}

func (f TimeMeasure) Equal(b interface{}) bool {
	switch b.(type) {
	case TimeMeasure:
		return time.Time(f).Equal(b.(time.Time))
	case MeasureMax:
		return true
	case MeasureMin:
		return false
	}
	return false
}

func (f TimeMeasure) BiggerOrEqual(b interface{}) bool {
	switch b.(type) {
	case TimeMeasure:
		return time.Time(f).After(b.(time.Time)) || time.Time(f).Equal(b.(time.Time))
	case MeasureMax:
		return false
	case MeasureMin:
		return true
	}
	return false
}

func (f TimeMeasure) SmallerOrEqual(b interface{}) bool {
	switch b.(type) {
	case TimeMeasure:
		return time.Time(f).Before(b.(time.Time)) || time.Time(f).Equal(b.(time.Time))
	case MeasureMax:
		return true
	case MeasureMin:
		return false
	}
	return false
}

type MeasureMin struct{}

func (f MeasureMin) Bigger(b interface{}) bool {
	return false
}

func (f MeasureMin) Smaller(b interface{}) bool {
	switch b.(type) {
	case MeasureMin:
		return false
	default:
		return true
	}
}

func (f MeasureMin) Equal(b interface{}) bool {
	switch b.(type) {
	case MeasureMin:
		return true
	default:
		return false
	}
}

func (f MeasureMin) BiggerOrEqual(b interface{}) bool {
	switch b.(type) {
	case MeasureMin:
		return true
	default:
		return false
	}
}

func (f MeasureMin) SmallerOrEqual(b interface{}) bool {
	return true
}

func (f MeasureMin) String() string {
	return fmt.Sprintf("<-INFINITE>")
}

type MeasureMax struct{}

func (f MeasureMax) Bigger(b interface{}) bool {
	switch b.(type) {
	case MeasureMax:
		return false
	default:
		return true
	}
}

func (f MeasureMax) Smaller(b interface{}) bool {
	return false
}

func (f MeasureMax) Equal(b interface{}) bool {
	switch b.(type) {
	case MeasureMax:
		return true
	default:
		return false
	}
}

func (f MeasureMax) BiggerOrEqual(b interface{}) bool {
	return true
}

func (f MeasureMax) SmallerOrEqual(b interface{}) bool {
	switch b.(type) {
	case MeasureMax:
		return true
	default:
		return true
	}
}

func (f MeasureMax) String() string {
	return fmt.Sprintf("<+INFINITE>")
}

type Interval [2]Measure

func (i Interval) Contains(p Measure) bool {
	return p.BiggerOrEqual(i[0]) && p.SmallerOrEqual(i[1])
}
