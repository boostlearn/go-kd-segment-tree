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
var dimNum int = 2

func init() {
	for i := 0 ; i < rectNum; i++ {
		rect := Rect{}
		point := Point{}
		for j := 0; j < dimNum; j++ {
			k := rand.Float64()
			rect = append(rect, [2]float64{k, k+0.0001})
			point = append(point, k)
		}
		testRects = append(testRects, rect)
		searchPoint = append(searchPoint, point)
	}
}

func TestNewTree(t *testing.T) {
	testSize := 100
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


	tree := NewTree(12, 16)
	for i, rect := range testRects[:testSize] {
		tree.Add(rect, "data" + strconv.FormatInt(int64(i), 10))
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


func BenchmarkTree_Search(b *testing.B) {
	tree := NewTree(12, 16)
	for i, rect := range testRects {
		tree.Add(rect, "data" + strconv.FormatInt(int64(i), 10))
	}
	tree.Build()

	b.ReportAllocs()
	b.ResetTimer()


	for i := 0; i < b.N; i++ {
		total := 0
		for _, p := range searchPoint {
			d := tree.Search(p)
			total += len(d)
		}
	}

}

func BenchmarkNoTree_Search(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		total := 0
		for _, p := range searchPoint {
			for _, rect := range testRects {
				if rect.Contains(p) {
					total += 1
				}
			}
		}
	}

}