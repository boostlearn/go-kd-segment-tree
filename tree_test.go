package go_kd_segment_tree

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

var testRects []Rect
var searchPoint []Point
var dimType map[interface{}]DimType
var rectNum int = 10000
var realDimNum int = 0
var realTargetNum int = 0
var scatterDimNum int = 30
var scatterTargetNum int = 3
var scatterDimSize = 100000

func init() {
	dimType = make(map[interface{}]DimType)
	for j := 0; j < realDimNum; j++ {
		dimType[j] = DimTypeReal
	}

	for j := 0; j < scatterDimNum; j++ {
		dimType[j+realDimNum] = DimTypeDiscrete
	}

	for i := 0; i < rectNum; i++ {
		rect := make(Rect)
		point := make(Point)
		for j := 0; j < realDimNum; j++ {
			k := rand.Float64()
			rect[j] = Interval{MeasureFloat(k), MeasureFloat(k + 0.001)}
			point[j] = MeasureFloat(k)
		}

		for j := 0; j < realDimNum-realTargetNum; j++ {
			for {
				dim := rand.Intn(realDimNum)
				if _, ok := rect[dim]; ok {
					delete(rect, dim)
					break
				}
			}
		}

		for j := 0; j < scatterDimNum; j++ {
			k := rand.Intn(scatterDimSize)
			rect[j+realDimNum] = Scatters{MeasureFloat(k)}
			point[j+realDimNum] = MeasureFloat(k)
		}

		for j := 0; j < scatterDimNum-scatterTargetNum; j++ {
			for {
				dim := rand.Intn(scatterDimNum)
				if _, ok := rect[dim+realDimNum]; ok {
					delete(rect, dim+realDimNum)
					break
				}
			}
		}

		testRects = append(testRects, rect)
		searchPoint = append(searchPoint, point)
	}
}


func TestNewTree(t *testing.T) {
	testSize := 100

	noTreeTotal := 0
	for _, p := range searchPoint {
		for _, rect := range testRects[:testSize] {
			if rect.Contains(p) {
				//fmt.Printf("notree: %v %v %v %v\n", i, j, rect, p)
				noTreeTotal += 1
			}
		}
	}
	//fmt.Println("notree: ", noTreeTotal)

	tree := NewTree(dimType, &TreeOptions{
		TreeLevelMax:                1,
		LeafNodeMin:                 1,
		BranchingDecreasePercentMin: 0,
	})
	for i, rect := range testRects[:testSize] {
		tree.Add(rect, "data"+strconv.FormatInt(int64(i), 10))
	}
	tree.Build()

	//fmt.Println("tree:", tree.Dumps())

	treeTotal := 0
	for i, p := range searchPoint {
		d := tree.Search(p)
		if len(d) > 0 {
			fmt.Printf("tree: %v %v, %v\n", i, p, d)
		}
		treeTotal += len(d)
	}

	if noTreeTotal != treeTotal {
		t.Fatal("miss match:", noTreeTotal, " ", treeTotal)
	} else {
		//t.Log("match:", noTreeTotal, " ", treeTotal)
	}
}


func BenchmarkTree_Search(b *testing.B) {
	tree := NewTree(dimType, &TreeOptions{
		TreeLevelMax:                16,
		LeafNodeMin:                 16,
		BranchingDecreasePercentMin: 0.1,
	})
	for i, rect := range testRects {
		_ = tree.Add(rect, "data"+strconv.FormatInt(int64(i), 10))
	}
	tree.Build()

	fmt.Println("tree:", tree.Dumps())

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
