package go_kd_segment_tree

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

var testRects []Rect
var searchPoint []Point
var rectNum int = 10000
var dimNum int = 3
var targetRate float64 = 1.0

func init() {
	for i := 0; i < rectNum; i++ {
		rect := Rect{}
		point := Point{}
		for j := 0; j < dimNum; j++ {
			k := rand.Float64()
			if rand.Float64() < targetRate {
				rect = append(rect, [2]Measure{FloatMeasure(k), FloatMeasure(k + 0.0001)})
			} else {
				rect = append(rect, [2]Measure{MeasureMin{}, MeasureMax{}})
			}
			point = append(point, FloatMeasure(k))
		}
		testRects = append(testRects, rect)
		searchPoint = append(searchPoint, point)
	}
}

/*
func TestNewTree(t *testing.T) {
	testSize := 10
	total := 0
	for i, p := range searchPoint {
		for j, rect := range testRects[:testSize] {
			if rect.Contains(p) {
				fmt.Printf("notree: %v %v %v %v\n", i, j, rect, p)
				total += 1
			}
		}
	}
	fmt.Println("notree: ", total)

	tree := NewTree(&TreeOptions{
		TreeLevelMax:     0,
		LeafNodeMin:      0,
		BranchingGiniMin: 1.0,
	})
	for i, rect := range testRects[:testSize] {
		tree.Add(rect, "data"+strconv.FormatInt(int64(i), 10))
	}
	tree.Build()
	fmt.Println(tree)
	total = 0
	for i, p := range searchPoint {
		d := tree.Search(p)
		if len(d) > 0 {
			fmt.Printf("tree: %v %v, %v\n", i, p, d)
		}
		total += len(d)
	}
	fmt.Println("tree: ", total)
}
*/

func BenchmarkTree_Search(b *testing.B) {
	tree := NewTree(&TreeOptions{
		TreeLevelMax:     16,
		LeafNodeMin:      16,
		BranchingGiniMin: 0.2,
	})
	for i, rect := range testRects {
		tree.Add(rect, "data"+strconv.FormatInt(int64(i), 10))
	}
	tree.Build()

	fmt.Println("tree:", tree.String())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p := searchPoint[i%len(searchPoint)]
		_ = tree.Search(p)
	}

}

func BenchmarkNoTree_Search(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p := searchPoint[i%len(searchPoint)]
		for _, rect := range testRects {
			if rect.Contains(p) {
			}
		}
	}

}
