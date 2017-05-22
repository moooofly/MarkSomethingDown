# ElasticSearch cURL API

## 针对数据操作

### 增

```shell
curl -u <usrname>:<passwd> -XPOST http://xxx/<_index>/<_type> -d '{<JSON>}'
```

### 删

```shell
# 删除整个 _index 
curl -u <usrname>:<passwd> -XDELETE http://xxx/<_index>

# 删除指定 _index 下的整个 _type
curl -u <usrname>:<passwd> -XDELETE http://xxx/<_index>/<_type>

# 删除单条数据
curl -u <usrname>:<passwd> -XDELETE http://xxx/<_index>/<_type>/<_id>

# 基于通配方式阐述所有 _index
curl -u <usrname>:<passwd> -XDELETE http://xxx/<_index_regex>
```

### 查

```shell
# 获取指定 _index 下 _type 中 _id 对应内容
curl -u <usrname>:<passwd> -XGET http://xxx/<_index>/<_type>/<_id>

# 只获取源数据部分
curl -u <usrname>:<passwd> -XGET http://xxx/<_index>/<_type>/<_id>/_source

# 只获取源数据中的特定字段
curl -u <usrname>:<passwd> -XGET http://xxx/<_index>/<_type>/<_id>?fields=<fieldM>,...,<fieldN>
```

### 改

```shell
# 全量提交
curl -u <usrname>:<passwd> -XPOST http://xxx/<_index>/<_type>/<_id> -d '{<JSON>}'

# 局部更新
curl -u <usrname>:<passwd> -XPOST http://xxx/<_index>/<_type>/<_id>/_update -d '{
    "doc": {
        "fieldN":"yyy"
    }
}'
```

### 搜索

```shell
# 全文搜索（针对 _source 中包含的全部 fields）
curl -u <usrname>:<passwd> -XGET http://xxx/<_index>/<_type>/_search?q=<key>

# 针对单个字段（field）上的全文搜索
curl -u <usrname>:<passwd> -XGET http://xxx/<_index>/<_type>/_search?q=<fieldN>:<key>

# 针对单个字段（field）上的全文精确搜索
curl -u <usrname>:<passwd> -XGET http://xxx/<_index>/<_type>/_search?q=<fieldN>:"<key>"
```

另外，如下两种形式相互等价

```
curl -u <usrname>:<passwd> -XGET http://xxx/<_index>/<_type>/_search
```

等价于

```
curl -u <usrname>:<passwd> -XGET http://xxx/<_index>/_search?type=<_type>
```


### 聚合

略

## 实际使用

- 空搜索（在所有索引的所有类型中搜索）

```
GET /_search
```

- 多 _index + 多 _type

```
# 在索引 gb 的所有类型中搜索
GET /gb/_search

# 在索引 gb 和 us 的所有类型中搜索
GET /gb,us/_search

# 在以 g 或 u 开头的索引的所有类型中搜索
GET /g*,u*/_search

# 在索引 gb 的类型 user 中搜索
GET /gb/user/_search

# 在 gb 和 us 索引的 user 和 tweet 类型中搜索
GET /gb,us/user,tweet/_search

# 在所有索引的 user 和 tweet 类型中搜索
GET /_all/user,tweet/_search
```

> 当搜索仅包含单一 _index 时，Elasticsearch 会转发搜索请求到该 _index 的主分片或每个分片的复制分片上，然后再聚合每个分片的结果。当搜索包含多个 _index 时，也是同样的方式，只不过会有更多的分片被关联。

- 分页

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
curl -u <usrname>:<passwd> -XPUT http://xxx/<_index>/_mapping -d '{<JSON>}'
```

### 删

```shell
curl -u <usrname>:<passwd> -XDELETE http://xxx/<_index>/_mapping/<_type>
```

> 注意：写入数据会自动添加映射，但删除数据不会删除数据的映射；存在一种特例：若删除整个索引 <_index> ，则映射将同时被删除；

### 查

```shell
curl -u <usrname>:<passwd> -XGET http://xxx/<_index>/_mapping/<_type>
```

### 改

```shell
curl -u <usrname>:<passwd> -XPUT http://xxx/<_index>/_mapping/<_type> -d '{<JSON>}'
```

> 注意：更新只对新字段有效，已经生成的字段映射是不可变更的；如果需要变更，则需要使用 `reindex` 方法解决；


## 针对模版（template）操作


### 增

```
curl -XPUT 'localhost:9200/_template/template_1?pretty' -H 'Content-Type: application/json' -d'
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
'
```

### 删

```
curl -XDELETE 'localhost:9200/_template/template_1?pretty'
```

### 查

```
# 查询单个 index template
curl -XGET 'localhost:9200/_template/template_1?pretty'

# 使用通配符匹配多个 index template
curl -XGET 'localhost:9200/_template/temp*?pretty'

# 直接指定多个 index template 名
curl -XGET 'localhost:9200/_template/template_1,template_2?pretty'

# 获取全部 index template 列表
curl -XGET 'localhost:9200/_template?pretty'

# 仅确定目标 index template 是否存在
curl -XHEAD 'localhost:9200/_template/template_1?pretty'
```

### 获取指定字段

```
curl -XGET 'localhost:9200/_template/template_1?filter_path=*.version&pretty'
```







