[中文](./README-ZH.md)

## Introduction

Generally, we need to establish data indexes to speed up data retrieval in the following two situations:
* Build indexes to optimize data retrieval in some database.
* Build indexes to optimize constraints matching in some strategy engine, advertising engine, experimental platform engine and so on.

>This project is working in the second situation.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/index_common.png)

Database indexes are mainly to be used to accelerate data point search, and can be constructed based on B-Tree, B+Tree, R-Tree, KD-Tree, Radix-Tree, HashTree or other data structures.
The basic principle of these data point tree indexes is to narrow the query scope by designing decomposition hyperplanes.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/point_index.png)

Constraint retrieval can also speed up retrieval by cutting the plane, but it requires some special processing. This project optimized the algorithm for selecting the cutting plane and achieved good results.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/segment_index.png)

## Performance

**Select 10 in 10-dimensional discrete space**
>there are 100 optional values per discrete space

|Number of constraints to be queried|QPS without index|QPS with index|speedup|
|----:|----:|---:|----:|
|100|80998|485201|6|
|1000|7767|427899|55|
|10000|704|239578|340|
|100000|25|201207|8140|

**Select 5 in 5-dimensional numerical space**

|Number of constraints to be queried|QPS without index|QPS with index|speedup|
|----:|----:|---:|----:|
|100|93110|270197|3|
|1000|8904|162522|18|
|10000|863|90926|105|
|100000|43|49111|1153|

**Select 3 for 5-dimensional numerical space and 5 for 20-dimensional discrete space**

|Number of constraints to be queried |QPS without index|QPS with index|speedup|
|----:|----:|---:|----:|
|100|69842|173792|2|
|1000|6313|88660|14|
|10000|482|43420|90|
|100000|17|24226|1399|

## Example

    // build a new tree
	tree1 := NewTree(DimTypes{
		"Field1": DimTypeDiscrete, // discrete space
		"Field2": DimTypeReal, // real space
	}, &TreeOptions{
		TreeLevelMax:                16, // tree's max height
		LeafNodeMin:                 4, // max data number within leaf node
		BranchingDecreasePercentMin: 0.1, // min split ratio
	})

	err := tree1.Add(Rect{
		"Field1": Scatters{MeasureString("one"), MeasureString("two"), MeasureString("three")}, // targeting discrete string 
		"Field2": Interval{MeasureFloat(0.1), MeasureFloat(2.0)}}, // target real interval
		"target1")
	if err != nil {
		log.Fatal("node add error:", err)
	}
	tree1.Build() // build tree

    // search point
	result := tree1.Search(Point{"Field1": MeasureString("one"), "Field2": MeasureFloat(0.3)})

	if len(result) != 1 || result[0].(string) != "target1" {
		log.Fatal("tree search error")
	}