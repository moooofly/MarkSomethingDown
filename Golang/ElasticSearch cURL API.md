# ElasticSearch cURL API

## 针对数据操作

### 增

```shell
POST /<_index>/<_type>
{<JSON>}
```

### 删

```shell
# 删除整个 _index （删库）
DELETE /<_index>

# 基于通配方式删除所有匹配的 _index （删所有匹配库）
DELETE /<_index_regex>

# 删除指定 _index 下的指定 _type （删表）
DELETE /<_index>/<_type>

# 基于通配方式删除指定 _index 下所有匹配 _type （删所有匹配表）
DELETE /<_index>/<_type_regex>

# 删除单条数据 （删行）
DELETE /<_index>/<_type>/<_id>

## 在不删除 index 和 type 的前提下，清除其中的所有数据
DELETE /<_index>/<_type>/_query
{
　　"query" : {
　　    "match_all" : {}
　　}
}
```

> 注意：针对 index 或 index/type 的删除，会导致相应的 mapping 也同时被删除；

### 查

```shell
# 获取指定 _index 内容（内容像 index template 但不是）
GET /<_index>

# 基于通配方式获取 _index 内容
GET /<_index_regex>

# 获取指定 _index 下 _type 中 _id 对应内容 （获取行数据）
GET /<_index>/<_type>/<_id>

# 只获取 JSON 源数据部分（存入什么得到什么）
GET /<_index>/<_type>/<_id>/_source

# 只获取源数据中的特定字段 （获取指定行下指定的列）
GET /<_index>/<_type>/<_id>?fields=<fieldM>,...,<fieldN>
```

### 改

```shell
# 全量提交
POST /<_index>/<_type>/<_id>
{<JSON>}

# 局部更新 （更新列）
POST /<_index>/<_type>/<_id>/_update
{
    "doc": {
        "fieldN":"yyy"
    }
}
```

### _search 搜索

- 空搜索（在所有 _index 的所有 _type 中搜索）

```shell
# 简化形式
GET /_search

# 完整形式
GET /_all/_search
```

- 其他搜索

```shell
# 全文搜索（针对 _source 中包含的全部 fields）
GET /<_index>/<_type>/_search?q=<key>

# 针对单个字段（field）上的全文搜索
GET /<_index>/<_type>/_search?q=<fieldN>:<key>

# 针对单个字段（field）上的全文精确搜索
GET /<_index>/<_type>/_search?q=<fieldN>:"<key>"

# 在索引 _index 的所有 _type 中搜索
GET /<_index>/_search

# 在索引 _indexM 和 _indexN 的所有 _type 中搜索
GET /<_indexM>,<_indexN>/_search

# 在以 g 或 u 开头索引的所有 _type 中搜索
GET /g*,u*/_search

# 在索引 _indexM 的 user 类型中搜索
GET /<_indexM>/user/_search

# 在 _indexM 和 _indexN 索引的 user 和 tweet 类型中搜索
GET /<_indexM>,<_indexN>/user,tweet/_search

# 在所有索引的 user 和 tweet 类型中搜索
GET /_all/user,tweet/_search
```

另外，如下两种形式相互等价

```
GET /<_index>/<_type>/_search
```

等价于

```
GET /<_index>/_search?type=<_type>
```


> 当搜索仅包含单一 _index 时，Elasticsearch 会转发搜索请求到该 _index 的主分片或每个分片的复制分片上，然后再聚合每个分片的结果。当搜索包含多个 _index 时，也是同样的方式，只不过会有更多的分片被关联。

### 聚合

略


### 分页

```
# 每页只显示 5 条结果
GET /_search?size=5

# 每页只显示 5 条结果，从第 6 条开始显示
GET /_search?size=5&from=5
```

> 分页时应该避免分页太深（即 from 值过大）或者一次请求太多的结果（即 size 值过大）；原因在于：结果在返回前会被排序，而一个搜索请求常常涉及多个分片，每个分片会生成自己排好序的结果，之后需要集中起来进行排序，以确保提供整体排序正确的结果。因此，若请求 `GET /_search?size=10&from=1000` 涉及 5 个分片，则每个分片都必须产生排序后的 10010 个结果，然后聚合操作还要排序这 50050 个结果，并最终丢弃 50040 个！


## 针对映射（mapping）操作

### 增

```shell
PUT /<_index>/_mapping
{<JSON>}
```

### 删

```shell
DELETE /<_index>/_mapping/<_type>
```

> 注意：写入数据会自动添加映射，但删除数据不会删除数据的映射；存在一种特例：若删除整个索引 <_index> ，则映射将同时被删除；

### 查

```shell
GET /<_index>/_mapping/<_type>
```

### 改

```shell
PUT /<_index>/_mapping/<_type>
{<JSON>}
```

> 注意：更新只对新字段有效，已经生成的字段映射是不可变更的；如果需要变更，则需要使用 `reindex` 方法解决；


## 针对模版（template）操作


### 增

```
PUT /_template/<template_name>
{
  "template": "te*",
  "settings": {
    "number_of_shards": 1
  },
  "mappings": {
    "type1": {
      "_source": {
        "enabled": false
      },
      "properties": {
        "host_name": {
          "type": "keyword"
        },
        "created_at": {
          "type": "date",
          "format": "EEE MMM dd HH:mm:ss Z YYYY"
        }
      }
    }
  }
}
```

### 删

```
DELETE /_template/<template_name>
```

### 查

```
# 获取全部 index template
GET /_template

# 获取指定 index template
GET /_template/<template_name>

# 基于通配符获取多个 index template
GET /_template/<template_name_regex>

# 直接指定多个 index template 名进行查询
GET /_template/<template_M>,<template_N>

# 仅确定目标 index template 是否存在
HEAD /_template/<template_name>
```

### 获取指定字段

```
GET /_template/<template_name>?filter_path=*.version
```


## ingest node pipeline

```
# 创建 ingest node pipeline
PUT /_ingest/pipeline/<pipeline_name>
{<JSON>}
```


## 获取 elasticsearch 元信息

```
GET /
```



