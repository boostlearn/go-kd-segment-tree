package go_kd_segment_tree

import (
	"fmt"
	"testing"
)

func TestInterval_Contains(t *testing.T) {
	a := MeasureFloat(0.1)
	var r = 1
	b := Interval{MeasureFloat(0.2), MeasureFloat(r)}

	var c interface{}
	c = a
	switch c.(type) {
	case Measure:
		fmt.Println("Measure")
	case Interval:
		fmt.Println("Interval")
	}

	c = b
	switch c.(type) {
	case Measure:
		fmt.Println("Measure")
	case Interval:
		fmt.Println("Interval")
	}
}
