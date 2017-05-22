# .kibana 文件说明

对于 kibana 4 或 5 来说，在 `kibana.yml` 中设置的默认 index 值为 ".kibana" ；

kibana 配置文件中的内容如下所示：

```
# Kibana uses an index in Elasticsearch to store saved searches, visualizations and
# dashboards. Kibana creates a new index if the index doesn’t already exist.
#kibana.index: ".kibana"
```

可以通过 REST API 接口获取和 `.kibana` 相关的内容；

## 获取 `.kibana` 本身的内容

```
curl -u elastic:changeme -XGET http://localhost:9200/.kibana?pretty
```

返回值如下

```
{
  ".kibana" : {
    "aliases" : { },
    "mappings" : {
      "dashboard" : {
            ...
      },
      "server" : {
            ...
      },
      "timelion-sheet" : {
            ...
      },
      "visualization" : {
            ...
      },
      "search" : {
            ...
      },
      "graph-workspace" : {
            ...
      },
      "config" : {
            ...
      },
      "index-pattern" : {
            ...
      }
    },
    "settings" : {
        ...
    }
  }
}
```

可以看到

- 索引 `.kibana` 的内容主要由 "mappings" 和 "settings" 构成（**index template** 正是由这两部分构成的）；
- "mappings" 的内容由 `dashboard`/`server`/`timelion-sheet`/`visualization`/`search`/`graph-workspace`/`config`/`index-pattern` 等**映射类型（mapping type）**构成；
- 针对每一种映射类型，通过 properties 定义针对每一个 field 的**映射选项**（类型/索引方式/多重索引等）；


> 关于 mapping 的补充说明：
>
> - mapping 用于定义 documents 以及其包含的 fields 是如何被**存储**和**索引**的；
>
> ------
> 关于 **mapping type** 的补充说明：
> 
> - 分为 meta-fields (\_index/\_type/\_id/\_source) 和 fields/properties 两部分；
> - 每种 index 均具有一种或多种映射类型，用于将 index 中的所有 documents 分为不同的逻辑分组；


## 获取 `.kibana` 索引下所有 _type 信息

```
curl -u elastic:changeme -XGET http://localhost:9200/.kibana/_search?pretty
```

返回值如下

```
{
  "took" : 1,
  "timed_out" : false,
  "_shards" : {
    "total" : 1,
    "successful" : 1,
    "failed" : 0
  },
  "hits" : {
    "total" : 268,
    "max_score" : 1.0,
    "hits" : [
      {
        "_index" : ".kibana",
        "_type" : "index-pattern",
        "_id" : "filebeat-*",
        "_score" : 1.0,
        "_source" : {
            ...
        }
      },
      {
        "_index" : ".kibana",
        "_type" : "dashboard",
        "_id" : "Filebeat-Apache2-Dashboard",
        "_score" : 1.0,
        "_source" : {
            ...
        }
      },
      {
        "_index" : ".kibana",
        "_type" : "visualization",
        "_id" : "Apache2-response-codes-of-top-URLs",
        "_score" : 1.0,
        "_source" : {
            ...
        }
      },
      {
        "_index" : ".kibana",
        "_type" : "search",
        "_id" : "Apache2-access-logs",
        "_score" : 1.0,
        "_source" : {
            ...
        }
      },
      ...
```

可以看到

- 能够获取到 `.kibana` 下包含的 `index-pattern`/`dashboard`/`visualization`/`search` 相关信息；


----------


另外，可以通过如下 API 获取指定 _index 下指定 _type 的信息

```
GET /.kibana/_search?type=dashboard&pretty=1

GET /.kibana/dashboard/_search?pretty=1
```





