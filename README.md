## Introduction

For scenarios with many queries, index is used to optimize query performance generally.

There are two conditions:
* Build indexes to optimize data retrieval in the database.
* For some strategy engine, advertising engine, experimental platform engine and so on, need build index to optimize match pre-designed strategy constraints and retrieve quickly.

>This project mainly optimizes the second condition.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/index_common.png)

Data index is used to fast searching of data points, which implemented based on BTree, B+Tree, R-Tree, KD-Tree, Radix-Tree and other data structures.
The basic principle of these tree indexes is to narrow the query range step by step to achieve fast query by designing split hyperplanes.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/point_index.png)

For the scenarios of directional condition retrieval, the index can also be established by cutting the plane.
Cutting the hyperplane may not be able to achieve complete cutting of the area. In this project, the section that cannot be cut is assigned to a query branch to process it, and good results have been achieved.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/segment_index.png)

## Performance

**Select 10 in 10-dimensional discrete space**
>there are 100 optional values per discrete space

|Number of targeted queries|QPS without index|QPS with index|speedup|
|----|----|---|----|
|100|80998|485201|6|
|1000|7767|427899|55|
|10000|704|239578|340|
|100000|25|201207|8140|

**Select 5 in 5-dimensional numerical space**

|Number of targeted queries|QPS without index|QPS with index|speedup|
|----|----|---|----|
|100|93110|270197|3|
|1000|8904|162522|18|
|10000|863|90926|105|
|100000|43|49111|1153|

**Select 3 for 5-dimensional numerical space and 5 for 20-dimensional discrete space**

|Number of targeted queries|QPS without index|QPS with index|speedup|
|----|----|---|----|
|100|69842|173792|2|
|1000|6313|88660|14|
|10000|482|43420|90|
|100000|17|24226|1399|

## Example
