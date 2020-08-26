## Introduction

For scenarios with a large number of queries, an index is generally required to optimize query performance.
* In the data warehouse, query data is optimized by establishing indexes.
* In some strategy engines, advertising engines or traffic experiment systems, build indexes to optimize targeting strategies that meet pre-made conditions from data points.

>This project is mainly for the second scenario.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/index_common.png)

For the index of data points, common indexes include BTree, B+Tree, R-Tree, KD-Tree, Radix-Tree, etc.
The basic principle of these indexes is to narrow the query range and speed up the query by cutting the plane.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/point_index.png)

The index of directional condition retrieval can also be indexed by cutting plane, but the area cannot be completely cut according to the cutting hyperplane.
This project makes additional treatment for these areas that cannot be cut, and the cutting effect is guaranteed.

![avatar](https://github.com/boostlearn/go-kd-segment-tree/raw/master/doc/segment_index.png)

## Performance

**Select 10 in 10-dimensional discrete space**
>100 optional values ​​per discrete space

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
