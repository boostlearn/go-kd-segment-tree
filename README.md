## Introduction

Generally, we need to establish data indexes to speed up data retrieval in the following two situations:
* Build indexes to optimize data retrieval in some database.
* Build indexes to optimize constraints matching in some strategy engine, advertising engine, experimental platform engine and so on.

>This project is working in the second situation.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/index_common.png)

Database index is mainly used to accelerate data point search, and can be constructed based on BTree, B + Tree, R-Tree, KD-Tree, Radix-Tree or other data structures.
The basic principle of these data point tree indexes is to narrow the query scope by designing decomposition hyperplanes.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/point_index.png)

Constraint retrieval can also speed up retrieval by cutting the plane, but it requires some special processing. This project optimized the algorithm for selecting the cutting plane and achieved good results.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/segment_index.png)

## Performance

**Select 10 in 10-dimensional discrete space**
>there are 100 optional values per discrete space

|Number of constraints|QPS without index|QPS with index|speedup|
|----|----|---|----|
|100|80998|485201|6|
|1000|7767|427899|55|
|10000|704|239578|340|
|100000|25|201207|8140|

**Select 5 in 5-dimensional numerical space**

|Number of constraints|QPS without index|QPS with index|speedup|
|----|----|---|----|
|100|93110|270197|3|
|1000|8904|162522|18|
|10000|863|90926|105|
|100000|43|49111|1153|

**Select 3 for 5-dimensional numerical space and 5 for 20-dimensional discrete space**

|Number of constraints |QPS without index|QPS with index|speedup|
|----|----|---|----|
|100|69842|173792|2|
|1000|6313|88660|14|
|10000|482|43420|90|
|100000|17|24226|1399|

## Example
